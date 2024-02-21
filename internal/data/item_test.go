package data

import (
	"testing"
)

func TestDataItem_Id(t *testing.T) {
	var no uint32 = 1 << 31
	var offset uint16 = 1 << 15

	id := wrapDataItemId(no, offset)
	no1, offset1 := parseDataItemId(id)

	t.Logf("id %d", id)
	t.Logf("src no %d, off: %d", no, offset)
	t.Logf("dst no %d, off: %d", no1, offset1)
	if no != no1 || offset != offset1 {
		t.Fatal("parseDataItemId/wrapDataItemId error")
	}
	t.Log("parseDataItemId/wrapDataItemId success")
}

func TestDataItem_parseDataItemId(t *testing.T) {
	no, off := parseDataItemId(132838)
	t.Logf("no: %d off: %d", no, off)
}
