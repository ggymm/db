package index

import (
	"path/filepath"
	"testing"

	"db/internal/data"
	"db/internal/data/page"
	"db/internal/opt"
	"db/internal/tx"
	"db/pkg/utils"
)

func TestNewIndex(t *testing.T) {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/index")

	o := &opt.Option{
		Open:   false,
		Path:   path,
		Name:   "test",
		Memory: page.Size * 10,
	}

	txManage := tx.NewMockManage()
	dataManage := data.NewManage(o, txManage)

	index, err := NewIndex(o, dataManage)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	t.Logf("%+v", index)
}

func TestIndex_Func(t *testing.T) {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/index")

	o := &opt.Option{
		Open:   false,
		Path:   path,
		Name:   "test",
		Memory: page.Size * 10,
	}

	txManage := tx.NewMockManage()
	dataManage := data.NewManage(o, txManage)

	var (
		err error

		index  Index
		result []uint64
	)

	index, err = NewIndex(o, dataManage)
	if err != nil {
		t.Fatalf("new index err %v", err)
	}

	// 测试插入
	limit := 99999
	for i := limit; i >= 0; i-- {
		err = index.Insert(uint64(i), uint64(i))
		if err != nil {
			t.Fatalf("insert index err %v", err)
		}
	}

	// 测试搜索
	for i := 0; i < limit; i++ {
		result, err = index.Search(uint64(i))
		if err != nil {
			t.Fatalf("search index err %v", err)
			return
		}
		if len(result) != 1 {
			t.Fatalf("search index err %v", err)
			return
		}
		if result[0] != uint64(i) {
			t.Fatalf("search index err %v", err)
			return
		}
	}
	t.Log("search index success")
}
