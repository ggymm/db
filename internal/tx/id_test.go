package tx

import (
	"testing"
)

func Test_TID(t *testing.T) {
	var src uint64 = 128
	buf := make([]byte, 8)
	writeTID(buf, src)
	t.Logf("%+v", buf)

	dst := readTID(buf)
	t.Logf("%+v", buf)

	t.Logf("src %+v, dst %+v", src, dst)
	if src == dst {
		t.Log("TID test pass")
	} else {
		t.Error("TID test failed")
	}
}
