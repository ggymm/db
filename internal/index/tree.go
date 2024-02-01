package index

import (
	"sync"

	"db/internal/data"
)

type Index interface {
	Insert(key, itemId uint64) error
	Search(key uint64) ([]uint64, error)
	SearchRange(leftKey, rightKey uint64) ([]uint64, error)
}

type index struct {
	lock     sync.Mutex
	bootId   uint64
	bootItem data.Item

	DataManage data.Manage
}

// Insert
// 插入 key（字段计算的hash值） 和 itemId（数据项的Id） 的索引关系
func (i *index) Insert(key, itemId uint64) error {
	panic("implement me")
}

func (i *index) Search(key uint64) ([]uint64, error) {
	panic("implement me")
}

func (i *index) SearchRange(leftKey, rightKey uint64) ([]uint64, error) {
	panic("implement me")
}
