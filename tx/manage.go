package tx

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/ggymm/db"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/file"
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
// 事务Id(tid) 起始为 1 按顺序递增
//
// 每个事务有三种状态
//         0. active     事务正在进行中
//         1. committed  事务已经提交
//         2. aborted    事务已经终止

var (
	ErrBadIdFile = errors.New("bad Id File")
)

const (
	name = "DB.TXN" // 事务 id 文件

	Super  uint64 = 0
	TidLen        = 8

	Active    byte = 0 // 事务正在进行中
	Committed byte = 1 // 事务已经提交
	Aborted   byte = 2 // 事务已经终止

	fieldLen  = 1      // 事务状态字段长度
	headerLen = TidLen // 文件头长度
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
	mu sync.Mutex

	seq  uint64   // 当前事务Id
	file *os.File // 文件句柄

	filepath string // 文件名称
}

func pos(tid uint64) int64 {
	return headerLen + int64(tid-1)*fieldLen
}

func open(m *txManager) {
	// 打开文件
	f, err := os.OpenFile(m.filepath, os.O_RDWR, file.Mode)
	if err != nil {
		panic(err)
	}

	// 解析文件
	buf := make([]byte, headerLen)
	_, err = f.ReadAt(buf, 0)
	if err != nil {
		panic(err)
	}
	tid := bin.Uint64(buf)

	// 获取 tid 对应的状态位置
	off := pos(tid)
	stat, _ := f.Stat()
	if off != stat.Size() {
		panic(ErrBadIdFile)
	}

	// 字段信息
	m.seq = tid
	m.file = f
}

func create(m *txManager) {
	p := m.filepath

	// 创建父文件夹
	dir := filepath.Dir(p)
	if !file.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode)
	if err != nil {
		panic(err)
	}

	// 写入文件头
	buf := make([]byte, headerLen)
	bin.PutUint64(buf, 1) // tid 从 1 开始
	_, err = f.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}

	// 字段信息
	m.seq = 1
	m.file = f
}

func NewManager(opt *db.Option) Manage {
	m := new(txManager)
	m.filepath = filepath.Join(opt.GetPath(name))

	// 判断文件是否存在
	if opt.Open {
		open(m)
	} else {
		create(m)
	}
	return m
}

func (m *txManager) inc() {
	m.seq++
	buf := make([]byte, 8)
	bin.PutUint64(buf, m.seq)

	// 写入并同步文件
	m.write(buf, 0)
}

func (m *txManager) state(tid uint64) byte {
	off := pos(tid)

	// 读取对应位置状态
	buf := make([]byte, 1)
	_, err := m.file.ReadAt(buf, off)
	if err != nil {
		panic(err)
	}
	return buf[0]
}

func (m *txManager) update(tid uint64, state byte) {
	off := pos(tid)

	// 写入并同步文件
	m.write([]byte{state}, off)
}

func (m *txManager) write(buf []byte, off int64) {
	_, err := m.file.WriteAt(buf, off)
	if err != nil {
		panic(err)
	}
	err = m.file.Sync()
	if err != nil {
		panic(err)
	}
}

func (m *txManager) Close() {
	err := m.file.Close()
	if err != nil {
		panic(err)
	}
}

func (m *txManager) Begin() uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	tid := m.seq
	m.inc()
	m.update(tid, Active)
	return tid
}

func (m *txManager) Abort(tid uint64) {
	m.update(tid, Aborted)
}

func (m *txManager) Commit(tid uint64) {
	m.update(tid, Committed)
}

func (m *txManager) IsActive(tid uint64) bool {
	if tid == Super {
		return false
	}
	return m.state(tid) == Active
}

func (m *txManager) IsCommitted(tid uint64) bool {
	if tid == Super {
		return true
	}
	return m.state(tid) == Committed
}

func (m *txManager) IsAborted(tid uint64) bool {
	if tid == Super {
		return false
	}
	return m.state(tid) == Aborted
}
