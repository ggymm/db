package data

import (
	"bytes"

	"db/internal/data/page"
	"db/pkg/utils"
)

// 从缓存中移除
const (
	checkPos = 100
	checkLen = 8
)

func InitPage1() []byte {
	data := make([]byte, page.Size)
	setVcOpen(data)
	return data
}

func SetVcOpen(p page.Page) {
	p.SetDirty(true)
	setVcOpen(p.Data())
}

func SetVcClose(p page.Page) {
	p.SetDirty(true)
	setVcClose(p.Data())
}

func CheckVc(p page.Page) bool {
	return checkVc(p.Data())
}

func setVcOpen(data []byte) {
	// 将随机 byte 数组写入
	// 100 ~ 107 字节
	copy(data[checkPos:checkPos+checkLen], utils.RandBytes(checkLen))
}

func setVcClose(data []byte) {
	// 将 100 ~ 107 字节
	// 复制到 108 ~ 115 字节
	copy(data[checkPos+checkLen:checkPos+checkLen*2], data[checkPos:checkPos+checkLen])
}

func checkVc(data []byte) bool {
	return bytes.Compare(
		data[checkPos:checkPos+checkLen],
		data[checkPos+checkLen:checkPos+checkLen*2],
	) == 0
}
