package table

import (
	"hash/fnv"
)

func hash(val any) uint64 {
	switch val.(type) {
	case uint32:
		return uint64(val.(uint32))
	case uint64:
		return val.(uint64)
	case string:
		h := fnv.New64a()
		_, _ = h.Write([]byte(val.(string)))
		return h.Sum64()
	}
	return 0
}
