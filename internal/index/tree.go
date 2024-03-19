package index

import (
	"sync"

	"db/internal/data"
	"db/internal/opt"
	"db/internal/tx"
	"db/pkg/bin"
)

type Index interface {
	Close()

	Insert(key, itemId uint64) error
	Search(key uint64) ([]uint64, error)
	SearchRange(prev, next uint64) ([]uint64, error)

	GetRootId() uint64
}

type tree struct {
	lock sync.Mutex
	root data.Item

	DataManage data.Manage
}

func NewIndex(dm data.Manage, opt *opt.Option) (Index, error) {
	var (
		ok  bool
		err error

		itemId uint64
		rootId uint64

		root data.Item
	)

	if opt.Open {
		rootId = opt.RootId
	} else {
		item := initRoot()
		itemId, err = dm.Insert(tx.Super, item)
		if err != nil {
			return nil, err
		}

		raw := bin.Uint64Raw(itemId)
		rootId, err = dm.Insert(tx.Super, raw)
		if err != nil {
			return nil, err
		}
	}

	// 读取根节点
	root, ok, err = dm.Read(rootId)
	if !ok || err != nil {
		return nil, err
	}
	return &tree{root: root, DataManage: dm}, nil
}

func (t *tree) rootId() uint64 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return bin.Uint64(t.root.DataBody())
}

func (t *tree) updateRootId(key, prev, next uint64) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 插入根节点
	root := createRoot(key, prev, next)
	rootId, err := t.DataManage.Insert(tx.Super, root)
	if err != nil {
		return err
	}

	// 更新根节点Id
	t.root.Before()
	raw := bin.Uint64Raw(rootId)
	copy(t.root.DataBody(), raw)
	t.root.After(tx.Super)
	return nil
}

// insert
func (t *tree) insert(nodeId, key, itemId uint64) (uint64, uint64, error) {
	var (
		nd  *node
		err error
	)

	nd, err = wrapNode(t, nodeId)
	if err != nil {
		return 0, 0, err
	}
	isLeaf := nd.IsLeaf()

	// 释放 node 引用
	release(nd)

	// 判断是否是叶子节点
	if isLeaf {
		return t.insertNode(nodeId, key, itemId)
	} else {
		var (
			child            uint64
			newKey, newChild uint64
		)
		// 查找可以插入的子节点，一直查找到叶子节点
		child, err = t.searchNode(nodeId, key)
		if err != nil {
			return 0, 0, err
		}
		newKey, newChild, err = t.insert(child, key, itemId)
		if err != nil {
			return 0, 0, err
		}

		// 如果新的子节点不为 0 则代表下一层产生了分裂
		// 需要在当前层插入分裂的节点信息 newKey 和 newChild
		if newChild != 0 {
			// 此处可以判断 newKey 或者 newChild 是否为 0
			// 如果不为 0 则代表插入数据后产生了分裂，需要继续向上层插入
			return t.insertNode(nodeId, newKey, newChild)
		}
	}
	return 0, 0, err
}

// insertNode
// 向 node 中插入 key 和 itemId
// 如果需要分裂，则返回新的 key 和新的 child
func (t *tree) insertNode(nodeId, key, itemId uint64) (uint64, uint64, error) {
	var (
		nd  *node
		err error

		sibling          uint64
		newKey, newChild uint64
	)
	for {
		nd, err = wrapNode(t, nodeId)
		if err != nil {
			return 0, 0, err
		}
		sibling, newKey, newChild, err = nd.Insert(key, itemId)

		// 释放 node 引用
		release(nd)

		// 判断是否需要继续查找下一个节点
		if sibling != 0 {
			nodeId = sibling
		} else {
			return newKey, newChild, err
		}
	}
}

// search
// 从 node 的子节点中查找 key 直到找到对应的叶子节点 id（itemId）
func (t *tree) search(nodeId, key uint64) (uint64, error) {
	var (
		nd  *node
		err error
	)

	nd, err = wrapNode(t, nodeId)
	if err != nil {
		return 0, err
	}
	isLeaf := nd.IsLeaf()

	// 释放 node 引用
	release(nd)

	// 判断是否是叶子节点
	if isLeaf {
		return nodeId, nil
	} else {
		var next uint64
		next, err = t.searchNode(nodeId, key)
		if err != nil {
			return 0, err
		}
		return t.search(next, key)
	}
}

func (t *tree) searchNode(nodeId, key uint64) (uint64, error) {
	for {
		nd, err := wrapNode(t, nodeId)
		if err != nil {
			return 0, err
		}
		child, sibling := nd.Search(key)

		// 释放 node 引用
		release(nd)

		// 如果找到符合条件的节点，则返回
		// 如果没有找到，则继续查找下一个节点
		if child != 0 {
			return child, nil
		}
		nodeId = sibling
	}
}

func (t *tree) Close() {
	t.root.Release()
}

// Insert
// 插入 key（字段计算的hash值） 和 itemId（数据项的Id） 的索引关系
func (t *tree) Insert(key, itemId uint64) error {
	rootId := t.rootId()

	newKey, newChild, err := t.insert(rootId, key, itemId)
	if err != nil {
		return err
	}

	if newChild != 0 {
		// 需要变更根节点
		err = t.updateRootId(newKey, rootId, newChild)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tree) Search(key uint64) ([]uint64, error) {
	return t.SearchRange(key, key)
}

func (t *tree) SearchRange(prevKey, nextKey uint64) ([]uint64, error) {
	var (
		err error

		nd     *node
		prevId uint64
	)

	prevId, err = t.search(t.rootId(), prevKey)
	if err != nil {
		return nil, err
	}

	var res []uint64
	for {
		nd, err = wrapNode(t, prevId)
		if err != nil {
			return nil, err
		}
		tmp, sibling := nd.SearchRange(prevKey, nextKey)
		res = append(res, tmp...)

		// 释放 node 引用
		release(nd)

		// 判断是否需要继续查找下一个节点
		if sibling == 0 {
			break
		}
		prevId = sibling
	}
	return res, nil
}

func (t *tree) GetRootId() uint64 {
	return t.rootId()
}
