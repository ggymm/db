package data

import (
	"testing"

	"db/internal/data/page"
)

func Test_PageX(t *testing.T) {
	p := page.NewPage(1, initPageX(), nil)

	t.Logf("%+v", p)
	t.Logf("max free: %d", maxPageFree())
	t.Logf("current fso: %d", parsePageFSO(p))
	t.Logf("current free: %d", calcPageFree(p))

	// 模拟插入数据
	data := []byte("hello world")
	insertPageData(p, data)
	t.Logf("%+v", p)
	t.Logf("data len: %d", len(data))
	t.Logf("current fso: %d", parsePageFSO(p))
	t.Logf("current free: %d", calcPageFree(p))
}
