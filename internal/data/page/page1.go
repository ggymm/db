package page

import (
	"bytes"
	"math/rand"
	"time"
)

// 第一页
//
// 特殊处理，用于验证数据库是否被正常关闭
// 初始化时，在 100 ~ 107 字节写入随机 byte 数组
// 正常关闭时，将 100 ~ 107 字节复制到 108 ~ 115 字节
// 每次启动时，检查 100 ~ 107 字节和 108 ~ 115 字节是否相同

const (
	checkPos = 100
	checkLen = 8
)

func randB(n int) []byte {
	b := make([]byte, n)
	// 为随机数生成器提供一个种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = byte(r.Intn(256)) // 随机生成一个字节
	}
	return b
}

func setVcOpen(data []byte) {
	// 将随机 byte 数组写入
	// 100 ~ 107 字节
	copy(
		data[checkPos:checkPos+checkLen], // dst
		randB(checkLen),                  // src
	)
}

func setVcClose(data []byte) {
	// 将 100 ~ 107 字节
	// 复制到 108 ~ 115 字节
	copy(
		data[checkPos+checkLen:checkPos+checkLen*2], // dst
		data[checkPos:checkPos+checkLen],            // src
	)
}

func checkVc(data []byte) bool {
	return bytes.Compare(
		data[checkPos:checkPos+checkLen],            // 100 ~ 107 字节
		data[checkPos+checkLen:checkPos+checkLen*2], // 108 ~ 115 字节
	) == 0
}

func InitPage1() []byte {
	data := make([]byte, Size)
	setVcOpen(data)
	return data
}

func SetVcOpen(p Page) {
	p.SetDirty(true)
	setVcOpen(p.Data())
}

func SetVcClose(p Page) {
	p.SetDirty(true)
	setVcClose(p.Data())
}

func CheckVc(p Page) bool {
	return checkVc(p.Data())
}
