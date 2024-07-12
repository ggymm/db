package page

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/ggymm/db"
	"github.com/ggymm/db/pkg/cache"
	"github.com/ggymm/db/pkg/file"
)

var ErrMemoryNotEnough = errors.New("memory not enough")

const (
	name = "DB.BIN"

	Size  = 1 << 13 // 页面大小 8KB
	Limit = 10
)

type Manage interface {
	Close()

	NewPage(data []byte) uint32        // 创建新的页面，返回页面编号
	ObtainPage(n uint32) (Page, error) // 跟觉页号获取页面
	ReleasePage(p Page)                // 释放页面

	// 以下方法在 recovery 时使用

	PageNum() int              // 返回当前缓存的页面数量
	PageFlush(p Page)          // 刷新页面到磁盘
	PageTruncate(maxNo uint32) // 截断页面
}

type pageManage struct {
	mu sync.Mutex

	num   uint32      // page 总数
	file  *os.File    // 文件句柄
	cache cache.Cache // 缓存

	filepath string // 文件名称
}

func pos(no uint32) int64 {
	return int64(no-1) * Size
}

func open(m *pageManage) {
	// 打开文件
	f, err := os.OpenFile(m.filepath, os.O_RDWR, file.Mode)
	if err != nil {
		panic(err)
	}

	// 获取文件大小
	stat, _ := f.Stat()
	size := stat.Size()

	// 字段信息
	m.num = uint32(size / Size)
	m.file = f
}

func create(m *pageManage) {
	// 创建父文件夹
	dir := filepath.Dir(m.filepath)
	if !file.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	f, err := os.OpenFile(m.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode)
	if err != nil {
		panic(err)
	}

	// 字段信息
	m.num = 0
	m.file = f
}

func NewManage(opt *db.Option) Manage {
	if opt.Memory/Size < Limit {
		panic(ErrMemoryNotEnough)
	}
	m := new(pageManage)
	m.filepath = filepath.Join(opt.GetPath(name))

	// 构造缓存对象
	m.cache = cache.NewCache(&cache.Option{
		Obtain:   m.obtainForCache,
		Release:  m.releaseForCache,
		MaxCount: uint32(opt.Memory / Size),
	})

	// 判断文件是否存在
	if opt.Open {
		open(m)
	} else {
		create(m)
	}
	return m
}

// 将页内容刷新到磁盘
func (m *pageManage) flush(p Page) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 写入数据
	_, err := m.file.WriteAt(p.Data(), pos(p.No()))
	if err != nil {
		panic(err)
	}

	// 刷新文件
	err = m.file.Sync()
	if err != nil {
		panic(err)
	}
}

// obtainForCache 需要支持并发
// 缓存中不存在时，从磁盘中获取，并且包装成 Page 对象
func (m *pageManage) obtainForCache(key uint64) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	no := uint32(key)

	// 读取数据
	buf := make([]byte, Size)
	_, err := m.file.ReadAt(buf, pos(no))
	if err != nil {
		panic(err)
	}
	return NewPage(no, buf, m), nil
}

// releaseForCache 需要是同步方法
// 释放缓存，需要将 Page 对象内存刷新到磁盘
func (m *pageManage) releaseForCache(data any) {
	p := data.(Page)
	if p.Dirty() {
		m.flush(p)
		p.SetDirty(false)
	}
}

func (m *pageManage) Close() {
	m.cache.Close()
}

func (m *pageManage) NewPage(data []byte) uint32 {
	no := atomic.AddUint32(&m.num, 1)

	// 创建页面
	p := NewPage(no, data, m)
	m.flush(p)
	return no
}

func (m *pageManage) ObtainPage(n uint32) (Page, error) {
	data, err := m.cache.Obtain(uint64(n))
	if err != nil {
		return nil, err
	}
	return data.(Page), nil
}

func (m *pageManage) ReleasePage(p Page) {
	m.cache.Release(uint64(p.No()))
}

func (m *pageManage) PageNum() int {
	return int(m.num)
}

func (m *pageManage) PageFlush(p Page) {
	m.flush(p)
}

func (m *pageManage) PageTruncate(maxNo uint32) {
	size := pos(maxNo + 1)
	err := m.file.Truncate(size)
	if err != nil {
		panic(err)
	}
	m.num = maxNo
}
