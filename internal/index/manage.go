package index

import (
	"db/internal/data"
	"sync"
)

type Manage interface {
	Insert(key, nodeId uint64) error
	Search(key uint64) ([]uint64, error)
	SearchRange(leftKey, rightKey uint64) ([]uint64, error)
}

type indexManage struct {
	lock     sync.Mutex
	bootId   uint64
	bootItem data.Item

	dataManage data.Manage
}

func (i *indexManage) Insert(key, nodeId uint64) error {
	panic("implement me")
}

func (i *indexManage) Search(key uint64) ([]uint64, error) {
	panic("implement me")
}

func (i *indexManage) SearchRange(leftKey, rightKey uint64) ([]uint64, error) {
	panic("implement me")
}
