package table

import (
	"testing"
)

func Test_Encode_Decode(t *testing.T) {
	buf1 := encode(uint64(1))
	t.Logf("%v", buf1)

	buf2 := encode("hello")
	t.Logf("%v", buf2)

	v1, n1 := decode[uint64](buf1)
	t.Logf("%v %d", v1, n1)

	v2, n2 := decode[string](buf2)
	t.Logf("%v %d", v2, n2)
}
