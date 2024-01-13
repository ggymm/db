package tx

import "encoding/binary"

const (
	IdLen        = 8
	Super uint64 = 0
)

func ReadTID(buf []byte) uint64 {
	return readTID(buf)
}

func WriteTID(buf []byte, tid uint64) {
	writeTID(buf, tid)
}

func readTID(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

func writeTID(buf []byte, tid uint64) {
	binary.LittleEndian.PutUint64(buf, tid)
}
