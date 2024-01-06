package txn

import (
	"db/pkg/utils"
	"os"
	"path/filepath"
	"sync"
)

const (
	Active    byte = 0 // 事务正在进行中
	Committed byte = 1 // 事务已经提交
	Aborted   byte = 2 // 事务已经终止

	Suffix = ".tid" // tid 文件后缀

	FieldLen  = 1       // 事务状态字段长度
	HeaderLen = TIDSize // TID 文件头长度
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
// tid 文件中每个事务使用 1 byte 存储其状态，位移为 (tid - 1) + HeaderLen
type txnManager struct {
	lock sync.Mutex

	tidSeq   TID
	tidFile  *os.File
	filename string
}

func pos(tid TID) int64 {
	return HeaderLen + int64(tid-1)*FieldLen
}

func open(tm *txnManager) {
	// 打开文件
	file, err := os.OpenFile(tm.filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// 解析文件
	buf := make([]byte, HeaderLen)
	_, err = file.ReadAt(buf, 0)
	if err != nil {
		panic(err)
	}
	tid := readTID(buf)

	// 获取 tid 对应的状态位置
	off := pos(tid)
	stat, _ := file.Stat()
	if off != stat.Size() {
		panic(ErrBadTIDFile)
	}

	// 构造结构体
	tm.tidSeq = tid
	tm.tidFile = file
}

func create(tm *txnManager) {
	filename := tm.filename

	// 判断文件夹是否存在
	dir := filepath.Dir(filename)
	if !utils.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	// 写入文件头
	buf := make([]byte, HeaderLen)
	writeTID(buf, 1) // tid 从 1 开始
	_, err = file.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}

	// 构造结构体
	tm.tidSeq = 1
	tm.tidFile = file
}

func NewTxnManager(filename string) Manager {
	tm := new(txnManager)
	tm.filename = filename + Suffix

	// 判断文件是否存在
	if utils.IsExist(tm.filename) {
		open(tm)
	} else {
		create(tm)
	}
	return tm
}

func (tm *txnManager) inc() {
	tm.tidSeq++
	buf := make([]byte, 8)
	writeTID(buf, tm.tidSeq)

	// 写入并同步文件
	var err error
	_, err = tm.tidFile.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}
	err = tm.tidFile.Sync()
	if err != nil {
		panic(err)
	}
}

func (tm *txnManager) sync(tid TID, state byte) {
	off := pos(tid)

	// 写入并同步文件
	var err error
	_, err = tm.tidFile.WriteAt([]byte{state}, off)
	if err != nil {
		return
	}
	err = tm.tidFile.Sync()
	if err != nil {
		panic(err)
	}
}

func (tm *txnManager) state(tid TID, state byte) bool {
	off := pos(tid)

	// 读取对应位置状态
	buf := make([]byte, 1)
	_, err := tm.tidFile.ReadAt(buf, off)
	if err != nil {
		panic(err)
	}
	return buf[0] == state
}

func (tm *txnManager) Close() {
	var err error
	err = tm.tidFile.Sync()
	if err != nil {
		panic(err)
	}
	err = tm.tidFile.Close()
	if err != nil {
		panic(err)
	}
}

func (tm *txnManager) Begin() TID {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	tid := tm.tidSeq
	tm.inc()
	tm.sync(tid, Active)
	return tid
}

func (tm *txnManager) Abort(tid TID) {
	tm.sync(tid, Aborted)
}

func (tm *txnManager) Commit(tid TID) {
	tm.sync(tid, Committed)
}

func (tm *txnManager) IsActive(tid TID) bool {
	if tid == Super {
		return false
	}
	return tm.state(tid, Active)
}

func (tm *txnManager) IsCommitted(tid TID) bool {
	if tid == Super {
		return true
	}
	return tm.state(tid, Committed)
}

func (tm *txnManager) IsAborted(tid TID) bool {
	if tid == Super {
		return false
	}
	return tm.state(tid, Aborted)
}
