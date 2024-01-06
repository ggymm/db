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
}

func TestTxnManager_Begin(t *testing.T) {

}
