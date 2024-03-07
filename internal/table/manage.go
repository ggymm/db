package table

import (
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
		tbm.tables.Set(t.name, t)

		// 读取下一个表的信息
		id = t.nextId
	}
	return tbm
}

func (tbm *tableManage) readTableId() uint64 {
	return bin.Uint64(tbm.boot.Load())
}

func (tbm *tableManage) updateTableId(id uint64) {
	buf := make([]byte, 8)
	bin.PutUint64(buf, id)
	tbm.boot.Update(buf)
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
		TxId:   txId,
		NextId: tbm.readTableId(),
		Stmt:   stmt,
	})
	if err != nil {
		return err
	}

	tbm.tables.Set(t.name, t)
	tbm.updateTableId(t.id)
	return nil
}

func (tbm *tableManage) Insert(txId uint64, stmt *sql.InsertStmt) error {
	var (
		ok  bool
		err error

		t    *table
		maps []map[string]string
	)

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
	es := make([]entry, len(t.fields))
	for _, ins := range maps {
		e := entry{
			raw:    make([]byte, 0),
			value:  make([]any, len(t.fields)),
			fields: make([]*field, len(t.fields)),
		}

		val := ""
		for i, f := range t.fields {
			e.fields[i] = f

			// 获取字段值
			val, ok = ins[f.name]
			switch {
			case ok:
				e.value[i] = val
			case len(f.defaultVal) != 0:
				e.value[i] = f.defaultVal
			case f.allowNull:
				e.value[i] = nil
			default:
				return fmt.Errorf("field %s is not allowed to be null", f.name)
			}

			// 获取字段二进制值
			e.raw = append(e.raw, sql.FieldRaw(f.dataType, e.value[i])...)
		}
		es = append(es, e)
	}

	// 遍历插入的数据条目，写入数据
	id := uint64(0)
	for _, e := range es {
		// 写入数据
		id, err = tbm.verManage.Insert(txId, e.raw)
		if err != nil {
			return err
		}

		// 判断是否有字段需要索引
		for _, f := range e.fields {
			if f.isIndexed() {
				err = f.idx.Insert(utils.Hash(e.value[f.index]), id)
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
	_, ok := tbm.tables.Get(stmt.Table)
	if !ok {
		return nil, ErrNoSuchTable
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
	v := newView()
	v.setHead(thead)
	v.setBody(tbody)
	return v.string(singleLine)
}

func (tbm *tableManage) ShowField(table string) string {
	thead := []string{"Field", "Type", "Index"}
	tbody := make([][]string, 0)
	t, exist := tbm.tables.Get(table)
	if !exist {
		return "no such table"
	}

	for _, f := range t.fields {
		index := "NO"
		if f.index != 0 {
			index = "YES"
		}
		tbody = append(tbody, []string{
			f.name,
			f.dataType,
			index,
		})
	}

	// 表格形式输出
	v := newView()
	v.setHead(thead)
	v.setBody(tbody)
	return v.string(singleLine)
}

func (tbm *tableManage) VerManage() ver.Manage {
	return tbm.verManage
}

func (tbm *tableManage) DataManage() data.Manage {
	return tbm.dataManage
}
