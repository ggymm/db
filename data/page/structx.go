package page

import (
	"github.com/ggymm/db/pkg/bin"
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

func readPageOffset(data []byte) uint16 {
	return bin.Uint16(data[0:headLen])
}

func writePageOffset(data []byte, off uint16) {
	bin.PutUint16(data[0:headLen], off)
}

func NewPageX() []byte {
	data := make([]byte, Size)
	writePageOffset(data, headLen) // 初始化写入 FSO
	return data
}

func MaxPageFree() uint32 {
	return Size - headLen
}

func CalcPageFree(p Page) uint32 {
	return Size - uint32(readPageOffset(p.Data()))
}

func ParsePageFSO(p Page) uint16 {
	return readPageOffset(p.Data())
}

func InsertPageData(p Page, data []byte) uint16 {
	p.SetDirty(true)
	off := readPageOffset(p.Data())
	copy(p.Data()[off:], data)
	writePageOffset(p.Data(), off+uint16(len(data)))
	return off
}

func RecoverPageInsert(p Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)

	// 更新 FSO
	fso := readPageOffset(p.Data())
	if off+uint16(len(data)) > fso {
		writePageOffset(p.Data(), off+uint16(len(data)))
	}
}

func RecoverPageUpdate(p Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)
}
