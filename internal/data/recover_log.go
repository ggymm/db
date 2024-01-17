package data

import (
	"db/internal/data/page"
	"db/internal/tx"
)

// 数据日志
//
// 插入日志，数据的结构如下：
// +----------------+----------------+----------------+----------------+
// |      type      |      tid       |     item_id    |      data      |
// +----------------+----------------+----------------+----------------+
// |     1 byte     |     8 byte     |     8 byte     |    * byte      |
// +----------------+----------------+----------------+----------------+
//
// 更新日志，数据的结构如下：
// +----------------+----------------+----------------+----------------+----------------+
// |    type        |      tid       |     item_id    |     old_data   |     new_data   |
// +----------------+----------------+----------------+----------------+----------------+
// |    1 byte      |     8 byte     |     8 byte     |     * byte     |     * byte     |
// +----------------+----------------+----------------+----------------+----------------+
//
// 如何保证 old_data 和 new_data 的长度相同？
// 通过上层业务保证
// 因为此时 data 表示的是数据表中的每一行数据，所以只需要保证数据字段的长度为固定值即可

const (
	typeLen = 1

	redoLog = 1
	undoLog = 2

	InsertLog = 1
	UpdateLog = 2
)

func wrapInsertLog(tid uint64, p page.Page, data []byte) []byte {
	// type: 1; tid: 8; itemId: 8
	l := typeLen + tx.IdLen + itemIdLen + len(data)
	log := make([]byte, l)

	pos := 0
	log[pos] = InsertLog // type

	pos += typeLen
	tx.WriteTID(log[pos:], tid) // tid

	pos += tx.IdLen
	off := page.ParsePageFSO(p)
	itemId := wrapDataItemId(p.No(), off)
	writeDataItemId(log[pos:], itemId) // item_id

	pos += itemIdLen
	copy(log[pos:], data) // data
	return log
}

func parseInsertLog(log []byte) (uint64, uint32, uint16, []byte) {
	pos := typeLen
	tid := tx.ReadTID(log[pos:]) // tid

	pos += tx.IdLen
	itemId := readDataItemId(log[pos:]) // item_id
	no, off := parseDataItemId(itemId)

	pos += itemIdLen
	data := log[pos:]
	return tid, no, off, data
}

func wrapUpdateLog(tid uint64, item Item) []byte {
	// type: 1; tid: 8; itemId: 8
	l := typeLen + tx.IdLen + itemIdLen + len(item.Data())*2
	log := make([]byte, l)

	pos := 0
	log[pos] = UpdateLog // type

	pos += typeLen
	tx.WriteTID(log[pos:], tid) // tid

	pos += tx.IdLen
	writeDataItemId(log[pos:], item.Id()) // item_id

	pos += itemIdLen
	copy(log[pos:], item.DataOld()) // old_data

	pos += len(item.DataOld())
	copy(log[pos:], item.Data())
	return log
}

func parseUpdateLog(log []byte) (uint64, uint32, uint16, []byte, []byte) {
	pos := typeLen
	tid := tx.ReadTID(log[pos:]) // tid

	pos += tx.IdLen
	itemId := readDataItemId(log[pos:]) // item_id
	no, off := parseDataItemId(itemId)

	pos += itemIdLen
	dataLen := (len(log) - pos) / 2
	dataOld := log[pos : pos+dataLen]
	dataNew := log[pos+dataLen : pos+dataLen*2]
	return tid, no, off, dataOld, dataNew
}
