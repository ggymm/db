package table

import (
	"errors"
	"fmt"
	"github.com/ggymm/db"
	"github.com/ggymm/db/index"
	"math"
	"slices"
	"sync"

	"github.com/ggymm/db/boot"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/cmap"
	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/pkg/view"
	"github.com/ggymm/db/ver"
)

var (
	ErrNoSuchTable  = errors.New("no such table")
	ErrNoPrimaryKey = errors.New("no primary key")
)

type Manage interface {
	Begin(level int) uint64
	Commit(txId uint64) error
	Rollback(txId uint64)

	Create(txId uint64, stmt *sql.CreateStmt) error
	Insert(txId uint64, stmt *sql.InsertStmt) error
	Update(txId uint64, stmt *sql.UpdateStmt) error
	Delete(txId uint64, stmt *sql.DeleteStmt) error
	Select(txId uint64, stmt *sql.SelectStmt) ([]entry, error)

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

func (tbm *tableManage) Commit(txId uint64) error {
	return tbm.verManage.Commit(txId)
}

func (tbm *tableManage) Rollback(txId uint64) {
	tbm.verManage.Rollback(txId)
}

func (tbm *tableManage) Create(txId uint64, stmt *sql.CreateStmt) error {
	tbm.mu.Lock()
	defer tbm.mu.Unlock()
	if exist := tbm.tables.Has(stmt.Name); exist {
		return nil
	}

	t := new(table)
	t.tbm = tbm
	t.Name = stmt.Name
	t.Next = tbm.readTableId() // 上一张表的 itemId
	t.Fields = make([]*field, 0)

	// 读取 主键 和 索引
	indexList := make([]string, 0)
	for _, i := range stmt.Table.Index {
		indexList = append(indexList, i.Field)
	}

	// 读取 field
	for _, tf := range stmt.Table.Field {
		f := new(field)
		f.Name = tf.Name
		f.Type = tf.Type.String()
		f.TreeId = 0
		f.Default = tf.Default
		f.Nullable = tf.Nullable

		indexed := false
		if slices.Contains(indexList, tf.Name) {
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
			i, err := index.NewIndex(tbm.DataManage(), &db.Option{
				Open: false,
			})
			if err != nil {
				return err
			}
			f.index = i
			f.TreeId = i.GetBootId()
		}

		// 保存字段信息
		err := f.persist(txId)
		if err != nil {
			return err
		}
		t.Fields = append(t.Fields, f)
	}

	// 保存表信息
	err := t.persist(txId)
	if err != nil {
		return err
	}

	// 更新表信息
	tbm.tables.Set(t.Name, t)
	tbm.updateTableId(t.itemId) // 更新为当前表的 itemId
	return err
}

func (tbm *tableManage) Insert(txId uint64, stmt *sql.InsertStmt) error {
	var (
		err error
		row map[string]string
	)

	// 获取表对象
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	// 格式化插入数据
	row, err = stmt.FormatData()
	if err != nil {
		return err
	}

	// 构建数据
	raw := make([]byte, 0)
	for _, f := range t.Fields {
		// 获取字段值
		val, exist := row[f.Name]
		if !exist {
			if len(f.Default) != 0 {
				val = f.Default
			} else if !f.Nullable {
				return fmt.Errorf("field %s is not allowed to be null", f.Name)
			}
		}

		// 获取字段二进制值
		raw = sql.FieldRaw(f.Type, val)
	}

	// 写入数据
	itemId, err1 := tbm.verManage.Write(txId, raw)
	if err1 != nil {
		return err1
	}

	// 判断是否有字段需要索引
	for _, f := range t.Fields {
		if f.isIndex() { // 主键同时也是索引
			v, exist := row[f.Name]
			if !exist {
				return errors.New("index is not allowed to be null")
			}

			// 格式化索引字段
			val := sql.FieldFormat(f.Type, v)
			err = f.index.Insert(hash(val), itemId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tbm *tableManage) Update(txId uint64, stmt *sql.UpdateStmt) error {
	_, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}
	return nil
}

func (tbm *tableManage) Delete(txId uint64, stmt *sql.DeleteStmt) error {
	_, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	return nil
}

func (tbm *tableManage) Select(txId uint64, stmt *sql.SelectStmt) ([]entry, error) {
	// 获取表对象
	t, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return nil, ErrNoSuchTable
	}

	var (
		err     error
		itemIds []uint64
	)

	// 遍历条件，没有查询条件则查询全部数据
	if len(stmt.Where) == 0 {
		// 找到主键索引字段
		f := &field{}
		for _, tf := range t.Fields {
			if tf.PrimaryKey {
				f = tf
				break
			}
		}
		if f == nil {
			return nil, ErrNoPrimaryKey
		}

		// 查询全部数据的
		itemIds, err = f.index.SearchRange(0, math.MaxUint64)
		if err != nil {
			return nil, err
		}
	} else {
		for _, cond := range stmt.Where {
			fmt.Printf("cond: %+v\n", cond)
		}
	}

	// 读取数据
	raw := make([]byte, 0)
	rows := make([]entry, 0)
	for _, itemId := range itemIds {
		raw, ok, err = tbm.verManage.Read(txId, itemId)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		// 解析数据
		pos := 0
		row := make(entry)
		for _, f := range t.Fields {
			val, shift := sql.FieldParse(f.Type, raw[pos:])
			if err != nil {
				return nil, err
			}
			row[f.Name] = val
			pos += shift
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (tbm *tableManage) ShowTable() string {
	thead := []string{"Tables"}
	tbody := make([][]string, 0)

	tbm.tables.IterCb(func(name string, _ *table) {
		tbody = append(tbody, []string{name})
	})

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(thead)
	vt.SetBody(tbody)
	return vt.String()
}

func (tbm *tableManage) ShowField(table string) string {
	thead := []string{"Field", "Type", "Null", "Key", "Default"}
	tbody := make([][]string, 0)
	t, exist := tbm.tables.Get(table)
	if !exist {
		return "no such table"
	}

	for _, f := range t.Fields {
		indexed := ""
		if f.isIndex() {
			indexed = "YES"
			if f.PrimaryKey {
				indexed = "PRI"
			}
		}
		nullable := "NO"
		if f.Nullable {
			nullable = "YES"
		}
		tbody = append(tbody, []string{
			f.Name,
			f.Type,
			nullable,
			indexed,
			f.Default,
		})
	}

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(thead)
	vt.SetBody(tbody)
	return vt.String()
}

func (tbm *tableManage) ShowResult(table string, entries []entry) string {
	// 获取表对象
	t, ok := tbm.tables.Get(table)
	if !ok {
		return ""
	}

	thead := make([]string, 0)
	tbody := make([][]string, 0)

	for _, f := range t.Fields {
		thead = append(thead, f.Name)
	}

	for _, ent := range entries {
		row := make([]string, 0)
		for _, f := range thead {
			val, exist := ent[f]
			if !exist {
				val = ""
			}
			row = append(row, fmt.Sprintf("%v", val))
		}
		tbody = append(tbody, row)
	}

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(thead)
	vt.SetBody(tbody)
	return vt.String()
}

func (tbm *tableManage) VerManage() ver.Manage {
	return tbm.verManage
}

func (tbm *tableManage) DataManage() data.Manage {
	return tbm.dataManage
}
