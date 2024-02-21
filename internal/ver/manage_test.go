package ver

import (
	"path/filepath"
	"testing"

	"db/internal/data"

	"db/internal/opt"
	"db/internal/tx"
	"db/pkg/utils"
)

func TestVerManage_Handle(t *testing.T) {
	base := utils.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/ver")

	tm := tx.NewManager(&opt.Option{
		Open: false,
		Name: name,
		Path: path,
	})
	dm := data.NewManage(tm, &opt.Option{
		Open:   false,
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	})

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
	txId = vm.Begin(0)
	dst, exist, err = vm.Read(txId, dataId)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}
	if !exist {
		t.Fatalf("read error: not exist")
	}
	err = vm.Commit(txId)
	if err != nil {
		t.Fatalf("commit error: %s", err)
	}

	t.Logf("read data: %s", string(src))
	t.Logf("read data: %s", string(dst))
}
