package table

import (
	"sync"

	"db/internal/boot"
	"db/internal/data"
	"db/internal/ver"
	"db/pkg/bin"
	"db/pkg/sql"
)

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

	lock sync.Mutex

	tables map[string]*table
}

func NewManage(boot boot.Boot, verManage ver.Manage, dataManage data.Manage) Manage {
	tbm := &tableManage{
		boot:       boot,
		verManage:  verManage,
		dataManage: dataManage,
		tables:     make(map[string]*table),
	}

	id := tbm.readTableId()
	for id != 0 {
		t := readTable(tbm, id)
		tbm.tables[t.Name] = t

		// 读取下一个表的信息
		id = t.NextId
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

func (tbm *tableManage) Show() []byte {
	tbm.lock.Lock()
	defer tbm.lock.Unlock()

	result := ""
	printTable := func(t *table) string {
		str := "Table: " + t.Name + "\n"
		return str
	}

	for _, t := range tbm.tables {
		result += printTable(t)
	}
	return []byte(result)
}

func (tbm *tableManage) Create(txId uint64, stmt *sql.CreateStmt) error {
	tbm.lock.Lock()
	defer tbm.lock.Unlock()

	if _, exist := tbm.tables[stmt.Name]; exist {
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

	tbm.tables[t.Name] = t
	tbm.updateTableId(t.Id)
	return nil
}

func (tbm *tableManage) Insert(txId uint64, stmt *sql.InsertStmt) error {
	// TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Update(txId uint64, stmt *sql.UpdateStmt) error {
	// TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Delete(txId uint64, stmt *sql.DeleteStmt) error {
	// TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Select(txId uint64, stmt *sql.SelectStmt) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (tbm *tableManage) ShowTable() string {
	thead := []string{"Tables"}
	tbody := make([][]string, 0)
	for name := range tbm.tables {
		tbody = append(tbody, []string{name})
	}

	// 表格形式输出
	v := newView()
	v.setHead(thead)
	v.setBody(tbody)
	return v.string(singleLine)
}

func (tbm *tableManage) ShowField(table string) string {
	thead := []string{"Field", "Type", "Index"}
	tbody := make([][]string, 0)
	t, exist := tbm.tables[table]
	if !exist {
		return "no such table"
	}

	for _, f := range t.Fields {
		index := "NO"
		if f.Index != 0 {
			index = "YES"
		}
		tbody = append(tbody, []string{f.Name, f.Type, index})
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
