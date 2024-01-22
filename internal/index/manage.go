package index

import (
	"sync"

	"db/internal/data"
)

type Index interface {
	Insert(key, nodeId uint64) error
	Search(key uint64) ([]uint64, error)
	SearchRange(leftKey, rightKey uint64) ([]uint64, error)
}

type index struct {
	lock     sync.Mutex
	bootId   uint64
	bootItem data.Item

	dataManage data.Manage
}

func (i *index) Insert(key, nodeId uint64) error {
	panic("implement me")
}

func (i *index) Search(key uint64) ([]uint64, error) {
	panic("implement me")
}

func (i *index) SearchRange(leftKey, rightKey uint64) ([]uint64, error) {
	panic("implement me")
}
