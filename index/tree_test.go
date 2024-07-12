package index

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/ggymm/db"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/data/page"
	"github.com/ggymm/db/pkg/file"
	"github.com/ggymm/db/tx"
)

func TestNewIndex(t *testing.T) {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/index")
	opt.Memory = page.Size * 10

	tm := tx.NewMockManage()
	dm := data.NewManage(tm, opt)

	index, err := NewIndex(dm, opt)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	t.Logf("%+v", index)
}

func TestTree_Gen(t *testing.T) {
	// 生成包含 1000 个 随机 uint64 的文件
	abs := db.RunPath()
	path := filepath.Join(abs, "temp/index/keys")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, file.Mode)
	if err != nil {
		t.Fatalf("open file err %v", err)
	}
	for i := 0; i < 1000; i++ {
		key := rand.Uint64()
		t.Logf("%d %d", i, key)
		// uint64 转换为 string
		_, err = f.WriteString(fmt.Sprintf("%d\n", key))
	}
	_ = f.Sync()
	_ = f.Close()
}

func TestIndex_Func(t *testing.T) {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/index")
	opt.Memory = page.Size * 10

	tm := tx.NewMockManage()
	dm := data.NewManage(tm, opt)

	var (
		err   error
		index Index

		lines  []string
		result []uint64
	)

	index, err = NewIndex(dm, opt)
	if err != nil {
		t.Fatalf("new index err %v", err)
	}

	lines, err = file.ReadLines(filepath.Join(abs, "temp/index/keys"))
	if err != nil {
		t.Fatalf("read lines err %v", err)
	}
	// 测试插入
	for i := 0; i < len(lines)-1; i++ {
		key, _ := strconv.ParseUint(lines[i], 10, 64)
		err = index.Insert(key, key)
		if err != nil {
			t.Fatalf("insert index err %v", err)
		}
		// t.Logf("insert index %d", key)
	}
	// 测试搜索
	for i := 0; i < len(lines)-1; i++ {
		key, _ := strconv.ParseUint(lines[i], 10, 64)
		result, err = index.Search(key)
		if err != nil {
			t.Fatalf("search index err %v", err)
		} else {
			if len(result) == 0 {
				t.Fatalf("search index err %d %v", key, result)
			} else {
				if result[0] != key {
					t.Fatalf("search index err %d %v", key, result)
				}
			}
		}
	}
	t.Log("search index success")
}

func TestIndex_FuncAsync(t *testing.T) {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/index")
	opt.Memory = page.Size * 10

	txManage := tx.NewMockManage()
	dataManage := data.NewManage(txManage, opt)

	var (
		err   error
		index Index

		insertNum = 50
		searchNum = 50
		taskNum   = 1000

		lock      sync.Mutex
		waitGroup = sync.WaitGroup{}
		cacheMap  = make(map[uint64]int)
	)

	index, err = NewIndex(dataManage, opt)
	if err != nil {
		t.Fatalf("new index err %v", err)
	}

	waitGroup.Add(insertNum + searchNum)

	// 插入
	for i := 0; i < insertNum; i++ {
		go func() {
			for j := 0; j < taskNum; j++ {
				key := rand.Uint64()
				err = index.Insert(key, key)
				if err != nil {
					t.Errorf("insert key %d err %v", key, err)
					continue
				}
				lock.Lock()
				cacheMap[key]++
				lock.Unlock()

				// t.Logf("insert key %d", key)
			}
			waitGroup.Done()
		}()
	}

	// 搜索
	for i := 0; i < searchNum; i++ {
		go func() {
			for j := 0; j < taskNum; j++ {
				prev := rand.Uint64()
				next := rand.Uint64()
				if prev > next {
					prev, next = next, prev
				}
				if next-prev > 10000 {
					next = prev + 10000
				}
				_, err = index.SearchRange(prev, next)
				if err != nil {
					continue
				}
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	// 检查
	t.Log("index check")
	for key, children := range cacheMap {
		res, _ := index.Search(key)

		if len(res) != children {
			t.Fatalf("error index check key %d %v %d", key, res, children)
		}
	}
	t.Log("index check success")
}
