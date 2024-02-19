package table

import (
	"db/internal/boot"
	"db/internal/data"
	"db/internal/ver"
	"db/pkg/bin"
	"db/pkg/sql"
	"sync"
)

type Manage interface {
	Begin(level int) uint64
	Abort(txId uint64)
	Commit(txId uint64) error

	Show() ([]byte, error)
	Create(txId uint64, stmt *sql.CreateStmt) error

	Insert(txId uint64, stmt *sql.InsertStmt) error
	Update(txId uint64, stmt *sql.UpdateStmt) error
	Delete(txId uint64, stmt *sql.DeleteStmt) error
	Select(txId uint64, stmt *sql.SelectStmt) ([]byte, error)

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

func (tbm *tableManage) writeTableId(id uint64) {
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

func (tbm *tableManage) Show() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Create(txId uint64, stmt *sql.CreateStmt) error {
	tbm.lock.Lock()
	defer tbm.lock.Unlock()

	if _, exist := tbm.tables[stmt.Name]; exist {
		return nil
	}

	t, err := createTable(tbm, &newTable{TxId: txId, NextId: tbm.readTableId(), Stmt: stmt})
	if err != nil {
		return err
	}

	tbm.writeTableId(t.Id)
	tbm.tables[t.Name] = t
	return nil
}

func (tbm *tableManage) Insert(txId uint64, stmt *sql.InsertStmt) error {
	//TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Update(txId uint64, stmt *sql.UpdateStmt) error {
	//TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Delete(txId uint64, stmt *sql.DeleteStmt) error {
	//TODO implement me
	panic("implement me")
}

func (tbm *tableManage) Select(txId uint64, stmt *sql.SelectStmt) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (tbm *tableManage) VerManage() ver.Manage {
	return tbm.verManage
}

func (tbm *tableManage) DataManage() data.Manage {
	return tbm.dataManage
}
