package ver

import (
	"db/internal/ver/lock"
	"errors"
	"sync"

	"db/internal/data"
	"db/internal/tx"
	"db/pkg/cache"
)

// 版本管理器
//
// 抽象

var (
	ErrNotFound = errors.New("entry not found")
)

type Manage interface {
	Read(tid uint64, itemId uint64)
	Delete(tid uint64, itemId uint64)
	Insert(tid uint64, data []byte)

	Begin(level int) uint64
	Abort(tid uint64)
	Commit(tid uint64) error
}

type verManage struct {
	lock sync.Mutex

	txManage   tx.Manage
	dataManage data.Manage

	vLock lock.Lock
	cache cache.Cache
}

func NewManage(tm tx.Manage, dm data.Manage) Manage {
	vm := new(verManage)

	vm.txManage = tm
	vm.dataManage = dm

	vm.cache = cache.NewCache(&cache.Option{
		Obtain:   vm.obtainForCache,
		Release:  vm.releaseForCache,
		MaxCount: 0,
	})
	return vm
}

// obtainForCache 需要支持并发
// 缓存中不存在时，从 data_manage 中获取
func (vm *verManage) obtainForCache(key uint64) (any, error) {
	item, ok, err := vm.dataManage.Read(key)
	if err != nil {
		return nil, err
	}
	if ok == false {
		return nil, ErrNotFound
	}

	return &entry{id: key, item: item, manage: vm}, nil
}

// releaseForCache 需要是同步方法
// 释放缓存，需要将 Page 对象内存刷新到磁盘
func (vm *verManage) releaseForCache(data any) {
	ent := data.(*entry)
	ent.item.Release() // 将 entry 从内存中彻底释放
}

func (vm *verManage) Read(tid uint64, itemId uint64) {
	//TODO implement me
	panic("implement me")
}

func (vm *verManage) Delete(tid uint64, itemId uint64) {
	//TODO implement me
	panic("implement me")
}

func (vm *verManage) Insert(tid uint64, data []byte) {
	//TODO implement me
	panic("implement me")
}

func (vm *verManage) Begin(level int) uint64 {
	//TODO implement me
	panic("implement me")
}

func (vm *verManage) Abort(tid uint64) {
	//TODO implement me
	panic("implement me")
}

func (vm *verManage) Commit(tid uint64) error {
	//TODO implement me
	panic("implement me")
}
