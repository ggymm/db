package index

import (
	"db/internal/tx"
	"sync"

	"db/internal/data"
)

type Index interface {
	Insert(key, itemId uint64) error
	Search(key uint64) ([]uint64, error)
	SearchRange(prev, next uint64) ([]uint64, error)
}

type tree struct {
	lock     sync.Mutex
	rootItem data.Item

	DataManage data.Manage
}

func (t *tree) rootId() uint64 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return bin.Uint64(t.rootItem.DataBody())
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
	t.rootItem.Before()
	buf := make([]byte, 8)
	bin.PutUint64(buf, rootId)
	copy(t.rootItem.DataBody(), buf)
	t.rootItem.After(tx.Super)
	return nil
}

// insert
func (t *tree) insert(nodeId, key, itemId uint64) (uint64, uint64, error) {
	var (
		err error

		nd               *node
		child            uint64
		newKey, newChild uint64
	)

	nd, err = wrapNode(t, nodeId)
	if err != nil {
		return 0, 0, err
	}
	isLeaf := nd.Leaf()

	nd.Release()
	if isLeaf {
		// 叶子节点
		return t.insertNode(nodeId, key, itemId)
	} else {
		// 非叶子节点
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
		err error

		nd               *node
		sibling          uint64
		newKey, newChild uint64
	)
	for {
		nd, err = wrapNode(t, nodeId)
		if err != nil {
			return 0, 0, err
		}
		sibling, newKey, newChild, err = nd.Insert(key, itemId)

		nd.Release()
		if sibling != 0 {
			nodeId = sibling
		} else {
			return newKey, newChild, err
		}
	}
}

func (t *tree) searchNode(nodeId, key uint64) (uint64, error) {
	for {
		nd, err := wrapNode(t, nodeId)
		if err != nil {
			return 0, err
		}
		child, sibling := nd.Search(key)

		// 如果找到符合条件的节点，则返回
		// 如果没有找到，则继续查找下一个节点
		nd.Release()
		if child != 0 {
			return child, nil
		}
		nodeId = sibling
	}
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
	panic("implement me")
}

func (t *tree) SearchRange(prev, next uint64) ([]uint64, error) {
	panic("implement me")
}
