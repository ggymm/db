package sql

import (
	"slices"
	"strconv"

	"db/pkg/bin"
	"db/pkg/str"
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
		raw = str.Serialize(v.(string))
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
		v, shift = str.Deserialize(raw)
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
