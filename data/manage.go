package data

import (
	"errors"

	"github.com/ggymm/db"
	"github.com/ggymm/db/data/log"
	"github.com/ggymm/db/data/page"
	"github.com/ggymm/db/pkg/cache"
	"github.com/ggymm/db/tx"
)

// 数据管理器（磁盘管理器）
//
// 负责管理数据对象（dataItem）读取和写入
// 数据对象的编辑，通过调用数据对象的 Before 和 After 方法实现
//
// 数据管理（dataItem），与 pageManage 都会缓存数据
// dataManage 是一级缓存，缓存的是 dataItem 对象
// pageManage 是二级缓存，缓存的是 page 对象

var (
	ErrBusy         = errors.New("database is busy")
	ErrInitPage1    = errors.New("init page1 error")
	ErrDataTooLarge = errors.New("data too large")
)

const (
	maxTry = 5
)

type Manage interface {
	Close()

	Read(id uint64) (Item, bool, error)
	Write(tid uint64, data []byte) (uint64, error)

	LogDataItem(tid uint64, item Item)
	ReleaseDataItem(item Item)

	TxManage() tx.Manage
}

type dataManage struct {
	txManage tx.Manage // 用于 recover 操作

	log        log.Log     // 日志
	page1      page.Page   // page1
	pageIndex  page.Index  // page 索引
	pageManage page.Manage // page 管理

	cache cache.Cache // item 缓存
}

func open(m *dataManage) {
	// 读取 page1
	page1, err := m.pageManage.ObtainPage(1)
	if err != nil {
		panic(err)
	}
	m.page1 = page1

	// 根据 page1 校验数据
	if page.CheckVc(m.page1) == false {
		// 执行恢复操作
		// TODO
	}

	// 读取 page 数据，填充 pageIndex
	num := m.pageManage.PageNum()
	for i := 2; i <= num; i++ {
		p, e := m.pageManage.ObtainPage(uint32(i))
		if e != nil {
			panic(e)
		}
		m.pageIndex.Add(p.No(), page.CalcPageFree(p))
		p.Release()
	}

	// 重新设置 page1 的校验数据
	page.SetVcOpen(m.page1)
	m.pageManage.PageFlush(m.page1)
}

func create(m *dataManage) {
	// 创建 page1
	no := m.pageManage.NewPage(page.NewPage1())
	if no != 1 {
		panic(ErrInitPage1)
	}
	page1, err := m.pageManage.ObtainPage(no)
	if err != nil {
		panic(err)
	}

	// 刷新 page1 数据到磁盘
	m.page1 = page1
	m.pageManage.PageFlush(m.page1)
}

func NewManage(tm tx.Manage, opt *db.Option) Manage {
	m := new(dataManage)

	m.txManage = tm

	m.log = log.NewLog(opt)
	m.pageIndex = page.NewIndex()
	m.pageManage = page.NewManage(opt)

	m.cache = cache.NewCache(&cache.Option{
		Obtain:   m.obtainForCache,
		Release:  m.releaseForCache,
		MaxCount: 0,
	})

	if opt.Open {
		open(m)
	} else {
		create(m)
	}
	return m
}

// obtainForCache 需要支持并发
//
// 当 dataManage 缓存中没有数据时，需要从 pageManage 缓存中获取数据
// 此时若 pageManage 缓存中也没有数据，则会从磁盘加载数据
func (m *dataManage) obtainForCache(key uint64) (any, error) {
	no, off := parseDataItemId(key)
	p, err := m.pageManage.ObtainPage(no)
	if err != nil {
		return nil, err
	}
	return parseDataItem(p, off, m), nil
}

// releaseForCache 需要是同步方法
func (m *dataManage) releaseForCache(data any) {
	item := data.(Item)
	item.Page().Release()
}

func (m *dataManage) Close() {
	m.log.Close()
	m.cache.Close()

	page.SetVcClose(m.page1) // 设置 page1 的校验数据
	m.page1.Release()
	m.pageManage.Close()
}

// Read 读取数据对象，从缓存中读取数据对象
// 若缓存中没有数据对象，则从 pageManage 缓存中读取数据对象
func (m *dataManage) Read(id uint64) (Item, bool, error) {
	data, err := m.cache.Obtain(id)
	if err != nil {
		return nil, false, err
	}
	item := data.(Item)
	if !item.Flag() {
		item.Release()
		return nil, false, nil
	}
	return item, true, nil
}

func (m *dataManage) Write(tid uint64, data []byte) (uint64, error) {
	data = wrapDataItem(data)
	length := uint32(len(data))
	if length > page.MaxPageFree() {
		return 0, ErrDataTooLarge
	}

	var (
		err  error
		p    page.Page
		no   uint32
		free uint32
	)

	// 选择可以插入的 page
	for i := 0; i < maxTry; i++ {
		no, free = m.pageIndex.Select(length)
		if free > 0 {
			break
		} else {
			// 创建新页，等待下次选择
			newNo := m.pageManage.NewPage(page.NewPageX())
			m.pageIndex.Add(newNo, page.MaxPageFree())
		}
	}
	if no == 0 {
		return 0, ErrBusy
	}
	defer func() {
		if p == nil {
			m.pageIndex.Add(no, free)
		} else {
			m.pageIndex.Add(no, page.CalcPageFree(p))
		}
	}()

	// 获取 page
	p, err = m.pageManage.ObtainPage(no)
	if err != nil {
		return 0, err
	}

	// 保存日志
	m.log.Log(wrapInsertLog(tid, p, data))

	// 保存数据
	off := page.WritePageData(p, data)

	// 释放页面
	p.Release()

	// 返回 item_id
	return wrapDataItemId(no, off), nil
}

func (m *dataManage) LogDataItem(tid uint64, item Item) {
	// 包装 update log 数据
	data := wrapUpdateLog(tid, item)
	m.log.Log(data)
}

func (m *dataManage) ReleaseDataItem(item Item) {
	m.cache.Release(item.Id())
}

func (m *dataManage) TxManage() tx.Manage {
	return m.txManage
}
