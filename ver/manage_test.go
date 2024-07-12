package ver

import (
	"testing"

	"github.com/ggymm/db"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/tx"
)

func TestVerManage_Handle(t *testing.T) {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/ver")
	opt.Memory = (1 << 20) * 64

	tm := tx.NewManager(opt)
	dm := data.NewManage(tm, opt)

	vm := NewManage(tm, dm)
	txId := vm.Begin(0)

	// insert
	src := []byte("test")
	dataId, err := vm.Insert(txId, src)
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	err = vm.Commit(txId)
	if err != nil {
		t.Fatalf("commit error: %s", err)
	}

	// read
	var (
		dst   []byte
		exist bool
	)
	dst, exist, err = vm.Read(tx.Super, dataId)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}
	if !exist {
		t.Fatalf("read error: not exist")
	}
	if err != nil {
		t.Fatalf("commit error: %s", err)
	}

	t.Logf("read data: %s", string(src))
	t.Logf("read data: %s", string(dst))
}
