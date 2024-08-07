package table

import (
	"github.com/ggymm/db/pkg/bin"
)

func encodeUint64(v uint64) []byte {
	return bin.Uint64Raw(v)
}

func encodeString(v string) []byte {
	// length
	l := len(v)
	data := make([]byte, 4+l)
	data[0] = byte(l)
	data[1] = byte(l >> 8)
	data[2] = byte(l >> 16)
	data[3] = byte(l >> 24)

	// string
	copy(data[4:], v)
	return data
}

func decodeUint64(v []byte) (uint64, int) {
	return bin.Uint64(v), 8
}

func decodeString(v []byte) (string, int) {
	l := int(v[0]) |
		int(v[1])<<8 |
		int(v[2])<<16 |
		int(v[3])<<24
	return string(v[4 : 4+l]), 4 + l
}
