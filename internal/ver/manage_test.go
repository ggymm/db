package ver

import (
	"path/filepath"
	"testing"

	"db/internal/app"
	"db/internal/data"
	"db/internal/tx"
)

func TestVerManage_Handle(t *testing.T) {
	base := app.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/ver")

	tm := tx.NewManager(&app.Option{
		Open: false,
		Name: name,
		Path: path,
	})
	dm := data.NewManage(tm, &app.Option{
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
