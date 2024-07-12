package page

import (
	"testing"
)

func Test_PageX(t *testing.T) {
	p := NewPage(1, NewPageX(), nil)

	t.Logf("%+v", p)
	t.Logf("max free %d", MaxPageFree())
	t.Logf("current fso %d", ParsePageFSO(p))
	t.Logf("current free %d", CalcPageFree(p))

	// 模拟插入数据
	data := []byte("hello world")
	WritePageData(p, data)
	t.Logf("%+v", p)
	t.Logf("data len %d", len(data))
	t.Logf("current fso %d", ParsePageFSO(p))
	t.Logf("current free %d", CalcPageFree(p))
}
