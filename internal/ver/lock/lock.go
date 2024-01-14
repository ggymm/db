package lock

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrCheckDeadlock = errors.New("check deadlock error")
)

type Lock interface {
	Add(tid, itemId uint64) (bool, chan struct{})
	Remove(tid uint64)
}

type lock struct {
	lock sync.Mutex

	waitCh map[uint64]chan struct{}

	tx2item     map[uint64]*list.List // 事务 已经获取的 数据
	item2tx     map[uint64]uint64     // 数据 被哪个事务 获取
	waitTx2Item map[uint64]uint64     // 事务 在等待哪个 数据
	waitItem2Tx map[uint64]*list.List // 数据 被哪些事务 等待
}

func NewLock() Lock {
	return &lock{
		waitCh: make(map[uint64]chan struct{}),

		tx2item: make(map[uint64]*list.List),
		item2tx: make(map[uint64]uint64),

		waitTx2Item: make(map[uint64]uint64),
		waitItem2Tx: make(map[uint64]*list.List),
	}
}

func inList(data map[uint64]*list.List, key, val uint64) bool {
	if values, ok := data[key]; ok {
		e := values.Front()
		for e != nil {
			if e.Value.(uint64) == val {
				return true
			}
			e = e.Next()
		}
	}
	return false
}

func putList(data map[uint64]*list.List, key, val uint64) {
	if _, ok := data[key]; !ok {
		data[key] = new(list.List)
	}
	data[key].PushFront(val)
}

func removeList(data map[uint64]*list.List, key, val uint64) {
	if values, ok := data[key]; ok {
		e := values.Front()
		for e != nil {
			if e.Value.(uint64) == val {
				values.Remove(e)
				break
			}
		}
		if values.Len() == 0 {
			delete(data, key)
		}
	}
}

var (
	stamp    int
	tidStamp map[uint64]int
)

// dfs 遍历依赖关系
// 从 item2tx 和 waitTx2Item 取出的边
// 是否会构成环，即
// 是否会出现重复的 tid（通过 tidStamp 中 tid 的取值判断）
func (l *lock) dfs(tid uint64) bool {
	stp := tidStamp[tid]
	if stp != 0 {
		if stp == stamp {
			return true
		} else if stp < stamp {
			return false
		}
	}
	tidStamp[tid] = stamp

	itemId := l.waitTx2Item[tid]
	if itemId == 0 {
		return false
	}
	tid = l.item2tx[itemId]
	if tid == 0 {
		panic(ErrCheckDeadlock)
	}
	return l.dfs(tid)
}

func (l *lock) checkDeadlock() bool {
	stamp = 1
	tidStamp = make(map[uint64]int)
	for tid, _ := range l.tx2item {
		if tidStamp[tid] > 0 {
			continue
		}
		stamp++
		if l.dfs(tid) {
			return true
		}
	}
	return false
}

func (l *lock) selectNextTID(itemId uint64) {
	delete(l.item2tx, itemId)
	txList := l.waitItem2Tx[itemId]
	if txList == nil {
		return
	}
	for txList.Len() > 0 {
		e := txList.Front()
		v := txList.Remove(e)
		tid := v.(uint64)
		if _, ok := l.waitCh[tid]; !ok {
			continue
		} else {
			l.item2tx[itemId] = tid
			ch := l.waitCh[tid]
			delete(l.waitCh, tid)
			delete(l.waitTx2Item, tid)
			ch <- struct{}{}
			break
		}
	}
	if txList.Len() == 0 {
		delete(l.waitItem2Tx, itemId)
	}
}

func (l *lock) Add(tid, itemId uint64) (bool, chan struct{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// 判断 itemId 在 tx2item 是否存在
	if inList(l.tx2item, tid, itemId) {
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		return true, ch
	}

	// 判断 itemId 是否被 其他 tid 占用
	if _, ok := l.item2tx[itemId]; !ok {
		l.item2tx[itemId] = tid
		putList(l.tx2item, tid, itemId)
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		return true, ch
	}

	// 添加 tid -> itemId 的等待，并判断是否会死锁
	l.waitTx2Item[tid] = itemId
	putList(l.waitItem2Tx, itemId, tid) // 添加到等待列表
	if l.checkDeadlock() {
		delete(l.waitTx2Item, tid)
		removeList(l.waitItem2Tx, itemId, tid)
		return false, nil
	}

	// 无死锁，则等待
	ch := make(chan struct{})
	l.waitCh[tid] = ch
	return true, ch
}

func (l *lock) Remove(tid uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()

	vs := l.tx2item[tid]
	if vs != nil {
		for vs.Len() > 0 {
			e := vs.Front()
			v := vs.Remove(e)
			l.selectNextTID(v.(uint64))
		}
	}

	delete(l.waitCh, tid)
	delete(l.tx2item, tid)
	delete(l.waitTx2Item, tid)
}
