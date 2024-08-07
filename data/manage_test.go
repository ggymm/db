package data

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ggymm/db"
	"github.com/ggymm/db/tx"
)

func newOpt() *db.Option {
	opt := db.NewOption(db.RunPath(), "temp/data")
	opt.Memory = (1 << 20) * 64
	return opt
}

func randB(n int) []byte {
	b := make([]byte, n)
	// 为随机数生成器提供一个种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = byte(r.Intn(256)) // 随机生成一个字节
	}
	return b
}

func TestNewManage(t *testing.T) {
	opt := newOpt()
	tm := tx.NewManager(opt)
	dm := NewManage(tm, opt)
	t.Logf("%+v", dm)
}

func TestDataManage_DataHandle(t *testing.T) {
	abs := db.RunPath()
	path := filepath.Join(abs, "temp/data")

	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}
	opt := newOpt()

	tm := tx.NewManager(opt)
	dm := NewManage(tm, opt)
	t.Logf("%+v", dm)

	data := randB(60)
	t.Logf("data %+v", data)

	tid := tm.Begin()
	defer tm.Commit(tid)
	id, err := dm.Write(tid, data)
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
	t.Logf("data %+v", item.DataBody())

	if bytes.Equal(data, item.DataBody()) {
		t.Log("success")
	} else {
		t.Failed()
	}
}

func TestDataManage_DataHandleAsync(t *testing.T) {
	// 每次测试清空数据
	abs := db.RunPath()
	path := filepath.Join(abs, "temp/data")

	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("err %v", err)
		return
	}

	// 数据管理
	dm0 := NewManage(nil, newOpt())
	t.Logf("%+v", dm0)

	// 模拟数据管理
	dm1 := newMockManage()
	t.Logf("%+v", dm1)

	// 开始执行并发测试

	num := 100   // 协程总数
	work := 1000 // 每个协程循环次数

	tid := tx.Super // 此时不测试事务, 因此使用超级事务
	id0s := make([]uint64, 0)
	id1s := make([]uint64, 0)

	lock := new(sync.Mutex)          // 初始化互斥锁
	waitGroup := new(sync.WaitGroup) // 初始化任务等待组

	worker := func() {
		dataLen := 80 // 随机测试数据长度
		defer waitGroup.Done()
		for i := 0; i < work; i++ {
			op := rand.Int() % 100
			if op < 50 { // Write
				data := randB(dataLen)
				id0, e := dm0.Write(tid, data)
				if e != nil {
					continue
				}
				id1, e := dm1.Write(tid, data)
				if e != nil {
					continue
				}

				lock.Lock()
				id0s = append(id0s, id0)
				id1s = append(id1s, id1)
				lock.Unlock()
			} else { // Read and Update
				lock.Lock()
				if len(id0s) == 0 {
					lock.Unlock()
					continue
				}

				// 随机选择一个id
				tmp := rand.Int() % len(id0s)
				id0 := id0s[tmp]
				id1 := id1s[tmp]
				lock.Unlock()

				data0, ok, e := dm0.Read(id0)
				if e != nil || ok == false {
					continue
				}
				data1, ok, e := dm1.Read(id1)
				if e != nil || ok == false {
					continue
				}

				data0.RLock()
				data1.RLock()
				if !bytes.Equal(data0.DataBody(), data1.DataBody()) {
					t.Logf("%+v", data0.DataBody())
					t.Logf("%+v", data1.DataBody())
					t.Fatalf("check error data not equal")
				}
				data0.RUnlock()
				data1.RUnlock()

				// 更新数据
				data := randB(dataLen)
				data0.Before()
				data1.Before()
				copy(data0.DataBody(), data)
				copy(data1.DataBody(), data)
				data0.After(tid)
				data1.After(tid)
				data0.Release()
				data1.Release()
			}
		}
	}

	waitGroup.Add(num)
	for i := 0; i < num; i++ {
		go worker()
	}
	waitGroup.Wait()
}
