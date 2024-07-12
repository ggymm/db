package table

import (
	"github.com/ggymm/db/pkg/bin"
)

func encode(value any) []byte {
	switch v := value.(type) {
	case uint64:
		return encodeUint64(v)
	case string:
		return encodeString(v)
	default:
		panic("unsupported type")
	}
}

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

func decode[T any](data []byte) (T, int) {
	var value T
	switch any(value).(type) {
	case uint64:
		value = any(bin.Uint64(data)).(T)
		return value, 8
	case string:
		l := int(data[0]) |
			int(data[1])<<8 |
			int(data[2])<<16 |
			int(data[3])<<24
		return any(string(data[4 : 4+l])).(T), 4 + l
	default:
		panic("unsupported type")
	}
}

func decodeUint64(v []byte) uint64 {
	return bin.Uint64(v)
}

func decodeString(v []byte) (string, int) {
	l := int(v[0]) |
		int(v[1])<<8 |
		int(v[2])<<16 |
		int(v[3])<<24
	return string(v[4 : 4+l]), 4 + l
}
