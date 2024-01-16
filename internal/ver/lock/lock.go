package lock

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"sort"
	"sync"
)

var (
	ErrCheckDeadlock = errors.New("check deadlock error")
)

type Lock interface {
	Add(id, key uint64) (bool, chan struct{})
	Remove(id uint64)

	String() string
}

type lock struct {
	lock sync.Mutex

	keyId     map[uint64]uint64 // 数据 被哪个事务 获取
	waitIdKey map[uint64]uint64 // 事务 在等待哪个 数据

	keys    map[uint64]*list.List // 事务 已经获取的 数据
	waitIds map[uint64]*list.List // 数据 被哪些事务 等待

	waitCh map[uint64]chan struct{} // 事务 等待通道
}

func NewLock() Lock {
	return &lock{
		keyId:     make(map[uint64]uint64),
		waitIdKey: make(map[uint64]uint64),

		keys:    make(map[uint64]*list.List),
		waitIds: make(map[uint64]*list.List),

		waitCh: make(map[uint64]chan struct{}),
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
	stamp   int
	idStamp map[uint64]int
)

// dfs 遍历依赖关系
// 从 keyId 和 waitIdKey 取出的边
// 是否会构成环，即
// 是否会出现重复的 id（通过 idStamp 中 id 的取值判断）
func (l *lock) dfs(id uint64) bool {
	stp := idStamp[id]
	if stp != 0 {
		if stp == stamp {
			return true
		} else if stp < stamp {
			return false
		}
	}
	idStamp[id] = stamp

	key := l.waitIdKey[id]
	if key == 0 {
		return false
	}
	id = l.keyId[key]
	if id == 0 {
		panic(ErrCheckDeadlock)
	}
	return l.dfs(id)
}

func (l *lock) checkDeadlock() bool {
	stamp = 1
	idStamp = make(map[uint64]int)
	for id := range l.keys {
		if idStamp[id] > 0 {
			continue
		}
		stamp++
		if l.dfs(id) {
			return true
		}
	}
	return false
}

func (l *lock) selectNextTID(key uint64) {
	delete(l.keyId, key)
	txList := l.waitIds[key]
	if txList == nil {
		return
	}
	for txList.Len() > 0 {
		e := txList.Front()
		v := txList.Remove(e)
		id := v.(uint64)
		if _, ok := l.waitCh[id]; !ok {
			continue
		} else {
			l.keyId[key] = id
			ch := l.waitCh[id]
			delete(l.waitCh, id)
			delete(l.waitIdKey, id)
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

	// 判断 key 是否存在
	if inList(l.keys, id, key) {
		// 如果已经获取锁，直接返回成功
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		return true, ch
	}

	// 判断 key 是否被 其他 id 占用
	if _, ok := l.keyId[key]; !ok {
		// 如果没有被占用
		l.keyId[key] = id        // 添加 key -> id 的依赖关系
		putList(l.keys, id, key) // 保存 key 到 id 对应的 keys 中
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		return true, ch
	}

	l.waitIdKey[id] = key       // 添加 id -> key 的等待关系
	putList(l.waitIds, key, id) // 保存 id 到 key 对应的 waitIds 中
	// 判断是否会出现死锁
	if l.checkDeadlock() {
		// 如果出现死锁
		delete(l.waitIdKey, id)        // 删除 id -> key 的等待关系
		removeList(l.waitIds, key, id) // 删除 id 从 key 对应的 waitIds 中
		return false, nil              // 获取锁失败
	}

	// 没有死锁，创建等待通道，等待释放锁
	ch := make(chan struct{})
	l.waitCh[id] = ch
	return true, ch
}

func (l *lock) Remove(id uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()

	vs := l.keys[id]
	if vs != nil {
		for vs.Len() > 0 {
			e := vs.Front()
			v := vs.Remove(e)
			l.selectNextTID(v.(uint64))
		}
	}

	delete(l.waitCh, id)
	delete(l.keys, id)
	delete(l.waitIdKey, id)
}

func (l *lock) String() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	var mapIntString = func(data map[uint64]uint64) string {
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

	var mapListString = func(data map[uint64]*list.List) string {
		if len(data) == 0 {
			return ""
		}
		var buf bytes.Buffer
		for key, values := range data {
			buf.WriteString(fmt.Sprintf("%d: [", key))
			e := values.Front()
			for e != nil {
				buf.WriteString(fmt.Sprintf("%d", e.Value.(uint64)))
				e = e.Next()
				if e != nil {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]\n")
		}
		return buf.String()
	}

	var buf bytes.Buffer
	buf.WriteString("waitCh(事务 等待通道):\n")
	for key := range l.waitCh {
		buf.WriteString(fmt.Sprintf("%d\n", key))
	}

	buf.WriteString("keyId(数据 被哪个事务 获取):\n")
	buf.WriteString(mapIntString(l.keyId))
	buf.WriteString("waitIdKey(事务 在等待哪个 数据):\n")
	buf.WriteString(mapIntString(l.waitIdKey))

	buf.WriteString("keys(事务 已经获取的 数据):\n")
	buf.WriteString(mapListString(l.keys))
	buf.WriteString("waitIds(数据 被哪些事务 等待):\n")
	buf.WriteString(mapListString(l.waitIds))
	return buf.String()
}
