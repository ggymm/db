package page

import (
	"os"
	"testing"

	"github.com/ggymm/db"
)

func TestNewPage(t *testing.T) {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/page")
	opt.Memory = (1 << 20) * 64

	// 清空目录
	err := os.RemoveAll(opt.Path)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}

	cache := NewManage(opt)
	t.Logf("%+v", cache)

	for i := 1; i <= 100; i++ {
		data := randB(Size)
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

	opt = db.NewOption(abs, "temp/page")
	opt.Memory = (1 << 20) * 64
	cache = NewManage(opt)
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
