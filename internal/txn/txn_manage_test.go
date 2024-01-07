package txn

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"db/pkg/utils"
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

// TestTxnManager_StageSync 用于测试事务管理器在并发环境下的行为
func TestTxnManager_StageSync(t *testing.T) {
	base := utils.RunPath()
	filename := filepath.Join(base, "temp/txn/",
		strconv.FormatInt(time.Now().UnixMilli(), 10))

	tm := NewTxnManager(filename)
	t.Log(base)
	t.Logf("%+v", tm)

	num := 50                        // 协程总数
	works := 3000                    // 每个协程循环次数
	curr := 0                        // 当前事务数目
	temp := make(map[TID]byte)       // 事务状态映射
	lock := new(sync.Mutex)          // 初始化互斥锁
	waitGroup := new(sync.WaitGroup) // 初始化任务等待组
	worker := func() {
		var (
			tid     TID
			isBegin = false
		)
		for i := 0; i < works; i++ {
			op := rand.Int() % 6
			if op == 0 {
				lock.Lock()
				if !isBegin {
					// 如果没有在事务中
					// 开始一个新的事务
					tid = tm.Begin()
					curr++
					temp[tid] = 0 // 保存事务状态
					isBegin = true
				} else {
					// 如果已经在事务中
					// 随机提交或者回滚
					state := (rand.Int() % 2) + 1
					switch state {
					case 1:
						tm.Commit(tid)
					case 2:
						tm.Abort(tid)
					}
					temp[tid] = byte(state) // 更新事务状态
					isBegin = false
				}
				lock.Unlock()
			} else {
				lock.Lock()
				// 如果有活跃的事务，进行验证
				if curr > 0 {
					tid = TID((rand.Int() % curr) + 1)
					state := temp[tid]
					var ok bool
					switch state {
					case 0:
						ok = tm.IsActive(tid)
					case 1:
						ok = tm.IsCommitted(tid)
					case 2:
						ok = tm.IsAborted(tid)
					}
					if ok {
						t.Log("check ok")
					} else {
						t.Log("check error")
						t.Failed()
					}
				}
				lock.Unlock()
			}
		}
		waitGroup.Done()
	}

	waitGroup.Add(num)
	for i := 0; i < num; i++ {
		go worker()
	}
	waitGroup.Wait()
	tm.Close()
}
