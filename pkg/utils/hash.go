package utils

import (
	"hash/fnv"
)

// Hash 返回任意值对应的哈希值
func Hash(val any) uint64 {
	switch val.(type) {
	case uint32:
		return uint64(val.(uint32))
	case uint64:
		return val.(uint64)
	case string:
		hash := fnv.New64a()
		_, _ = hash.Write([]byte(val.(string)))
		return hash.Sum64()
	}
	return 0
}
