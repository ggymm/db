package tx

import (
	"testing"
)

func Test_Id(t *testing.T) {
	var src uint64 = 128
	buf := make([]byte, 8)
	writeId(buf, src)
	t.Logf("%+v", buf)

	dst := readId(buf)
	t.Logf("%+v", buf)

	t.Logf("src %+v, dst %+v", src, dst)
	if src == dst {
		t.Log("Id test pass")
	} else {
		t.Error("Id test failed")
	}
}
