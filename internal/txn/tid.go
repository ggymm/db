package txn

import "encoding/binary"

type TID uint64

const (
	TIDSize = 8
)

func readTID(buf []byte) TID {
	return TID(binary.LittleEndian.Uint64(buf))
}

func writeTID(buf []byte, tid TID) {
	binary.LittleEndian.PutUint64(buf, uint64(tid))
}
