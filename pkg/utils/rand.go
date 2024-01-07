package utils

import (
	"math/rand"
	"time"
)

func RandBytes(n int) []byte {
	b := make([]byte, n)
	// 为随机数生成器提供一个种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = byte(rand.Intn(256)) // 随机生成一个字节
	}
	return b
}
