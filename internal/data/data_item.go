package data

import (
	"db/internal/data/page"
	"db/internal/txn"
	"encoding/binary"
	"sync"
)

// 数据对象（保存在 page 中的数据对象）
//
// dataItem 的 Id 由 Page 编号和 Page 中的偏移量组成
// +----------------+----------------+
// |       no       |     offset     |
// +----------------+----------------+
// |     4 byte     |     4 byte     |
// +----------------+----------------+
// 其中 no 为 uint32 类型，offset 为 uint16 类型
//
// dataItem 的 data 数据结构如下：
// +----------------+----------------+----------------+
// |     flag       |      size      |      data      |
// +----------------+----------------+----------------+
// |    1 byte      |     2 byte     |     * byte     |
// +----------------+----------------+----------------+
//
// flag: 1 byte，标记数据是否合法（0 表示合法，1 表示非法）
// size: 2 byte，标记 data 的长度
// data: * byte，数据内容

const (
	offFlag = 0
	offSize = 1 // flag 占用 1 字节
	offData = 3 // size 占用 2 字节
)

type Item interface {
	Id() uint64
	Data() []byte
	Page() page.Page
	Valid() bool

	Before()
	UnBefore()
	After(tid txn.TID)
	Release()

	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type dataItem struct {
	lock sync.RWMutex

	id      uint64
	data    []byte
	oldData []byte

	page       page.Page
	dataManage Manage
}

func readSize(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf)
}

func writeSize(buf []byte, size uint16) {
	binary.LittleEndian.PutUint16(buf, size)
}

func wrapDataItem(data []byte) []byte {
	buf := make([]byte, len(data)+offData)
	writeSize(buf[offSize:], uint16(len(data)))
	copy(buf[offData:], data)
	return data
}

func parseDataItem(p page.Page, off uint16, m Manage) Item {
	id := wrapDataItemId(p.No(), off)

	data := p.Data()[off:]
	size := readSize(data[offSize:])
	return &dataItem{
		id:      id,
		data:    data[:size],
		oldData: make([]byte, size),

		page:       p,
		dataManage: m,
	}
}

func wrapDataItemId(no uint32, off uint16) uint64 {
	return uint64(no)<<16 | uint64(off)
}

func parseDataItemId(id uint64) (no uint32, off uint16) {
	no = uint32(id >> 16)
	off = uint16(id & 0xffff)
	return no, off
}

func (item *dataItem) Id() uint64 {
	return item.id
}

func (item *dataItem) Data() []byte {
	return item.data[offData:]
}

func (item *dataItem) Page() page.Page {
	return item.page
}

func (item *dataItem) Valid() bool {
	return item.data[offFlag] == 0
}

func (item *dataItem) Before() {
	item.lock.Lock()
	item.page.SetDirty(true)
	copy(item.oldData, item.data)
}

func (item *dataItem) UnBefore() {
	copy(item.data, item.oldData)
	item.lock.Unlock()
}

func (item *dataItem) After(tid txn.TID) {
	item.dataManage.LogDataItem(tid, item)
	item.lock.Unlock()
}

func (item *dataItem) Release() {
	item.dataManage.ReleaseDataItem(item)
}

func (item *dataItem) Lock() {
	item.lock.Lock()
}

func (item *dataItem) Unlock() {
	item.lock.Unlock()
}

func (item *dataItem) RLock() {
	item.lock.RLock()
}

func (item *dataItem) RUnlock() {
	item.lock.RUnlock()
}
