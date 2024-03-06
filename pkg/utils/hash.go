package utils

import "hash/fnv"

// Hash 返回任意值对应的哈希值
func Hash(val any) uint64 {
	switch val.(type) {
	case uint32, uint64:
		return uint64(val.(uint))
	case string:
		hash := fnv.New64a()
		_, _ = hash.Write([]byte(val.(string)))
		return hash.Sum64()
	}
	return 0
}
