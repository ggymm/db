package ver

import (
	"errors"
	"sync"

	"github.com/ggymm/db/data"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/cache"
	"github.com/ggymm/db/tx"
	"github.com/ggymm/db/ver/lock"
)

// 版本管理器
//
// 抽象

var (
	ErrNotFound     = errors.New("entry not found")
	ErrCannotHandle = errors.New("could not access due to concurrent update")
)

type Manage interface {
	Begin(level int) uint64
	Commit(tid uint64) error
	Rollback(tid uint64)

	Read(tid uint64, key uint64) ([]byte, bool, error)
	Write(tid uint64, data []byte) (uint64, error)
	Delete(tid uint64, key uint64) (bool, error)
}

type verManage struct {
	mu sync.Mutex

	txManage   tx.Manage
	dataManage data.Manage

	txLock  lock.Lock
	txCache map[uint64]*transaction

	cache cache.Cache
}

func NewManage(tm tx.Manage, dm data.Manage) Manage {
	vm := new(verManage)

	vm.txManage = tm
	vm.dataManage = dm

	vm.txLock = lock.NewLock()
	vm.txCache = make(map[uint64]*transaction)
	vm.txCache[tx.Super] = newTransaction(tx.Super, 0, nil)

	vm.cache = cache.NewCache(&cache.Option{
		Obtain:   vm.obtainForCache,
		Release:  vm.releaseForCache,
		MaxCount: 0,
	})
	return vm
}

// 回滚事务
//
// 手动回滚：
// 自动回滚：
func (vm *verManage) rollback(tid uint64, manual bool) {
	vm.mu.Lock()
	t := vm.txCache[tid]
	if manual {
		delete(vm.txCache, tid)
	}
	vm.mu.Unlock()

	if t.AutoRollback {
		return
	}

	vm.txLock.Remove(tid)
	vm.txManage.Rollback(tid)
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

// Begin 开启一个事务
//
// 保存当前处于激活状态的事务
func (vm *verManage) Begin(level int) uint64 {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// 开启一个事务，并且缓存当前处于激活状态该的事务
	tid := vm.txManage.Begin()
	vm.txCache[tid] = newTransaction(tid, level, vm.txCache)
	return tid
}

// Commit 提交一个事务
func (vm *verManage) Commit(tid uint64) error {
	vm.mu.Lock()
	t := vm.txCache[tid]
	vm.mu.Unlock()

	if t.Err != nil {
		return t.Err
	}

	vm.mu.Lock()
	delete(vm.txCache, tid)
	vm.mu.Unlock()

	vm.txLock.Remove(tid)
	vm.txManage.Commit(tid)
	return nil
}

// Rollback 回滚一个事务
func (vm *verManage) Rollback(tid uint64) {
	vm.rollback(tid, true)
}

func (vm *verManage) Read(tid uint64, key uint64) ([]byte, bool, error) {
	vm.mu.Lock()
	t := vm.txCache[tid]
	vm.mu.Unlock()

	if t.Err != nil {
		return nil, false, t.Err
	}

	// 读取数据
	val, err := vm.cache.Obtain(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	ent := val.(*entry)
	defer vm.cache.Release(tid) // 释放缓存

	if !t.IsVisible(vm.txManage, ent) {
		return nil, false, nil
	}
	return ent.Data(), true, nil
}

func (vm *verManage) Write(tid uint64, data []byte) (uint64, error) {
	vm.mu.Lock()
	t := vm.txCache[tid]
	vm.mu.Unlock()

	if t.Err != nil {
		return 0, t.Err
	}

	// 包装成 entry 数据
	ent := make([]byte, offData+len(data))
	bin.PutUint64(ent[offMin:], tid)
	copy(ent[offData:], data)
	return vm.dataManage.Write(tid, ent)
}

func (vm *verManage) Delete(tid uint64, key uint64) (bool, error) {
	vm.mu.Lock()
	t := vm.txCache[tid]
	vm.mu.Unlock()

	if t.Err != nil {
		return false, t.Err
	}

	// 获取数据
	val, err := vm.cache.Obtain(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	ent := val.(*entry)
	defer vm.cache.Release(tid) // 释放缓存

	// 判断是否可见
	if !t.IsVisible(vm.txManage, ent) {
		return false, nil
	}

	// 添加锁并判断是否死锁
	ok, ch := vm.txLock.Add(tid, key)
	if !ok {
		vm.rollback(tid, false) // 自动回滚
		t.Err = ErrCannotHandle
		t.AutoRollback = true
		return false, t.Err
	}
	<-ch // 等待锁释放

	// 判断是否被自己删除
	if ent.Max() == tid {
		return false, nil
	}

	// 判断是否发生了版本跳跃
	if t.IsSkip(vm.txManage, ent) {
		vm.rollback(tid, false) // 自动回滚
		t.Err = ErrCannotHandle
		t.AutoRollback = true
		return false, t.Err
	}

	// 更新 max
	ent.SetMax(tid) // 执行删除操作
	return true, nil
}
