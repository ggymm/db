package bin

import "encoding/binary"

var bin = binary.LittleEndian

func Uint16(data []byte) uint16 {
	return bin.Uint16(data)
}

func PutUint16(data []byte, num uint16) {
	bin.PutUint16(data, num)
}

func Uint32(data []byte) uint32 {
	return bin.Uint32(data)
}

func PutUint32(data []byte, num uint32) {
	bin.PutUint32(data, num)
}

func Uint64(data []byte) uint64 {
	return bin.Uint64(data)
}

func PutUint64(data []byte, num uint64) {
	bin.PutUint64(data, num)
}
