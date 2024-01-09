package data

import (
	"db/internal/cache"
	"db/internal/data/log"
	"db/internal/data/page"
	"db/internal/txn"
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

func open(d *dataManage) {

}

func create(d *dataManage) {

}

func NewCache(memory int64, filename string, tm txn.Manage) Manage {
	m := new(dataManage)

	m.log = log.NewLog(filename)
	m.txnManage = tm

	m.pageIndex = page.NewIndex()
	m.pageCache = page.NewCache(memory, filename)

	m.cache = cache.NewCache(&cache.Option{
		Obtain:   m.obtainForCache,
		Release:  m.releaseForCache,
		MaxCount: 0,
	})

	return nil
}

// obtainForCache 需要支持并发
func (d *dataManage) obtainForCache(key uint64) (any, error) {
	no, off := parseDataItemId(key)
	p, err := d.pageCache.ObtainPage(no)
	if err != nil {
		return nil, err
	}
	return parseDataItem(p, off, d), nil
}

// releaseForCache 需要是同步方法
func (d *dataManage) releaseForCache(data any) {
	item := data.(Item)
	item.Page().Release()
}

func (d *dataManage) Close() {
	d.log.Close()
	d.cache.Close()

	SetVcClose(d.page1) // 设置 page1 的校验数据
	d.page1.Release()
	d.pageCache.Close()
}

func (d *dataManage) Read(id uint64) (Item, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dataManage) Insert(tid txn.TID, data []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
