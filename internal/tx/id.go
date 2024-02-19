package tx

import (
	"db/pkg/bin"
)

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
	return bin.Uint64(buf)
}

func writeId(buf []byte, tid uint64) {
	bin.PutUint64(buf, tid)
}
