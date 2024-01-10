package data

import (
	"db/internal/cache"
	"db/internal/data/log"
	"db/internal/data/page"
	"db/internal/ops"
	"db/internal/txn"
	"errors"
)

// 数据管理器
//
// 负责管理数据对象（dataItem）读取和写入
// 数据对象的编辑，通过调用数据对象的 Before 和 After 方法实现
//
// 数据管理（dataItem），与 pageCache 都会缓存数据
// dataManage 是一级缓存，缓存的是 dataItem 对象
// pageCache 是二级缓存，缓存的是 page 对象

var (
	ErrInitPage1 = errors.New("init page1 error")
)

type Manage interface {
	Close()

	Read(id uint64) (Item, bool, error)
	Insert(tid txn.TID, data []byte) (uint64, error)

	LogDataItem(tid txn.TID, item Item)
	ReleaseDataItem(item Item)
}

type dataManage struct {
	log       log.Log
	page1     page.Page
	txnManage txn.Manage

	pageIndex page.Index
	pageCache page.Cache

	cache cache.Cache
}

func open(dm *dataManage) {
	// 读取 page1
	page1, err := dm.pageCache.ObtainPage(1)
	if err != nil {
		panic(err)
	}
	dm.page1 = page1

	// 根据 page1 校验数据
	if CheckVc(dm.page1) == false {
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
		dm.pageIndex.Add(p.No(), ParseFree(p))
		p.Release()
	}

	// 重新设置 page1 的校验数据
	SetVcOpen(dm.page1)
	dm.pageCache.PageFlush(dm.page1)
}

func create(dm *dataManage) {
	// 创建 page1
	no := dm.pageCache.NewPage(InitPage1())
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

func NewManage(ops *ops.Option, txn txn.Manage) Manage {
	dm := new(dataManage)

	dm.log = log.NewLog(ops)
	dm.txnManage = txn

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

	SetVcClose(dm.page1) // 设置 page1 的校验数据
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
	if !item.Valid() {
		item.Release()
		return nil, false, nil
	}
	return item, true, nil
}

func (dm *dataManage) Insert(tid txn.TID, data []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (dm *dataManage) LogDataItem(tid txn.TID, item Item) {
	// 包装 update log 数据
	// dm.log.Log()
}

func (dm *dataManage) ReleaseDataItem(item Item) {
	dm.cache.Release(item.Id())
}
