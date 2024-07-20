package table

import (
	"fmt"
	"slices"
	"sync"

	"github.com/ggymm/db"
	"github.com/ggymm/db/boot"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/index"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/cmap"
	"github.com/ggymm/db/pkg/hash"
	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/pkg/view"
	"github.com/ggymm/db/ver"
)

type Manage interface {
	Begin(level int) uint64
	Commit(tid uint64) error
	Rollback(tid uint64)

	Create(tid uint64, stmt *sql.CreateStmt) (err error)
	Insert(tid uint64, stmt *sql.InsertStmt) (err error)
	Update(tid uint64, stmt *sql.UpdateStmt) (err error)
	Delete(tid uint64, stmt *sql.DeleteStmt) (err error)
	Select(tid uint64, stmt *sql.SelectStmt) ([]entry, error)

	ShowTable() string
	ShowField(table string) string
	ShowResult(table string, entries []entry) string

	VerManage() ver.Manage
	DataManage() data.Manage
}

type tableManage struct {
	mu     sync.Mutex
	tables cmap.CMap[string, *table]

	boot       boot.Boot
	verManage  ver.Manage
	dataManage data.Manage
}

func NewManage(boot boot.Boot, verManage ver.Manage, dataManage data.Manage) Manage {
	tbm := &tableManage{
		tables: cmap.New[*table](),

		boot:       boot,
		verManage:  verManage,
		dataManage: dataManage,
	}

	itemId := tbm.readTableId()
	for itemId != 0 {
		t := readTable(tbm, itemId)
		tbm.tables.Set(t.Name, t)

		// 读取下一个表的信息
		itemId = t.Next
	}
	return tbm
}

func (tbm *tableManage) readTableId() uint64 {
	return bin.Uint64(tbm.boot.Load())
}

func (tbm *tableManage) updateTableId(id uint64) {
	tbm.boot.Update(bin.Uint64Raw(id))
}

func (tbm *tableManage) Begin(level int) uint64 {
	return tbm.verManage.Begin(level)
}

func (tbm *tableManage) Commit(tid uint64) error {
	return tbm.verManage.Commit(tid)
}

func (tbm *tableManage) Rollback(tid uint64) {
	tbm.verManage.Rollback(tid)
}

func (tbm *tableManage) Create(tid uint64, stmt *sql.CreateStmt) (err error) {
	tbm.mu.Lock()
	defer tbm.mu.Unlock()
	if exist := tbm.tables.Has(stmt.Name); exist {
		return
	}

	t := new(table)
	t.tbm = tbm
	t.Name = stmt.Name
	t.Next = tbm.readTableId() // 上一张表的 itemId
	t.Fields = make([]*field, 0)

	// 读取 主键 和 索引
	indexes := make([]string, 0)
	for _, i := range stmt.Table.Index {
		indexes = append(indexes, i.Field)
	}

	// 读取 field
	for _, tf := range stmt.Table.Field {
		f := new(field)
		f.tbm = tbm
		f.Name = tf.Name
		f.Type = tf.Type.String()
		f.TreeId = 0
		f.Default = tf.Default
		f.Nullable = tf.Nullable

		indexed := false
		if slices.Contains(indexes, tf.Name) {
			indexed = true
			f.Nullable = false
		}

		// 如果是主键
		// 则不允许为空，且是索引
		if stmt.Table.Pk.Field == tf.Name {
			indexed = true
			f.Nullable = false
			f.PrimaryKey = true
		}

		if indexed {
			i, err1 := index.NewIndex(tbm.DataManage(), &db.Option{
				Open: false,
			})
			if err1 != nil {
				return err1
			}
			f.index = i
			f.TreeId = i.GetBootId()
		}

		// 保存字段信息
		err = f.save(tid)
		if err != nil {
			return err
		}
		t.Fields = append(t.Fields, f)
	}

	// 保存表信息
	err = t.save(tid)
	if err != nil {
		return err
	}

	// 更新表信息
	tbm.tables.Set(t.Name, t)
	tbm.updateTableId(t.itemId) // 更新为当前表的 itemId
	return
}

func (tbm *tableManage) Insert(tid uint64, stmt *sql.InsertStmt) (err error) {
	// 获取表对象
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	// 格式化插入数据
	if len(stmt.Value) != len(stmt.Field) {
		return ErrInsertNotMatch
	}
	row := make(map[string]any)
	for _, f := range t.Fields {
		i := slices.Index(stmt.Field, f.Name)
		if i == -1 {
			continue
		}
		row[f.Name] = sql.FormatVal(f.Type, stmt.Value[i])
	}

	// 构建数据
	raw, err := t.wrapRaw(row)
	if err != nil {
		return err
	}

	// 写入数据
	rid, err := tbm.verManage.Write(tid, raw)
	if err != nil {
		return err
	}

	// 判断是否有字段需要索引
	for _, f := range t.Fields {
		if f.TreeId != 0 {
			v, exist := row[f.Name]
			if !exist {
				return NewError(ErrNotAllowNull, f.Name)
			}

			// 格式化索引字段
			err = f.index.Insert(hash.Sum64(v), rid)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (tbm *tableManage) Update(tid uint64, stmt *sql.UpdateStmt) (err error) {
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	if len(stmt.Where) == 0 {
		return ErrMustHaveCondition
	}

	var (
		raw  = make([]byte, 0)
		rids = make([]uint64, 0)
	)

	// 查询解析
	rids, err = resolveWhere(t, stmt.Where)
	if err != nil {
		return err
	}

	// 读取数据
	for _, rid := range rids {
		raw, ok, err = tbm.verManage.Read(tid, rid)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		// 解析数据
		row := t.wrapEntry(raw, stmt.Where)
		if row == nil {
			continue
		}

		// 删除数据
		_, err = tbm.verManage.Delete(tid, rid)
		if err != nil {
			return err
		}

		// 更新数据
		for k, v := range stmt.Value {
			if _, ok = row[k]; ok {
				row[k] = v
			}
		}
		raw, err = t.wrapRaw(row)
		if err != nil {
			return err
		}
		rid, err = tbm.verManage.Write(tid, raw)
		if err != nil {
			return err
		}

		// 更新索引
		for _, f := range t.Fields {
			if f.TreeId != 0 {
				v, exist := row[f.Name]
				if !exist {
					return NewError(ErrNotAllowNull, f.Name)
				}

				// 格式化索引字段
				err = f.index.Insert(hash.Sum64(v), rid)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

func (tbm *tableManage) Delete(tid uint64, stmt *sql.DeleteStmt) (err error) {
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	if len(stmt.Where) == 0 {
		return ErrMustHaveCondition
	}

	// 查询解析
	rids, err := resolveWhere(t, stmt.Where)
	if err != nil {
		return err
	}

	// 读取数据
	for _, rid := range rids {
		_, err = tbm.verManage.Delete(tid, rid)
		if err != nil {
			return err
		}
	}
	return
}

// Select 查询数据
func (tbm *tableManage) Select(tid uint64, stmt *sql.SelectStmt) ([]entry, error) {
	// 获取表对象
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return nil, ErrNoSuchTable
	}

	var (
		err  error
		rids = make([]uint64, 0)
	)

	// 查询条件
	rids, err = resolveWhere(t, stmt.Where)
	if err != nil {
		return nil, err
	}

	var (
		raw  = make([]byte, 0)
		rows = make([]entry, 0)
	)

	// 读取数据
	for _, rid := range rids {
		raw, ok, err = tbm.verManage.Read(tid, rid)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		// 解析数据
		row := t.wrapEntry(raw, stmt.Where)
		if row != nil {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (tbm *tableManage) ShowTable() string {
	head := []string{"Tables"}
	body := make([][]string, 0)

	tbm.tables.IterCb(func(name string, _ *table) {
		body = append(body, []string{name})
	})

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(head)
	vt.SetBody(body)
	return vt.String()
}

func (tbm *tableManage) ShowField(table string) string {
	head := []string{"Field", "Type", "Null", "Key", "Default"}
	body := make([][]string, 0)
	t, exist := tbm.tables.Get(table)
	if !exist {
		return "no such table"
	}

	for _, f := range t.Fields {
		indexed := ""
		if f.TreeId != 0 {
			indexed = "YES"
			if f.PrimaryKey {
				indexed = "PRI"
			}
		}
		nullable := "NO"
		if f.Nullable {
			nullable = "YES"
		}
		body = append(body, []string{
			f.Name,
			f.Type,
			nullable,
			indexed,
			f.Default,
		})
	}

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(head)
	vt.SetBody(body)
	return vt.String()
}

func (tbm *tableManage) ShowResult(table string, entries []entry) string {
	// 获取表对象
	t, ok := tbm.tables.Get(table)
	if !ok {
		return ""
	}

	head := make([]string, 0)
	body := make([][]string, 0)

	for _, f := range t.Fields {
		head = append(head, f.Name)
	}

	for _, ent := range entries {
		row := make([]string, 0)
		for _, f := range head {
			val, exist := ent[f]
			if !exist {
				val = ""
			}
			row = append(row, fmt.Sprintf("%v", val))
		}
		body = append(body, row)
	}

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(head)
	vt.SetBody(body)
	return vt.String()
}

func (tbm *tableManage) VerManage() ver.Manage {
	return tbm.verManage
}

func (tbm *tableManage) DataManage() data.Manage {
	return tbm.dataManage
}
