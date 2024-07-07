package page

import (
	"container/list"
	"sync"
)

// 页面空余空间缓存
//
// 将一个页面大小划分为 40 份，计算每份的大小
//
// 初始化时，遍历所有页面空余空间
// 使用空余空间大小除以每份大小，得到的商作为下标，将页面添加到对应的链表中
// 例如：空余空间大小为 100，每份大小为 20，商为 5，将页面添加到第 5 个链表中
//
// 选择时，与添加逻辑基本类似，只是注意需要向上取整
// 例如：需要 30 的空间，每份大小为 20，商为 1.5，向上取整为 2，选择第 2 个链表中的页面
// 如果第 2 个链表为空，则选择第 3 个链表中的页面，以此类推
// 最后，选择的页面需要从链表中删除

const (
	interval  = 40
	threshold = Size / interval
)

type Index interface {
	Add(no uint32, free int)
	Select(free int) (uint32, int)
}

type pageIndex struct {
	lock sync.Mutex

	spaceList [interval + 1]list.List
}

type indexItem struct {
	no   uint32
	free int
}

func NewIndex() Index {
	return &pageIndex{
		spaceList: [interval + 1]list.List{},
	}
}

func (pi *pageIndex) Add(no uint32, free int) {
	pi.lock.Lock()
	defer pi.lock.Unlock()

	i := free / threshold
	i = min(i, interval)
	pi.spaceList[i].PushBack(&indexItem{
		no:   no,
		free: free,
	})
}

func (pi *pageIndex) Select(free int) (uint32, int) {
	pi.lock.Lock()
	defer pi.lock.Unlock()

	i := free / threshold
	if i < interval {
		i++
	}
	for i <= interval {
		if pi.spaceList[i].Len() > 0 {
			itemEle := pi.spaceList[i].Front()
			itemVal := pi.spaceList[i].Remove(itemEle)
			item := itemVal.(*indexItem)
			return item.no, item.free
		}
		i++
	}
	return 0, 0
}
