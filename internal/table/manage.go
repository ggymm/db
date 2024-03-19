package table

import (
	"db/pkg/view"
	"errors"
	"fmt"
	"sync"

	"db/pkg/utils"

	"db/internal/boot"
	"db/internal/data"
	"db/internal/ver"
	"db/pkg/bin"
	"db/pkg/cmap"
	"db/pkg/sql"
)

var ErrNoSuchTable = errors.New("no such table")

type Manage interface {
	Begin(level int) uint64
	Abort(txId uint64)
	Commit(txId uint64) error

	Create(txId uint64, stmt *sql.CreateStmt) error
	Insert(txId uint64, stmt *sql.InsertStmt) error
	Update(txId uint64, stmt *sql.UpdateStmt) error
	Delete(txId uint64, stmt *sql.DeleteStmt) error
	Select(txId uint64, stmt *sql.SelectStmt) ([]byte, error)

	ShowTable() string
	ShowField(table string) string

	VerManage() ver.Manage
	DataManage() data.Manage
}

type tableManage struct {
	boot       boot.Boot
	verManage  ver.Manage
	dataManage data.Manage

	lock   sync.Mutex
	tables cmap.CMap[string, *table]
}

func NewManage(boot boot.Boot, verManage ver.Manage, dataManage data.Manage) Manage {
	tbm := &tableManage{
		boot:       boot,
		verManage:  verManage,
		dataManage: dataManage,
		tables:     cmap.New[*table](),
	}

	id := tbm.readTableId()
	for id != 0 {
		t := readTable(tbm, id)
		tbm.tables.Set(t.tableName, t)

		// 读取下一个表的信息
		id = t.tableNext
	}
	return tbm
}

func (tbm *tableManage) readTableId() uint64 {
	return bin.Uint64(tbm.boot.Load())
}

func (tbm *tableManage) updateTableId(id uint64) {
	raw := bin.Uint64Raw(id)
	tbm.boot.Update(raw)
}

func (tbm *tableManage) Begin(level int) uint64 {
	return tbm.verManage.Begin(level)
}

func (tbm *tableManage) Abort(txId uint64) {
	tbm.verManage.Abort(txId)
}

func (tbm *tableManage) Commit(txId uint64) error {
	return tbm.verManage.Commit(txId)
}

func (tbm *tableManage) Create(txId uint64, stmt *sql.CreateStmt) error {
	tbm.lock.Lock()
	defer tbm.lock.Unlock()

	if exist := tbm.tables.Has(stmt.Name); exist {
		return nil
	}

	t, err := createTable(tbm, &newTable{
		txId: txId,

		tableName: stmt.Name,
		tableNext: tbm.readTableId(),

		index: stmt.Table.Index,
		field: stmt.Table.Field,
	})
	if err != nil {
		return err
	}

	tbm.tables.Set(t.tableName, t)
	tbm.updateTableId(t.itemId)
	return nil
}

func (tbm *tableManage) Insert(txId uint64, stmt *sql.InsertStmt) error {
	var (
		ok  bool
		err error

		t    *table
		maps []map[string]string
	)

	// 获取表对象
	t, ok = tbm.tables.Get(stmt.Table)
	if !ok {
		return ErrNoSuchTable
	}

	// 格式化插入数据
	maps, err = stmt.Format()
	if err != nil {
		return err
	}

	// 构建插入的数据条目
	length := len(t.tableFields)
	entries := make([]entry, length)
	for _, item := range maps {
		ent := entry{
			raw:    make([]byte, 0),
			value:  make([]any, length),
			fields: make([]*field, length),
		}

		val := ""
		for i, f := range t.tableFields {
			ent.fields[i] = f

			// 获取字段值
			val, ok = item[f.fieldName]
			switch {
			case ok:
				ent.value[i] = val
			case len(f.defaultVal) != 0:
				ent.value[i] = f.defaultVal
			case f.allowNull:
				ent.value[i] = nil
			default:
				return fmt.Errorf("field %s is not allowed to be null", f.fieldName)
			}

			// 获取字段二进制值
			ent.raw = append(ent.raw, sql.FieldRaw(f.fieldType, ent.value[i])...)
		}
		entries = append(entries, ent)
	}

	// 遍历插入的数据条目，写入数据
	id := uint64(0)
	for _, ent := range entries {
		// 写入数据
		id, err = tbm.verManage.Insert(txId, ent.raw)
		if err != nil {
			return err
		}

		// 判断是否有字段需要索引
		for i, f := range ent.fields {
			if f.isIndex() {
				key := utils.Hash(ent.value[i])
				err = f.idx.Insert(key, id)
				if err != nil {
					return err
				}
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

func (tbm *tableManage) Select(txId uint64, stmt *sql.SelectStmt) ([]byte, error) {
	// 获取表对象
	_, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return nil, ErrNoSuchTable
	}

	// 遍历条件，如果有索引，则使用索引进行查询
	if len(stmt.Where) == 0 {
		// 根据主键索引，扫描所有数据

	} else {
		for _, cond := range stmt.Where {
			fmt.Printf("cond: %+v\n", cond)
		}
	}
	return nil, nil
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
	thead := []string{"Field", "Type", "Index"}
	tbody := make([][]string, 0)
	t, exist := tbm.tables.Get(table)
	if !exist {
		return "no such table"
	}

	for _, f := range t.tableFields {
		index := "NO"
		if f.isIndex() {
			index = "YES"
		}
		tbody = append(tbody, []string{
			f.fieldName,
			f.fieldType,
			index,
		})
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
