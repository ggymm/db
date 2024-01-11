package data

import (
	"os"
	"path/filepath"
	"testing"

	"db/internal/ops"
	"db/internal/txn"
	"db/pkg/utils"
)

func newOps(open bool) *ops.Option {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/data")
	if !open && !utils.IsEmpty(path) {
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}

	return &ops.Option{
		Open:   open,
		Path:   path,
		Memory: (1 << 20) * 64,
	}
}

func TestNewManage(t *testing.T) {
	o := newOps(false)
	tm := txn.NewManager(o)
	dm := NewManage(o, tm)
	t.Logf("%+v", dm)
}

func TestDataManage_DataHandle(t *testing.T) {
	o := newOps(true)
	tm := txn.NewManager(o)
	dm := NewManage(o, tm)
	t.Logf("%+v", dm)

	data := utils.RandBytes(60)
	t.Logf("data: %+v", data)

	tid := tm.Begin()
	defer tm.Commit(tid)
	id, err := dm.Insert(tid, data)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}
	t.Log(id)

	tid = tm.Begin()
	defer tm.Commit(tid)
	item, ok, err := dm.Read(id)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}
	if !ok {
		t.Fatalf("not ok")
		return
	}
	t.Logf("%+v", item)
}
