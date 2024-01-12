package tx

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"db/internal/ops"
	"db/pkg/utils"
)

// 事务管理器
//
// 事务管理器文件结构如下：
// +----------------+----------------+----------------+----------------+
// |      type      |      tid       |      tid       |      tid       |
// +----------------+----------------+----------------+----------------+
// |     8 byte     |    1 byte      |    1 byte      |    1 byte      |
// +----------------+----------------+----------------+----------------+
//
// 事务ID(tid) 起始为 1 按顺序递增
//
// 每个事务有三种状态
//         0. active     事务正在进行中
//         1. committed  事务已经提交
//         2. aborted    事务已经终止

var (
	ErrBadTIDFile = errors.New("bad TID File")
)

const (
	Active    byte = 0 // 事务正在进行中
	Committed byte = 1 // 事务已经提交
	Aborted   byte = 2 // 事务已经终止

	suffix = ".tid" // tid 文件后缀

	fieldLen  = 1      // 事务状态字段长度
	headerLen = TIDLen // 文件头长度
)

type Manage interface {
	Close() // 关闭事务管理器

	Begin() uint64     // 开启一个事务
	Abort(tid uint64)  // 取消一个事务
	Commit(tid uint64) // 提交一个事务

	IsActive(tid uint64) bool    // 判断事务是否处于进行中
	IsCommitted(tid uint64) bool // 判断事务是否已经提交
	IsAborted(tid uint64) bool   // 判断事务是否已经取消
}

type txManager struct {
	lock sync.Mutex

	seq  uint64   // 当前事务ID
	file *os.File // 文件句柄

	filepath string // 文件名称
}

func pos(tid uint64) int64 {
	return headerLen + int64(tid-1)*fieldLen
}

func open(tm *txManager) {
	// 打开文件
	file, err := os.OpenFile(tm.filepath, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// 解析文件
	buf := make([]byte, headerLen)
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

	// 字段信息
	tm.seq = tid
	tm.file = file
}

func create(tm *txManager) {
	path := tm.filepath

	// 创建父文件夹
	dir := filepath.Dir(path)
	if !utils.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	// 写入文件头
	buf := make([]byte, headerLen)
	writeTID(buf, 1) // tid 从 1 开始
	_, err = file.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}

	// 字段信息
	tm.seq = 1
	tm.file = file
}

func NewManager(ops *ops.Option) Manage {
	tm := new(txManager)
	tm.filepath = filepath.Join(ops.Path, suffix)

	// 判断文件是否存在
	if ops.Open {
		open(tm)
	} else {
		create(tm)
	}
	return tm
}

func (tm *txManager) incr() {
	tm.seq++
	buf := make([]byte, 8)
	writeTID(buf, tm.seq)

	// 写入并同步文件
	tm.write(buf, 0)
}

func (tm *txManager) state(tid uint64) byte {
	off := pos(tid)

	// 读取对应位置状态
	buf := make([]byte, 1)
	_, err := tm.file.ReadAt(buf, off)
	if err != nil {
		panic(err)
	}
	return buf[0]
}

func (tm *txManager) update(tid uint64, state byte) {
	off := pos(tid)

	// 写入并同步文件
	tm.write([]byte{state}, off)
}

func (tm *txManager) write(buf []byte, off int64) {
	_, err := tm.file.WriteAt(buf, off)
	if err != nil {
		panic(err)
	}
	err = tm.file.Sync()
	if err != nil {
		panic(err)
	}
}

func (tm *txManager) Close() {
	err := tm.file.Close()
	if err != nil {
		panic(err)
	}
}

func (tm *txManager) Begin() uint64 {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	tid := tm.seq
	tm.incr()
	tm.update(tid, Active)
	return tid
}

func (tm *txManager) Abort(tid uint64) {
	tm.update(tid, Aborted)
}

func (tm *txManager) Commit(tid uint64) {
	tm.update(tid, Committed)
}

func (tm *txManager) IsActive(tid uint64) bool {
	if tid == Super {
		return false
	}
	return tm.state(tid) == Active
}

func (tm *txManager) IsCommitted(tid uint64) bool {
	if tid == Super {
		return true
	}
	return tm.state(tid) == Committed
}

func (tm *txManager) IsAborted(tid uint64) bool {
	if tid == Super {
		return false
	}
	return tm.state(tid) == Aborted
}
