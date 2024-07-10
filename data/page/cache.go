package page

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"db"
	"db/pkg/cache"
	"db/pkg/file"
)

var ErrMemoryNotEnough = errors.New("memory not enough")

const (
	Size  = 1 << 13 // 页面大小 8KB
	Limit = 10

	suffix = ".data"
)

type Cache interface {
	Close()

	NewPage(data []byte) uint32        // 创建新的页面，返回页面编号
	ObtainPage(n uint32) (Page, error) // 跟觉页号获取页面
	ReleasePage(p Page)                // 释放页面

	// 以下方法在 recovery 时使用

	PageNum() int              // 返回当前缓存的页面数量
	PageFlush(p Page)          // 刷新页面到磁盘
	PageTruncate(maxNo uint32) // 截断页面
}

type pageCache struct {
	mu sync.Mutex

	num   uint32      // page 总数
	file  *os.File    // 文件句柄
	cache cache.Cache // 缓存

	filepath string // 文件名称
}

func pos(no uint32) int64 {
	return int64(no-1) * Size
}

func open(c *pageCache) {
	// 打开文件
	f, err := os.OpenFile(c.filepath, os.O_RDWR, file.Mode)
	if err != nil {
		panic(err)
	}

	// 获取文件大小
	stat, _ := f.Stat()
	size := stat.Size()

	// 字段信息
	c.num = uint32(size / Size)
	c.file = f
}

func create(c *pageCache) {
	p := c.filepath

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

	// 字段信息
	c.num = 0
	c.file = f
}

func NewCache(opt *db.Option) Cache {
	if opt.Memory/Size < Limit {
		panic(ErrMemoryNotEnough)
	}
	c := new(pageCache)
	c.filepath = filepath.Join(opt.GetPath(suffix))

	// 构造缓存对象
	c.cache = cache.NewCache(&cache.Option{
		Obtain:   c.obtainForCache,
		Release:  c.releaseForCache,
		MaxCount: uint32(opt.Memory / Size),
	})

	// 判断文件是否存在
	if opt.Open {
		open(c)
	} else {
		create(c)
	}
	return c
}

// 将页内容刷新到磁盘
func (c *pageCache) flush(p Page) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 写入数据
	_, err := c.file.WriteAt(p.Data(), pos(p.No()))
	if err != nil {
		panic(err)
	}

	// 刷新文件
	err = c.file.Sync()
	if err != nil {
		panic(err)
	}
}

// obtainForCache 需要支持并发
// 缓存中不存在时，从磁盘中获取，并且包装成 Page 对象
func (c *pageCache) obtainForCache(key uint64) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	no := uint32(key)

	// 读取数据
	buf := make([]byte, Size)
	_, err := c.file.ReadAt(buf, pos(no))
	if err != nil {
		panic(err)
	}
	return NewPage(no, buf, c), nil
}

// releaseForCache 需要是同步方法
// 释放缓存，需要将 Page 对象内存刷新到磁盘
func (c *pageCache) releaseForCache(data any) {
	p := data.(Page)
	if p.Dirty() {
		c.flush(p)
		p.SetDirty(false)
	}
}

func (c *pageCache) Close() {
	c.cache.Close()
}

func (c *pageCache) NewPage(data []byte) uint32 {
	no := atomic.AddUint32(&c.num, 1)

	// 创建页面
	p := NewPage(no, data, c)
	c.flush(p)
	return no
}

func (c *pageCache) ObtainPage(n uint32) (Page, error) {
	data, err := c.cache.Obtain(uint64(n))
	if err != nil {
		return nil, err
	}
	return data.(Page), nil
}

func (c *pageCache) ReleasePage(p Page) {
	c.cache.Release(uint64(p.No()))
}

func (c *pageCache) PageNum() int {
	return int(c.num)
}

func (c *pageCache) PageFlush(p Page) {
	c.flush(p)
}

func (c *pageCache) PageTruncate(maxNo uint32) {
	size := pos(maxNo + 1)
	err := c.file.Truncate(size)
	if err != nil {
		panic(err)
	}
	c.num = maxNo
}
