package sql

import (
	"slices"

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
		raw = bin.Uint32Raw(v.(uint32))
	case Int64:
		raw = bin.Uint64Raw(v.(uint64))
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
