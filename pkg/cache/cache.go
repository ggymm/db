package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrCacheFull = errors.New("cache is full")

type Cache interface {
	Close()

	Obtain(key uint64) (any, error)
	Release(key uint64)
}

type Option struct {
	// 当前缓存不存在时，进行获取操作
	// 该函数必须是并发安全的
	Obtain func(key uint64) (any, error)

	// 当前缓存被释放时，进行释放操作
	// 该函数必须是同步方法
	Release func(data any)

	// 当前缓存的最大数量，如果小于等于 0，则不限制
	MaxCount uint32
}

type cache struct {
	sync.Mutex

	opt   *Option
	count uint32

	refs   map[uint64]uint32
	cache  map[uint64]any
	obtain map[uint64]bool
}

func NewCache(opt *Option) Cache {
	return &cache{
		opt:    opt,
		refs:   make(map[uint64]uint32),
		cache:  make(map[uint64]any),
		obtain: make(map[uint64]bool),
	}
}

func (c *cache) Close() {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.cache {
		c.opt.Release(v)
		delete(c.refs, k)
		delete(c.cache, k)
	}
}

func (c *cache) Obtain(key uint64) (any, error) {
	for {
		c.Lock()

		// 如果正在被其他线程获取，则等待
		if _, ok := c.obtain[key]; ok {
			c.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}

		// 如果在缓存中，则直接返回
		if data, ok := c.cache[key]; ok {
			c.refs[key]++
			c.Unlock()
			return data, nil
		}

		// 如果缓存中的数据已经达到上限，则抛出异常
		if c.opt.MaxCount > 0 && c.count >= c.opt.MaxCount {
			c.Unlock()
			return nil, ErrCacheFull
		}

		// 尝试获取
		c.count++
		c.obtain[key] = true
		c.Unlock()
		break
	}

	// 获取数据
	data, err := c.opt.Obtain(key)
	if err != nil {
		c.Lock()
		c.count--
		delete(c.obtain, key)
		c.Unlock()
		return nil, err
	}

	// 缓存该数据
	c.Lock()
	c.refs[key] = 1
	c.cache[key] = data
	delete(c.obtain, key)
	c.Unlock()
	return data, nil
}

func (c *cache) Release(key uint64) {
	c.Lock()
	defer c.Unlock()

	c.refs[key]--
	if c.refs[key] == 0 {
		c.opt.Release(c.cache[key])
		delete(c.refs, key)
		delete(c.cache, key)
		c.count--
	}
}
