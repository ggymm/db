package page

import (
	"os"
	"path/filepath"
	"testing"

	"db/internal/ops"
	"db/pkg/utils"
)

func newOps(open bool) *ops.Option {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/db")
	return &ops.Option{
		Open:   open,
		Path:   path,
		Memory: (1 << 20) * 64,
	}
}

func TestNewPage(t *testing.T) {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/db")
	// 清空目录
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}

	cache := NewCache(newOps(false))
	t.Logf("%+v", cache)

	for i := 1; i <= 100; i++ {
		data := utils.RandBytes(Size)
		no := cache.NewPage(data)
		p, e := cache.ObtainPage(no)
		if e != nil {
			t.Fatalf("err %v", e)
			return
		}
		p.SetDirty(true)
		p.Data()[0] = byte(i)
		p.Release()
	}
	cache.Close()

	cache = NewCache(newOps(true))
	for i := 1; i <= 100; i++ {
		p, e := cache.ObtainPage(uint32(i))
		if e != nil {
			t.Fatalf("err %v", e)
			return
		}
		if p.Data()[0] == byte(i) {
			t.Logf("success %d %v", i, p.Data()[0])
		} else {
			t.Fatalf("err %d %v", i, p.Data()[0])
		}
		p.Release()
	}
	cache.Close()
}
