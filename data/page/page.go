package page

import (
	"sync"
)

// 内存中页面的结构
//
// no 为编号，读取和写入时，通过 no * pageSize 计算偏移量
// data 是数据内容
// dirty 为 true 代表数据被修改过，需要写入文件

type Page interface {
	Lock()
	Unlock()

	No() uint32
	Data() []byte
	Dirty() bool

	SetNo(no uint32)
	SetData(data []byte)
	SetDirty(dirty bool)

	Release()
}

type page struct {
	mu sync.Mutex

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
	p.mu.Lock()
}

func (p *page) Unlock() {
	p.mu.Unlock()
}

func (p *page) No() uint32 {
	return p.no
}

func (p *page) Data() []byte {
	return p.data
}

func (p *page) Dirty() bool {
	return p.dirty
}

func (p *page) SetNo(no uint32) {
	p.no = no
}

func (p *page) SetData(data []byte) {
	p.data = data
}

func (p *page) SetDirty(dirty bool) {
	p.dirty = dirty
}

func (p *page) Release() {
	p.cache.ReleasePage(p)
}
