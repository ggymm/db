package page

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"db/internal/cache"
	"db/pkg/utils"
)

var (
	ErrMemoryNotEnough = errors.New("memory not enough")
)

const (
	Size  = 1 << 13 // 页面大小 8KB
	Limit = 10

	Suffix = ".db"
)

type Cache interface {
	Close()

	NewPage(data []byte) uint32      // 创建新的页面，返回页面编号
	GetPage(no uint32) (Page, error) // 跟觉页号获取页面

	// 以下方法在 recovery 时使用

	PageNum() uint32           // 返回当前缓存的页面数量
	PageFlush(p Page)          // 刷新页面到磁盘
	PageTruncate(maxNo uint32) // TODO
}

type pageCache struct {
	lock sync.Mutex

	no       uint32
	file     *os.File
	memory   int64
	filename string

	cache cache.Cache
}

func open(pc *pageCache) {
	// 打开文件
	file, err := os.OpenFile(pc.filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	// 字段信息
	pc.no = uint32(size / Size)
	pc.file = file
}

func create(pc *pageCache) {
	filename := pc.filename

	// 创建父文件夹
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

	// 字段信息
	pc.no = 0
	pc.file = file
}

func NewCache(memory int64, filename string) Cache {
	if memory/Size < Limit {
		panic(ErrMemoryNotEnough)
	}
	pc := new(pageCache)
	pc.memory = memory
	pc.filename = filename + Suffix

	// 构造缓存对象
	pc.cache = cache.NewCache(&cache.Option{
		Obtain:   pc.obtainForCache,
		Release:  pc.releaseForCache,
		MaxCount: uint32(memory / Size),
	})

	// 判断文件是否存在
	if utils.IsExist(pc.filename) {
		open(pc)
	} else {
		create(pc)
	}
	return pc
}

func pos(no uint32) int64 {
	return int64(no-1) * Size
}

// 将页内容刷新到磁盘
func (c *pageCache) flush(p Page) {
	c.lock.Lock()
	defer c.lock.Unlock()

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

func (c *pageCache) release(p Page) {
	c.cache.Release(uint64(p.No()))
}

// Obtain 需要支持并发
// 缓存中不存在时，从磁盘中获取，并且包装成 Page 对象
func (c *pageCache) obtainForCache(key uint64) (any, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	no := uint32(key)

	// 读取数据
	buf := make([]byte, Size)
	_, err := c.file.ReadAt(buf, pos(no))
	if err != nil {
		panic(err)
	}
	return NewPage(no, buf, c), nil
}

// Release 需要是同步方法
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
	no := atomic.AddUint32(&c.no, 1)

	// 创建页面
	p := NewPage(no, data, c)
	c.flush(p)
	return no
}

func (c *pageCache) GetPage(no uint32) (Page, error) {
	data, err := c.cache.Obtain(uint64(no))
	if err != nil {
		return nil, err
	}
	return data.(Page), nil
}

func (c *pageCache) PageNum() uint32 {
	return c.no
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
	c.no = maxNo
}
