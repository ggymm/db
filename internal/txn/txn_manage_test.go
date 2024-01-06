package txn

import (
	"db/pkg/utils"
	"path/filepath"
	"testing"
)

func TestNewTxnManager(t *testing.T) {
	base := utils.RunPath()
	filename := filepath.Join(base, "temp/txn/test", "txn")

	tm := NewTxnManager(filename)
	t.Log(base)
	t.Logf("%+v", tm)

	tm.Close()
}

func TestTxnManager_State(t *testing.T) {
	base := utils.RunPath()
	filename := filepath.Join(base, "temp/txn/test", "txn")

	tm := NewTxnManager(filename)
	t.Log(base)
	t.Logf("%+v", tm)

	tid := tm.Begin()
	t.Logf("%d is active %t", tid, tm.IsActive(tid))

	tm.Commit(tid)
	t.Logf("%d is committed %t", tid, tm.IsCommitted(tid))

	tid2 := tm.Begin()
	t.Logf("%d is active %t", tid2, tm.IsActive(tid2))

	tm.Abort(tid2)
	t.Logf("%d is abord %t", tid2, tm.IsAborted(tid2))

	tm.Close()
}
