package index

import (
	"math"

	"db/data"
	"db/pkg/bin"
	"db/tx"
)

// node
// B+Tree 的节点结构
//
// +-------------------+-------------------+-------------------+-------------------+
// |      isLeaf       |      keysNum      |      sibling      |   key and child   |
// +-------------------+-------------------+-------------------+-------------------+
// |      1 byte       |      2 bytes      |      8 bytes      | 16 bytes * 32 * 2 |
// +-------------------+-------------------+-------------------+-------------------+
//
// isLeaf: 是否为叶子节点
// keysNum: 节点中的 key 的数量
// sibling: 兄弟节点的 id（itemId）（如果非叶子节点，则是子节点的 id）
// keys
//  key: 8 bytes（uint64）, 字段的 hash 值
//  child: 8 bytes（uint64）, 数据的 id（itemId）（如果非叶子节点，则是子节点的 id）
//
// 特殊处理
// 其他的 B+Tree 算法中，非叶子节点的会有一个指针指向最右边的子节点
// 这里取消了指针，同时将 keyN 设置为 math.MaxUint64，childN 指向最右侧的子节点
//
// 这样，非叶子节点和叶子节点的二进制结构保持一致
// 执行查询操作时，也保持了和叶子节点一致的查询逻辑

const (
	offLeaf    = 0
	offKeysNum = offLeaf + 1
	offSibling = offKeysNum + 2

	headerLen  = 1 + 2 + 8
	balanceNum = 32

	nodeSize = headerLen + (2*8)*(balanceNum*2+2)
)

type node struct {
	tree *tree

	id   uint64
	data []byte
	item data.Item
}

func getLeaf(data []byte) bool {
	return data[offLeaf] == byte(1)
}

func setLeaf(data []byte, leaf bool) {
	if leaf {
		data[offLeaf] = byte(1)
	} else {
		data[offLeaf] = byte(0)
	}
}

func getKeysNum(data []byte) int {
	return int(bin.Uint16(data[offKeysNum:]))
}

func setKeysNum(data []byte, num int) {
	bin.PutUint16(data[offKeysNum:], uint16(num))
}

func getSibling(data []byte) uint64 {
	return bin.Uint64(data[offSibling:])
}

func setSibling(data []byte, sibling uint64) {
	bin.PutUint64(data[offSibling:], sibling)
}

func getOff(i int) int {
	return headerLen + i*8*2
}

func getKey(data []byte, i int) uint64 {
	off := getOff(i)
	return bin.Uint64(data[off:])
}

func setKey(data []byte, i int, key uint64) {
	off := getOff(i)
	bin.PutUint64(data[off:], key)
}

func getChild(data []byte, i int) uint64 {
	off := getOff(i) + 8
	return bin.Uint64(data[off:])
}

func setChild(data []byte, i int, val uint64) {
	off := getOff(i) + 8
	bin.PutUint64(data[off:], val)
}

func shiftData(data []byte, i int) {
	start := getOff(i + 1)
	length := nodeSize - 1
	for n := length; n >= start; n-- {
		data[n] = data[n-8*2]
	}
}

func writeInitData(i int, dst, src []byte) {
	off := getOff(i)
	copy(dst[headerLen:], src[off:])
}

func initRoot() []byte {
	buf := make([]byte, nodeSize)
	setLeaf(buf, true)
	setKeysNum(buf, 0)
	setSibling(buf, 0)
	return buf
}

func createRoot(key, prev, next uint64) []byte {
	buf := make([]byte, nodeSize)
	setLeaf(buf, false) // 非叶子节点
	setKeysNum(buf, 2)  // 相当于有两个子节点
	setSibling(buf, 0)  // 没有兄弟节点

	// 左节点
	setKey(buf, 0, key)
	setChild(buf, 0, prev)

	// 右节点
	setKey(buf, 1, math.MaxUint64)
	setChild(buf, 1, next)
	return buf
}

func wrapNode(t *tree, id uint64) (*node, error) {
	item, ok, err := t.DataManage.Read(id)
	if !ok || err != nil {
		return nil, err
	}

	return &node{
		tree: t,
		id:   id,
		data: item.DataBody(),
		item: item,
	}, nil
}

func release(n *node) {
	n.item.Release()
}

// split 将节点分为 prev 和 next 两个节点
// 并且返回 next 节点的第一个 key 和新节点的 id
//
// newKey: next 节点的第一个 key
// newChild: next 节点的 itemId
/*
   key0, child0, key60, child60, INF, child99

                newKey                 newChild
                  ↓                       ↓
   key0, child0, key60, child60, key60, child11, INF, child99
*/
func (n *node) split() (uint64, uint64, error) {
	buf := make([]byte, nodeSize)

	// 设置新节点的属性
	setLeaf(buf, getLeaf(n.data))
	setKeysNum(buf, balanceNum)
	setSibling(buf, getSibling(n.data))
	writeInitData(balanceNum, buf, n.data)

	// 插入新节点
	newChild, err := n.tree.DataManage.Insert(tx.Super, buf)
	if err != nil {
		return 0, 0, err
	}

	// 修改原节点属性
	setKeysNum(n.data, balanceNum)
	setSibling(n.data, newChild)
	return getKey(buf, 0), newChild, nil
}

func (n *node) insert(key, itemId uint64) bool {
	// 遍历节点的 key，找到合适的位置
	num := getKeysNum(n.data)
	var i int
	for i < num {
		if key <= getKey(n.data, i) {
			break
		}
		i++
	}
	if i == num && getSibling(n.data) != 0 {
		// 如果是最后一个节点，且有兄弟节点，则需要向兄弟节点插入
		return false
	}

	if getLeaf(n.data) {
		// 叶子节点
		// 此时，直接将 key 和 itemId 插入到节点对应位置
		shiftData(n.data, i)
		setKey(n.data, i, key)
		setChild(n.data, i, itemId)
	} else {
		// 非叶子节点（子节点进行了分裂操作）
		// 需要将 key 插入到 i 的位置
		// 需要将 itemId 插入到 i+1 的位置
		shiftData(n.data, i)
		setKey(n.data, i, key)
		setChild(n.data, i+1, itemId)

		// nextKey := getKey(n.data, i)
		// setKey(n.data, i, key)
		// shiftData(n.data, i+1)
		// setKey(n.data, i+1, nextKey)
		// setChild(n.data, i+1, itemId)
	}
	setKeysNum(n.data, num+1)
	return true
}

func (n *node) IsLeaf() bool {
	n.item.RLock()
	defer n.item.RUnlock()

	return getLeaf(n.data)
}

// Insert 插入数据
//
// 返回值：
// sibling: 当前节点插入失败时，返回兄弟节点的 id
// newKey: 当前节点分裂时，返回新节点的第一个 key
// newChild: 当前节点分裂时，返回新节点的 itemId
func (n *node) Insert(key, itemId uint64) (uint64, uint64, uint64, error) {
	var (
		err     error
		success bool

		newKey   uint64
		newChild uint64
	)

	n.item.Before()
	defer func() {
		if err == nil && success {
			n.item.After(tx.Super)
		} else {
			n.item.UnBefore()
		}
	}()

	success = n.insert(key, itemId)
	if !success {
		return getSibling(n.data), 0, 0, nil
	}

	// 检查是否需要分裂
	if getKeysNum(n.data) == balanceNum*2 {
		newKey, newChild, err = n.split()
	}
	return 0, newKey, newChild, err
}

// Search 查找数据
//
// 返回值：
// childId: 子节点的 id
// siblingId: 兄弟节点的 id
func (n *node) Search(key uint64) (uint64, uint64) {
	n.item.RLock()
	defer n.item.RUnlock()

	nums := getKeysNum(n.data)
	for i := 0; i < nums; i++ {
		// 判断 key 是否小于当前节点的 key
		if key < getKey(n.data, i) {
			return getChild(n.data, i), 0
		}
	}
	return 0, getSibling(n.data)
}

// SearchRange 查找范围内的数据
//
// 返回值：
// []uint64: 满足条件的子节点 id
// uint64: 兄弟节点的 id
func (n *node) SearchRange(prevKey, nextKey uint64) ([]uint64, uint64) {
	n.item.RLock()
	defer n.item.RUnlock()

	res := make([]uint64, 0)
	num := getKeysNum(n.data)
	var i int
	// 查找大于等于 prevKey 的 key index
	for i < num {
		if prevKey <= getKey(n.data, i) {
			break
		}
		i++
	}
	// 查找小于等于 nextKey 的 key index
	// 将所有满足条件的 childId 加入到 res 中
	for i < num {
		if nextKey < getKey(n.data, i) {
			break
		}
		res = append(res, getChild(n.data, i))
		i++
	}

	var sibling uint64
	if i == num {
		sibling = getSibling(n.data)
	}
	return res, sibling
}
