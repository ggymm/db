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

func initPageX() []byte {
	data := make([]byte, page.Size)
	updatePageFSO(data, headLen) // 初始化写入 FSO
	return data
}

func maxPageFree() int {
	return page.Size - headLen
}

func calcPageFree(p page.Page) int {
	return int(page.Size - parsePageFSO(p))
}

func parsePageFSO(p page.Page) uint16 {
	return binary.LittleEndian.Uint16(p.Data()[0:headLen])
}

func updatePageFSO(data []byte, off uint16) {
	binary.LittleEndian.PutUint16(data[0:headLen], off)
}

func insertPageData(p page.Page, data []byte) uint16 {
	p.SetDirty(true)
	off := parsePageFSO(p)
	copy(p.Data()[off:], data)
	updatePageFSO(p.Data(), off+uint16(len(data)))
	return off
}

func recoverPageInsert(p page.Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)

	// 更新 FSO
	fso := parsePageFSO(p)
	if off+uint16(len(data)) > fso {
		updatePageFSO(p.Data(), off+uint16(len(data)))
	}
}

func recoverPageUpdate(p page.Page, off uint16, data []byte) {
	p.SetDirty(true)
	copy(p.Data()[off:], data)
}
