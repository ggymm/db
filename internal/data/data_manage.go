package data

import (
	"db/internal/cache"
	"db/internal/data/log"
	"db/internal/data/page"
	"db/internal/txn"
	"errors"
)

var (
	ErrInitPage1 = errors.New("init page1 error")
)

type Manage interface {
	Close()

	Read(id uint64) (Item, bool, error)
	Insert(tid txn.TID, data []byte) (uint64, error)
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

}

func create(dm *dataManage) {
	// 创建 page1
	no := dm.pageCache.NewPage(InitPage1())
	if no != 1 {
		panic(ErrInitPage1)
	}
	var err error
	dm.page1, err = dm.pageCache.ObtainPage(no)
	if err != nil {
		panic(err)
	}

	// 刷新 page1 数据到磁盘
	dm.pageCache.PageFlush(dm.page1)
}

func NewManage(path string, memory int64, tm txn.Manage) Manage {
	dm := new(dataManage)

	dm.log = log.NewLog(path)
	dm.txnManage = tm

	dm.pageIndex = page.NewIndex()
	dm.pageCache = page.NewCache(path, memory)

	dm.cache = cache.NewCache(&cache.Option{
		Obtain:   dm.obtainForCache,
		Release:  dm.releaseForCache,
		MaxCount: 0,
	})

	return nil
}

// obtainForCache 需要支持并发
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

func (dm *dataManage) Read(id uint64) (Item, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (dm *dataManage) Insert(tid txn.TID, data []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (dm *dataManage) logDataItem(tid txn.TID, item Item) {
	// 包装 update log 数据
	// dm.log.Log()
}

func (dm *dataManage) ReleaseDataItem(item Item) {
	dm.cache.Release(item.Id())
}
