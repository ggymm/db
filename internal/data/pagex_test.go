package data

import (
	"testing"

	"db/internal/data/page"
)

func Test_PageX(t *testing.T) {
	p := page.NewPage(1, InitPageX(), nil)

	t.Logf("%+v", p)
	t.Logf("max free: %d", MaxFree())
	t.Logf("current fso: %d", ParseFSO(p))
	t.Logf("current free: %d", ParseFree(p))

	// 模拟插入数据
	data := []byte("hello world")
	InsertData(p, data)
	t.Logf("%+v", p)
	t.Logf("data len: %d", len(data))
	t.Logf("current fso: %d", ParseFSO(p))
	t.Logf("current free: %d", ParseFree(p))
}
