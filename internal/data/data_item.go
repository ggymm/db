package data

import (
	"db/internal/data/page"
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

	Before()
	UnBefore()
	After()
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

func (d *dataItem) Id() uint64 {
	return d.id
}

func (d *dataItem) Data() []byte {
	return d.data[offData:]
}

func (d *dataItem) Page() page.Page {
	return d.page
}

func (d *dataItem) Before() {
	d.lock.Lock()
	d.page.SetDirty(true)
	copy(d.oldData, d.data)
}

func (d *dataItem) UnBefore() {
	copy(d.data, d.oldData)
	d.lock.Unlock()
}

func (d *dataItem) After() {
	//TODO implement me
	panic("implement me")
}

func (d *dataItem) Release() {
	//TODO implement me
	panic("implement me")
}

func (d *dataItem) Lock() {
	d.lock.Lock()
}

func (d *dataItem) Unlock() {
	d.lock.Unlock()
}

func (d *dataItem) RLock() {
	d.lock.RLock()
}

func (d *dataItem) RUnlock() {
	d.lock.RUnlock()
}
