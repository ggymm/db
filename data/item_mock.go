package data

import (
	"sync"

	"db/data/page"
)

type mockItem struct {
	lock sync.RWMutex

	id      uint64
	data    []byte
	oldData []byte
}

func newMockItem(id uint64, data []byte) *mockItem {
	return &mockItem{
		id:      id,
		data:    data,
		oldData: make([]byte, len(data)),
	}
}

func (item *mockItem) Id() uint64 {
	return item.id
}

func (item *mockItem) Flag() bool {
	return item.data[offFlag] == 0
}

func (item *mockItem) Page() page.Page {
	return nil
}

func (item *mockItem) Data() []byte {
	return item.data
}

func (item *mockItem) DataOld() []byte {
	return item.oldData
}

func (item *mockItem) DataBody() []byte {
	return item.data
}

func (item *mockItem) Before() {
	item.lock.Lock()
	copy(item.oldData, item.data)
}

func (item *mockItem) UnBefore() {
	copy(item.data, item.oldData)
	item.lock.Unlock()
}

func (item *mockItem) After(_ uint64) {
	item.lock.Unlock()
}

func (item *mockItem) Release() {
}

func (item *mockItem) Lock() {
	item.lock.Lock()
}

func (item *mockItem) Unlock() {
	item.lock.Unlock()
}

func (item *mockItem) RLock() {
	item.lock.RLock()
}

func (item *mockItem) RUnlock() {
	item.lock.RUnlock()
}
