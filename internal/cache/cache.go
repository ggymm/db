package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCacheFull = errors.New("cache is full")
)

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
	lock sync.Mutex

	ops   *Option
	count uint32

	refs   map[uint64]uint32
	cache  map[uint64]any
	obtain map[uint64]bool
}

func NewCache(ops *Option) Cache {
	return &cache{
		ops:    ops,
		refs:   make(map[uint64]uint32),
		cache:  make(map[uint64]any),
		obtain: make(map[uint64]bool),
	}
}

func (c *cache) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.cache {
		c.ops.Release(v)
		delete(c.refs, k)
		delete(c.cache, k)
	}
}

func (c *cache) Obtain(key uint64) (any, error) {
	for {
		c.lock.Lock()

		// 如果正在被其他线程获取，则等待
		if _, ok := c.obtain[key]; ok {
			c.lock.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}

		// 如果在缓存中，则直接返回
		if data, ok := c.cache[key]; ok {
			c.refs[key]++
			c.lock.Unlock()
			return data, nil
		}

		// 如果缓存中的数据已经达到上限，则抛出异常
		if c.ops.MaxCount > 0 && c.count >= c.ops.MaxCount {
			c.lock.Unlock()
			return nil, ErrCacheFull
		}

		// 尝试获取
		c.count++
		c.obtain[key] = true
		c.lock.Unlock()
		break
	}

	// 获取数据
	data, err := c.ops.Obtain(key)
	if err != nil {
		c.lock.Lock()
		c.count--
		delete(c.obtain, key)
		c.lock.Unlock()
		return nil, err
	}

	// 缓存该数据
	c.lock.Lock()
	c.refs[key] = 1
	c.cache[key] = data
	delete(c.obtain, key)
	c.lock.Unlock()
	return data, nil
}

func (c *cache) Release(key uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.refs[key]--
	if c.refs[key] == 0 {
		c.ops.Release(c.cache[key])
		delete(c.refs, key)
		delete(c.cache, key)
		c.count--
	}
}
