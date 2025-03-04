package sql

import (
	"strconv"
)

func FormatVal(t string, v string) any {
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
