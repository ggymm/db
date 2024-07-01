package sql

import (
	"slices"
	"strconv"

	"db/pkg/bin"
)

const (
	Null byte = iota
	NotNull
)

func FieldRaw(t string, v any) []byte {
	var raw []byte

	if v == nil {
		return []byte{Null}
	}
	switch typeMapping[t] {
	case Int32:
		i, _ := strconv.ParseUint(v.(string), 10, 32)
		raw = bin.Uint32Raw(uint32(i))
	case Int64:
		i, _ := strconv.ParseUint(v.(string), 10, 64)
		raw = bin.Uint64Raw(i)
	case Varchar:
		// length
		l := len(v.(string))
		raw = make([]byte, 4+l)
		raw[0] = byte(l)
		raw[1] = byte(l >> 8)
		raw[2] = byte(l >> 16)
		raw[3] = byte(l >> 24)

		// string
		copy(raw[4:], v.(string))
	}
	return slices.Insert(raw, 0, NotNull)
}

func FieldParse(t string, raw []byte) (any, int) {
	var v any
	var shift int

	if raw[0] == Null {
		return nil, 1
	}
	raw = raw[1:]
	switch typeMapping[t] {
	case Int32:
		v = bin.Uint32(raw)
		shift = 4
	case Int64:
		v = bin.Uint64(raw)
		shift = 8
	case Varchar:
		l := int(raw[0]) |
			int(raw[1])<<8 |
			int(raw[2])<<16 |
			int(raw[3])<<24
		v = string(raw[4 : 4+l])
		shift = l + 4
	}
	return v, shift + 1
}

func FieldFormat(t string, v string) any {
	switch typeMapping[t] {
	case Int32:
		i, _ := strconv.ParseUint(v, 10, 32)
		return uint32(i)
	case Int64:
		i, _ := strconv.ParseUint(v, 10, 64)
		return i
	case Varchar:
		return v
	}
	return nil
}
