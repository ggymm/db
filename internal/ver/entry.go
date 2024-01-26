package ver

import (
	"db/internal/data"
	"db/internal/tx"
)

// 数据记录（带版本）
//
// 数据记录的机构如下
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |     8 byte     |     8 byte     |     * byte     |
// +----------------+----------------+----------------+
//
// mix：代表创建该记录的事务Id
// max：代表删除该记录的事务Id（或者出现新版本）
// data：数据内容
//
// 可以理解为，entry 就是 data_item 结构中的数据内容
// 同时，结构内的 min 和 max 共同约束的是该数据记录的可见性（版本）

const (
	offMin  = 0
	offMax  = offMin + tx.IdLen
	offData = offMax + tx.IdLen
)

type entry struct {
	id     uint64 // data_item 的 id
	item   data.Item
	manage Manage
}

func (e *entry) Min() uint64 {
	e.item.RLock()
	defer e.item.RUnlock()

	return tx.ReadId(e.item.DataBody()[offMin:])
}

func (e *entry) Max() uint64 {
	e.item.RLock()
	defer e.item.RUnlock()

	return tx.ReadId(e.item.DataBody()[offMax:])
}

func (e *entry) Data() (data []byte) {
	e.item.RLock()
	defer e.item.RUnlock()

	data = make([]byte, len(e.item.DataBody())-offData)
	copy(data, e.item.DataBody()[offData:])
	return
}

func (e *entry) SetMax(tid uint64) {
	e.item.Before()
	defer e.item.After(tid)

	tx.WriteId(e.item.DataBody()[offMax:], tid)
}
