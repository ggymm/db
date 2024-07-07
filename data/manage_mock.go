package data

import (
	"math/rand"
	"sync"

	"db/tx"
)

type mockManage struct {
	lock  sync.Mutex
	cache map[uint64]*mockItem
}

func newMockManage() Manage {
	return &mockManage{
		cache: make(map[uint64]*mockItem),
	}
}

func (m *mockManage) Close() {
}

func (m *mockManage) Read(id uint64) (Item, bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.cache[id]; ok == false {
		return nil, false, nil
	}
	return m.cache[id], true, nil
}

func (m *mockManage) Insert(_ uint64, data []byte) (uint64, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var id uint64
	for {
		id = uint64(rand.Uint32())
		if id == 0 {
			continue
		}
		if _, ok := m.cache[id]; ok {
			continue
		}
		break
	}
	item := newMockItem(id, data)
	m.cache[id] = item
	return id, nil
}

func (m *mockManage) LogDataItem(tid uint64, item Item) {
	println(tid, item)
}

func (m *mockManage) ReleaseDataItem(item Item) {
	println(item)
}

func (m *mockManage) TxManage() tx.Manage {
	return nil
}
