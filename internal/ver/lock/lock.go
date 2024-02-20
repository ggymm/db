package lock

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"sort"
	"sync"
)

var ErrCheckDeadlock = errors.New("check deadlock error")

type Lock interface {
	Add(id, key uint64) (bool, chan struct{})
	Remove(id uint64)

	String() string
}

type lock struct {
	lock sync.Mutex

	keyOwn   map[uint64]uint64        // 数据 被哪个事务 获取
	waitKey  map[uint64]uint64        // 事务 在等待哪个 数据
	waitIds  map[uint64]*list.List    // 数据 被哪些事务 等待（用于数据解除占用后，选择下一个事务）
	currKeys map[uint64]*list.List    // 事务 已经获取的 数据
	waitLock map[uint64]chan struct{} // 事务 等待通道
}

func NewLock() Lock {
	return &lock{
		keyOwn:   make(map[uint64]uint64),
		waitKey:  make(map[uint64]uint64),
		waitIds:  make(map[uint64]*list.List),
		currKeys: make(map[uint64]*list.List),
		waitLock: make(map[uint64]chan struct{}),
	}
}

func hasItem(data map[uint64]*list.List, key, val uint64) bool {
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

func putItem(data map[uint64]*list.List, key, val uint64) {
	if _, ok := data[key]; !ok {
		data[key] = new(list.List)
	}
	data[key].PushFront(val)
}

func removeItem(data map[uint64]*list.List, key, val uint64) {
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
	state   int
	idState map[uint64]int
)

// dfs 遍历依赖关系
// 从 keyOwn 和 waitKey 取出的边
// 是否会构成环，即
// 是否会出现重复的 id（通过 idState 中 id 的取值判断）
func (l *lock) dfs(id uint64) bool {
	if id == 0 {
		panic(ErrCheckDeadlock)
	}
	if st := idState[id]; st != 0 {
		if st == state {
			return true
		} else if st < state {
			return false
		}
	}
	idState[id] = state

	// 继续遍历
	// 遍历属于同一个依赖关系的下一个节点 id
	// 首先获取当前 id 所等待的 key
	// 其次获取当前 key 所属的 id （即当前 id 所等待的 id）
	key := l.waitKey[id]
	if key == 0 {
		return false
	}
	return l.dfs(l.keyOwn[key]) // 继续遍历 key 所属的 id
}

func (l *lock) checkDeadlock() bool {
	state = 1
	idState = make(map[uint64]int)
	for id := range l.currKeys {
		if idState[id] > 0 {
			continue
		}
		state++
		if l.dfs(id) {
			return true
		}
	}
	return false
}

func (l *lock) selectNextId(key uint64) {
	delete(l.keyOwn, key)
	txList := l.waitIds[key]
	if txList == nil {
		return
	}
	for txList.Len() > 0 {
		e := txList.Front()
		v := txList.Remove(e)
		id := v.(uint64)
		if _, ok := l.waitLock[id]; !ok {
			continue
		} else {
			l.keyOwn[key] = id
			ch := l.waitLock[id]
			delete(l.waitLock, id)
			delete(l.waitKey, id)
			ch <- struct{}{}
			break
		}
	}
	if txList.Len() == 0 {
		delete(l.waitIds, key)
	}
}

// Add 添加 id -> key 的依赖关系
//
// id: 事务 id
// key: 数据 key
//
// 返回值:
// bool: 是否添加成功
// chan struct{}: 事务等待通道
func (l *lock) Add(id, key uint64) (bool, chan struct{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	success := func() (bool, chan struct{}) {
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		return true, ch
	}

	// 判断 key 是否已经被当前 id 占用
	if hasItem(l.currKeys, id, key) {
		return success()
	}

	// 判断 key 是否被 其他 id 占用
	if _, ok := l.keyOwn[key]; !ok {
		// 如果没有被占用
		l.keyOwn[key] = id           // 添加 key -> id 的依赖关系
		putItem(l.currKeys, id, key) // 添加 id -> key 的占用关系
		return success()
	}

	l.waitKey[id] = key         // 添加 id -> key 的等待关系
	putItem(l.waitIds, key, id) // 添加 key -> id 的等待关系
	// 判断是否会出现死锁
	if l.checkDeadlock() {
		// 如果出现死锁
		delete(l.waitKey, id)          // 删除 id -> key 的等待关系
		removeItem(l.waitIds, key, id) // 删除 key -> id 的等待关系
		return false, nil              // 获取锁失败
	}

	// 没有死锁，创建等待通道，等待释放锁
	ch := make(chan struct{})
	l.waitLock[id] = ch // 添加 id -> 等待通道 的关系
	return true, ch
}

func (l *lock) Remove(id uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()

	keys := l.currKeys[id]
	if keys != nil {
		for keys.Len() > 0 {
			ele := keys.Front()
			key := keys.Remove(ele)
			// 释放一个 key，选择下一个事务
			l.selectNextId(key.(uint64))
		}
	}

	delete(l.waitKey, id)
	delete(l.currKeys, id)
	delete(l.waitLock, id)
}

func (l *lock) String() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	mapIntString := func(data map[uint64]uint64) string {
		if len(data) == 0 {
			return ""
		}
		var keys []uint64
		for key := range data {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		var buf bytes.Buffer
		for _, key := range keys {
			buf.WriteString(fmt.Sprintf("%d: %d\n", key, data[key]))
		}
		return buf.String()
	}

	mapListString := func(data map[uint64]*list.List) string {
		if len(data) == 0 {
			return ""
		}
		var keys []uint64
		for key := range data {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		var buf bytes.Buffer
		for _, key := range keys {
			buf.WriteString(fmt.Sprintf("%d: [", key))
			item := data[key].Front()
			for item != nil {
				buf.WriteString(fmt.Sprintf("%d", item.Value.(uint64)))
				item = item.Next()
				if item != nil {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]\n")
		}
		return buf.String()
	}

	var buf bytes.Buffer
	buf.WriteString("waitLock(事务 等待通道):\n")
	for key := range l.waitLock {
		buf.WriteString(fmt.Sprintf("%d\n", key))
	}

	buf.WriteString("currKeys(事务 已经获取的 数据):\n")
	buf.WriteString(mapListString(l.currKeys))
	buf.WriteString("keyOwn(数据 被哪个事务 获取):\n")
	buf.WriteString(mapIntString(l.keyOwn))
	buf.WriteString("waitKey(事务 在等待哪个 数据):\n")
	buf.WriteString(mapIntString(l.waitKey))
	buf.WriteString("waitIds(数据 被哪些事务 等待):\n")
	buf.WriteString(mapListString(l.waitIds))
	return buf.String()
}
