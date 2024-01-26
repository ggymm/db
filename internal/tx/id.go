package tx

import "encoding/binary"

const (
	IdLen        = 8
	Super uint64 = 0
)

func ReadId(buf []byte) uint64 {
	return readId(buf)
}

func WriteId(buf []byte, tid uint64) {
	writeId(buf, tid)
}

func readId(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

func writeId(buf []byte, tid uint64) {
	binary.LittleEndian.PutUint64(buf, tid)
}
