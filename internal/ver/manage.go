package ver

import (
	"errors"
	"sync"

	"db/internal/data"
	"db/internal/tx"
	"db/internal/ver/lock"
	"db/pkg/cache"
)

// 版本管理器
//
// 抽象

var (
	ErrNotFound     = errors.New("entry not found")
	ErrCannotHandle = errors.New("could not access due to concurrent update")
)

type Manage interface {
	Read(tid uint64, key uint64) ([]byte, bool, error)
	Insert(tid uint64, data []byte) (uint64, error)
	Delete(tid uint64, key uint64) (bool, error)

	Begin(level int) uint64
	Abort(tid uint64)
	Commit(tid uint64) error
}

type verManage struct {
	lock sync.Mutex

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

	vm.cache = cache.NewCache(&cache.Option{
		Obtain:   vm.obtainForCache,
		Release:  vm.releaseForCache,
		MaxCount: 0,
	})
	return vm
}

// 撤销事务
//
// 手动撤销：
// 自动撤销：
func (vm *verManage) abort(tid uint64, manual bool) {
	vm.lock.Lock()
	t := vm.txCache[tid]
	if manual {
		delete(vm.txCache, tid)
	}
	vm.lock.Unlock()

	if t.AutoAborted {
		return
	}

	vm.txLock.Remove(tid)
	vm.txManage.Abort(tid)
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

func (vm *verManage) Read(tid uint64, key uint64) ([]byte, bool, error) {
	vm.lock.Lock()
	t := vm.txCache[tid]
	vm.lock.Unlock()

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

func (vm *verManage) Insert(tid uint64, data []byte) (uint64, error) {
	vm.lock.Lock()
	t := vm.txCache[tid]
	vm.lock.Unlock()

	if t.Err != nil {
		return 0, t.Err
	}

	// 包装成 entry 数据
	ent := make([]byte, offData+len(data))
	tx.WriteId(ent[offMin:], tid)
	copy(ent[offData:], data)
	return vm.dataManage.Insert(tid, ent)
}

func (vm *verManage) Delete(tid uint64, key uint64) (bool, error) {
	vm.lock.Lock()
	t := vm.txCache[tid]
	vm.lock.Unlock()

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
		vm.abort(tid, false) // 自动取消
		t.Err = ErrCannotHandle
		t.AutoAborted = true
		return false, t.Err
	}
	<-ch // 等待锁释放

	// 判断是否被自己删除
	if ent.Max() == tid {
		return false, nil
	}

	// 判断是否发生了版本跳跃
	if t.IsSkip(vm.txManage, ent) {
		vm.abort(tid, false) // 自动取消
		t.Err = ErrCannotHandle
		t.AutoAborted = true
		return false, t.Err
	}

	// 更新 max
	ent.SetMax(tid) // 执行删除操作
	return true, nil
}

// Begin 开启一个事务
//
// 保存当前处于激活状态的事务
func (vm *verManage) Begin(level int) uint64 {
	vm.lock.Lock()
	defer vm.lock.Unlock()

	// 开启一个事务，并且缓存当前处于激活状态该的事务
	tid := vm.txManage.Begin()
	vm.txCache[tid] = newTransaction(tid, level, vm.txCache)
	return tid
}

// Abort 取消一个事务
func (vm *verManage) Abort(tid uint64) {
	vm.abort(tid, true)
}

// Commit 提交一个事务
func (vm *verManage) Commit(tid uint64) error {
	vm.lock.Lock()
	t := vm.txCache[tid]
	vm.lock.Unlock()

	if t.Err != nil {
		return t.Err
	}

	vm.lock.Lock()
	delete(vm.txCache, tid)
	vm.lock.Unlock()

	vm.txLock.Remove(tid)
	vm.txManage.Commit(tid)
	return nil
}
