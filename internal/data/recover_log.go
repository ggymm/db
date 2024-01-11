package data

import (
	"db/internal/data/page"
	"db/internal/txn"
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

const (
	typeLen = 1

	redoLog = 1
	undoLog = 2

	InsertLog = 1
	UpdateLog = 2
)

func wrapInsertLog(tid uint64, p page.Page, data []byte) []byte {
	// type: 1; tid: 8; itemId: 8
	l := typeLen + txn.TIDLen + itemIdLen + len(data)
	log := make([]byte, l)

	pos := 0
	log[pos] = InsertLog // type

	pos += typeLen
	txn.WriteTID(log[pos:], tid) // tid

	pos += txn.TIDLen
	off := parsePageFSO(p)
	itemId := wrapDataItemId(p.No(), off)
	writeDataItemId(log[pos:], itemId) // item_id

	pos += itemIdLen
	copy(log[pos:], data) // data
	return log
}

func parseInsertLog(log []byte) (uint64, uint32, uint16, []byte) {
	pos := typeLen
	tid := txn.ReadTID(log[pos:]) // tid

	pos += txn.TIDLen
	itemId := readDataItemId(log[pos:]) // item_id
	no, off := parseDataItemId(itemId)

	pos += itemIdLen
	data := log[pos:]
	return tid, no, off, data
}

func wrapUpdateLog(tid uint64, item Item) []byte {
	// type: 1; tid: 8; itemId: 8
	l := typeLen + txn.TIDLen + itemIdLen + len(item.Data())*2
	log := make([]byte, l)

	pos := 0
	log[pos] = UpdateLog // type

	pos += typeLen
	txn.WriteTID(log[pos:], tid) // tid

	pos += txn.TIDLen
	writeDataItemId(log[pos:], item.Id()) // item_id

	pos += itemIdLen
	copy(log[pos:], item.DataOld()) // old_data

	pos += len(item.DataOld())
	copy(log[pos:], item.Data())
	return log
}

func parseUpdateLog(log []byte) (uint64, uint32, uint16, []byte, []byte) {
	pos := typeLen
	tid := txn.ReadTID(log[pos:]) // tid

	pos += txn.TIDLen
	itemId := readDataItemId(log[pos:]) // item_id
	no, off := parseDataItemId(itemId)

	pos += itemIdLen
	dataLen := (len(log) - pos) / 2
	dataOld := log[pos : pos+dataLen]
	dataNew := log[pos+dataLen : pos+dataLen*2]
	return tid, no, off, dataOld, dataNew
}
