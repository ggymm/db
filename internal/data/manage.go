package data

import (
	"errors"

	"db/internal/data/log"
	"db/internal/data/page"
	"db/internal/ops"
	"db/internal/tx"
	"db/pkg/cache"
)

// 数据管理器（磁盘管理器）
//
// 负责管理数据对象（dataItem）读取和写入
// 数据对象的编辑，通过调用数据对象的 Before 和 After 方法实现
//
// 数据管理（dataItem），与 pageCache 都会缓存数据
// dataManage 是一级缓存，缓存的是 dataItem 对象
// pageCache 是二级缓存，缓存的是 page 对象

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
	Insert(tid uint64, data []byte) (uint64, error)

	LogDataItem(tid uint64, item Item)
	ReleaseDataItem(item Item)
}

type dataManage struct {
	txManage tx.Manage // 用于 recover 操作

	log       log.Log    // 日志
	page1     page.Page  // page1
	pageIndex page.Index // page 索引
	pageCache page.Cache // page 缓存

	cache cache.Cache // item 缓存
}

func open(dm *dataManage) {
	// 读取 page1
	page1, err := dm.pageCache.ObtainPage(1)
	if err != nil {
		panic(err)
	}
	dm.page1 = page1

	// 根据 page1 校验数据
	if page.CheckVc(dm.page1) == false {
		// 执行恢复操作
		// TODO
	}

	// 读取 page 数据，填充 pageIndex
	num := dm.pageCache.PageNum()
	for i := 2; i <= num; i++ {
		p, e := dm.pageCache.ObtainPage(uint32(i))
		if e != nil {
			panic(e)
		}
		dm.pageIndex.Add(p.No(), page.CalcPageFree(p))
		p.Release()
	}

	// 重新设置 page1 的校验数据
	page.SetVcOpen(dm.page1)
	dm.pageCache.PageFlush(dm.page1)
}

func create(dm *dataManage) {
	// 创建 page1
	no := dm.pageCache.NewPage(page.InitPage1())
	if no != 1 {
		panic(ErrInitPage1)
	}
	page1, err := dm.pageCache.ObtainPage(no)
	if err != nil {
		panic(err)
	}

	// 刷新 page1 数据到磁盘
	dm.page1 = page1
	dm.pageCache.PageFlush(dm.page1)
}

func NewManage(ops *ops.Option, txm tx.Manage) Manage {
	dm := new(dataManage)

	dm.txManage = txm

	dm.log = log.NewLog(ops)
	dm.pageIndex = page.NewIndex()
	dm.pageCache = page.NewCache(ops)

	dm.cache = cache.NewCache(&cache.Option{
		Obtain:   dm.obtainForCache,
		Release:  dm.releaseForCache,
		MaxCount: 0,
	})

	if ops.Open {
		open(dm)
	} else {
		create(dm)
	}
	return dm
}

// obtainForCache 需要支持并发
//
// 当 dataManage 缓存中没有数据时，需要从 pageCache 缓存中获取数据
// 此时若 pageCache 缓存中也没有数据，则会从磁盘加载数据
func (dm *dataManage) obtainForCache(key uint64) (any, error) {
	no, off := parseDataItemId(key)
	p, err := dm.pageCache.ObtainPage(no)
	if err != nil {
		return nil, err
	}
	return parseDataItem(p, off, dm), nil
}

// releaseForCache 需要是同步方法
func (dm *dataManage) releaseForCache(data any) {
	item := data.(Item)
	item.Page().Release()
}

func (dm *dataManage) Close() {
	dm.log.Close()
	dm.cache.Close()

	page.SetVcClose(dm.page1) // 设置 page1 的校验数据
	dm.page1.Release()
	dm.pageCache.Close()
}

// Read 读取数据对象，从缓存中读取数据对象
// 若缓存中没有数据对象，则从 pageCache 缓存中读取数据对象
func (dm *dataManage) Read(id uint64) (Item, bool, error) {
	data, err := dm.cache.Obtain(id)
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

func (dm *dataManage) Insert(tid uint64, data []byte) (uint64, error) {
	data = wrapDataItem(data)
	if len(data) > page.MaxPageFree() {
		return 0, ErrDataTooLarge
	}

	var (
		err  error
		p    page.Page
		no   uint32
		free int
	)

	// 选择可以插入的 page
	for i := 0; i < maxTry; i++ {
		no, free = dm.pageIndex.Select(len(data))
		if free > 0 {
			break
		} else {
			// 创建新页，等待下次选择
			newNo := dm.pageCache.NewPage(page.InitPageX())
			dm.pageIndex.Add(newNo, page.MaxPageFree())
		}
	}
	if no == 0 {
		return 0, ErrBusy
	}
	defer func() {
		if p == nil {
			dm.pageIndex.Add(no, free)
		} else {
			dm.pageIndex.Add(no, page.CalcPageFree(p))
		}
	}()

	// 获取 page
	p, err = dm.pageCache.ObtainPage(no)
	if err != nil {
		return 0, err
	}

	// 保存日志
	dm.log.Log(wrapInsertLog(tid, p, data))

	// 保存数据
	off := page.InsertPageData(p, data)

	// 释放页面
	p.Release()

	// 返回 item_id
	return wrapDataItemId(no, off), nil
}

func (dm *dataManage) LogDataItem(tid uint64, item Item) {
	// 包装 update log 数据
	data := wrapUpdateLog(tid, item)
	dm.log.Log(data)
}

func (dm *dataManage) ReleaseDataItem(item Item) {
	dm.cache.Release(item.Id())
}
