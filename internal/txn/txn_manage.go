package txn

import (
	"db/pkg/utils"
	"os"
	"sync"
)

const (
	Active    = 0 // 事务正在进行中
	Committed = 1 // 事务已经提交
	Aborted   = 2 // 事务已经终止

	Suffix = ".tid" // tid 文件后缀

	TIDFieldLen  = 1 // 事务状态字段长度
	TIDHeaderLen = 8 // TID 文件头长度
)

type Manager interface {
	Close()         // 关闭事务管理器
	Begin() TID     // 开启一个事务
	Abort(tid TID)  // 取消一个事务
	Commit(tid TID) // 提交一个事务

	IsActive(tid TID) bool    // 判断事务是否处于进行中
	IsCommitted(tid TID) bool // 判断事务是否已经提交
	IsAborted(tid TID) bool   // 判断事务是否已经取消
}

// txnManager 用于管理 tid 文件
//
// 事务ID(tid) 起始为 1 按顺序递增
//
// 每个事务有三种状态
// 0. active     事务正在进行中
// 1. committed  事务已经提交
// 2. aborted    事务已经终止
//
// tid 文件起始位置为 8 type，存储 tid 序列号
// tid 文件中每个事务使用 1 byte 存储其状态，位移为 (tid - 1) + TIDHeaderLen
type txnManager struct {
	lock sync.Mutex

	tidSeq  TID
	tidFile *os.File
}

func NewTxnManager(filename string) Manager {
	tm := new(txnManager)
	// 判断文件是否存在
	if utils.IsExist(filename) {
		// 打开文件
		file, err := os.OpenFile(filename, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		tm.tidFile = file
	} else {
		// 创建文件
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		tm.tidFile = file
	}

	return tm
}

func (tm *txnManager) check() {

}

func (tm *txnManager) Close() {

}

func (tm *txnManager) Begin() TID {
	return 0
}

func (tm *txnManager) Abort(tid TID) {

}

func (tm *txnManager) Commit(tid TID) {

}

func (tm *txnManager) IsActive(tid TID) bool {
	return false
}

func (tm *txnManager) IsCommitted(tid TID) bool {
	return false
}

func (tm *txnManager) IsAborted(tid TID) bool {
	return false
}
