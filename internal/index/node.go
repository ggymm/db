package index

import (
	"db/internal/data"
	"math"
)

// node
// B+Tree 的节点结构
//
// +-------------------+-------------------+-------------------+-------------------+
// |      isLeaf       |      keysNum      |      sibling      |    key and son    |
// +-------------------+-------------------+-------------------+-------------------+
// |      1 byte       |      2 bytes      |      8 bytes      | 16 bytes * 32 * 2 |
// +-------------------+-------------------+-------------------+-------------------+
//
// isLeaf: 是否为叶子节点
// keysNum: 节点中的 key 的数量
// sibling: 兄弟节点的 id（itemId）（如果非叶子节点，则是第一个子节点的 id）
// keys
//  key: 8 bytes（uint64）, 字段的 hash 值
//  child: 8 bytes（uint64）, 数据的 id（itemId）（如果非叶子节点，则是子节点的 id）

const (
	offLeaf    = 0
	offKeysNum = offLeaf + 1
	offSibling = offKeysNum + 2

	headerLen  = 1 + 2 + 8
	balanceNum = 32

	nodeSize = headerLen + (2*8)*(balanceNum*2+2)
)

type node struct {
	tree Index

	id   uint64
	data []byte
	item data.Item
}

func getLeaf(data []byte) bool {
	return data[offLeaf] == 1
}

func setLeaf(data []byte, leaf bool) {
	if leaf {
		data[offLeaf] = 1
	} else {
		data[offLeaf] = 0
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

func setSibling(data []byte, child uint64) {
	bin.PutUint64(data[offSibling:], child)
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
	for k := length; k >= start; k-- {
		data[k] = data[k-8*2]
	}
}

func writeData(i int, src, dst []byte) {
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

func createRoot(left, right uint64, key uint64) []byte {
	buf := make([]byte, nodeSize)
	setLeaf(buf, false) // 非叶子节点
	setKeysNum(buf, 2)  // 相当于有两个子节点
	setSibling(buf, 0)  // 没有兄弟节点

	// 左节点
	setKey(buf, 0, key)
	setChild(buf, 0, left)

	// 右节点
	setKey(buf, 1, math.MaxUint64)
	setChild(buf, 1, right)
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

func (n *node) split() {

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
		// 直接将 key 和 itemId 插入到节点中
		shiftData(n.data, i)
		setKey(n.data, i, key)
		setChild(n.data, i, itemId)
	} else {
		// 非叶子节点
		// 获取当前 i 位置的 nextKey（大于 key 的）
		nextKey := getKey(n.data, i)

		// 将 key 插入到 i 的位置
		// 不需要插入 itemId，是因为此时 sibling 是子节点的 id
		setKey(n.data, i, key)

		// 将 nextKey 插入到 i+1 的位置
		// 因为 nextKey 是大于 key 的，查找时可以根据 nextKey 找到对应的子节点
		shiftData(n.data, i+1)
		setKey(n.data, i+1, nextKey)
		setChild(n.data, i+1, itemId)
	}
	setKeysNum(n.data, num+1)
	return true
}

func (n *node) Release() {
	n.item.Release() // 从缓存中释放引用
}

func (n *node) Leaf() bool {
	n.item.RLock()
	defer n.item.RUnlock()

	return getLeaf(n.data)
}

func (n *node) Insert(key, itemId uint64) (uint64, uint64, uint64, error) {

	return 0, 0, 0, nil
}

// Search 查找数据
// 如果无法找到，则返回下一个节点
// childId: 子节点的 id
// siblingId: 兄弟节点的 id
func (n *node) Search(key uint64) (childId, siblingId uint64) {
	n.item.RLock()
	defer n.item.RUnlock()

	nums := getKeysNum(n.data)
	for i := 0; i < nums; i++ {
		if key < getKey(n.data, i) {
			return getChild(n.data, i), 0
		}
	}
	return 0, getSibling(n.data)
}
