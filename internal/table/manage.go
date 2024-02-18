package table

import (
	"db/internal/boot"
	"db/internal/data"
	"db/internal/ver"
	"db/pkg/sql"
	"sync"
)

type Manage interface {
	Begin(level int) uint64
	Abort(id uint64)
	Commit(id uint64) error

	Show() ([]byte, error)
	Create(id uint64, stmt sql.CreateStmt) error

	Insert(id uint64, stmt sql.InsertStmt) error
	Update(id uint64, stmt sql.UpdateStmt) error
	Delete(id uint64, stmt sql.DeleteStmt) error
	Select(id uint64, stmt sql.SelectStmt) ([]byte, error)
}

type tableManage struct {
	boot       *boot.Boot
	verManage  ver.Manage
	dataManage data.Manage

	lock sync.Mutex

	tables map[string]*Table
}

func NewManage(boot *boot.Boot, verManage ver.Manage, dataManage data.Manage) Manage {
	return &tableManage{
		boot:       boot,
		verManage:  verManage,
		dataManage: dataManage,
		tables:     make(map[string]*Table),
	}
}

func (t *tableManage) Begin(level int) uint64 {
	return t.verManage.Begin(level)
}

func (t *tableManage) Abort(id uint64) {
	t.verManage.Abort(id)
}

func (t *tableManage) Commit(id uint64) error {
	return t.verManage.Commit(id)
}

func (t *tableManage) Show() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableManage) Create(id uint64, stmt sql.CreateStmt) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableManage) Insert(id uint64, stmt sql.InsertStmt) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableManage) Update(id uint64, stmt sql.UpdateStmt) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableManage) Delete(id uint64, stmt sql.DeleteStmt) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableManage) Select(id uint64, stmt sql.SelectStmt) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
