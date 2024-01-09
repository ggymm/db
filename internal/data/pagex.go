package data

import (
	"encoding/binary"

	"db/internal/data/page"
)

// 普通页
//
// 页面数据结构如下：
// +----------------+----------------+
// |      FSO       |      data      |
// +----------------+----------------+
// |     2 byte     |    * byte      |
// +----------------+----------------+
// FSO（Free Space Offset）：空闲空间偏移量
//
// 使用 uint16 存储 FSO（最大可支持 64k 页面大小）

const (
	headLen = 2
)

func parseFSO(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data[0:headLen])
}

func updateFSO(data []byte, off uint16) {
	binary.LittleEndian.PutUint16(data[0:headLen], off)
}

func InitPageX() []byte {
	data := make([]byte, page.Size)
	updateFSO(data, headLen) // 初始化写入 FSO
	return data
}

func MaxFree() int {
	return page.Size - headLen
}

func ParseFSO(p page.Page) uint16 {
	return parseFSO(p.Data())
}

func ParseFree(p page.Page) int {
	return int(page.Size - parseFSO(p.Data()))
}

func InsertData(p page.Page, data []byte) uint16 {
	p.SetDirty(true)
	off := parseFSO(p.Data())
	copy(p.Data()[off:], data)
	updateFSO(p.Data(), off+uint16(len(data)))
	return off
}

func RecoverInsert(p page.Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)

	// 更新 FSO
	fso := parseFSO(p.Data())
	if off+uint16(len(data)) > fso {
		updateFSO(p.Data(), off+uint16(len(data)))
	}
}

func RecoverUpdate(p page.Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)
}
