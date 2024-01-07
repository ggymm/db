package page

import "sync"

type Page interface {
	Lock()
	Unlock()

	No() uint32
	Data() []byte
	Dirty()
	Release()
}

// 保存在内存中的页面结构
type page struct {
	lock sync.Mutex

	no    uint32 // 编号（从 1 开始）
	data  []byte // 数据内容
	dirty bool   // 是否是脏页面
	cache Cache  // 页面缓存
}

func NewPage(no uint32, data []byte, cache Cache) Page {
	return &page{
		no:    no,
		data:  data,
		cache: cache,
	}
}

func (p *page) Lock() {
	p.lock.Lock()
}

func (p *page) Unlock() {
	p.lock.Unlock()
}

func (p *page) No() uint32 {
	return p.no
}

func (p *page) Data() []byte {
	return p.data
}

func (p *page) Dirty() {
	p.dirty = true
}

func (p *page) Release() {
	// p.cache.Release()
}
