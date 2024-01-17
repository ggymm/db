package page

import (
	"encoding/binary"
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

func readPageFSO(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data[0:headLen])
}

func writePageFSO(data []byte, off uint16) {
	binary.LittleEndian.PutUint16(data[0:headLen], off)
}

func InitPageX() []byte {
	data := make([]byte, Size)
	writePageFSO(data, headLen) // 初始化写入 FSO
	return data
}

func MaxPageFree() int {
	return Size - headLen
}

func CalcPageFree(p Page) int {
	return int(Size - ParsePageFSO(p))
}

func ParsePageFSO(p Page) uint16 {
	return readPageFSO(p.Data())
}

func InsertPageData(p Page, data []byte) uint16 {
	p.SetDirty(true)
	off := ParsePageFSO(p)
	copy(p.Data()[off:], data)
	writePageFSO(p.Data(), off+uint16(len(data)))
	return off
}

func RecoverPageInsert(p Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)

	// 更新 FSO
	fso := ParsePageFSO(p)
	if off+uint16(len(data)) > fso {
		writePageFSO(p.Data(), off+uint16(len(data)))
	}
}

func RecoverPageUpdate(p Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)
}
