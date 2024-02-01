package index

import (
	"db/internal/data"
	"encoding/binary"
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
// sibling: 兄弟节点的 id（itemId）
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

var (
	bin = binary.LittleEndian
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

func setLeaf(leaf bool, data []byte) {
	if leaf {
		data[offLeaf] = 1
	} else {
		data[offLeaf] = 0
	}
}

func getKeysNum(data []byte) int {
	return int(bin.Uint16(data[offKeysNum:]))
}

func setKeysNum(num int, data []byte) {
	bin.PutUint16(data[offKeysNum:], uint16(num))
}

func getSibling(data []byte) uint64 {
	return bin.Uint64(data[offSibling:])
}

func setSibling(child uint64, data []byte) {
	bin.PutUint64(data[offSibling:], child)
}

func getOff(i int) int {
	return headerLen + i*8*2
}

func getKey(i int, data []byte) uint64 {
	off := getOff(i)
	return bin.Uint64(data[off:])
}

func setKey(i int, key uint64, data []byte) {
	off := getOff(i)
	bin.PutUint64(data[off:], key)
}

func getChild(i int, data []byte) uint64 {
	off := getOff(i) + 8
	return bin.Uint64(data[off:])
}

func setChild(i int, val uint64, data []byte) {
	off := getOff(i) + 8
	bin.PutUint64(data[off:], val)
}

func shiftData(i int, data []byte) {
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
	setLeaf(true, buf)
	setKeysNum(0, buf)
	setSibling(0, buf)
	return buf
}

func createRoot(left, right uint64, key uint64) []byte {
	buf := make([]byte, nodeSize)
	setLeaf(false, buf) // 非叶子节点
	setKeysNum(2, buf)  // 相当于有两个子节点
	setSibling(0, buf)  // 没有兄弟节点

	// 左节点
	setKey(0, key, buf)
	setChild(0, left, buf)

	// 右节点
	setKey(1, math.MaxUint64, buf)
	setChild(1, right, buf)
	return buf
}

func wrapNode(tree *index, id uint64) (*node, error) {
	item, ok, err := tree.DataManage.Read(id)
	if !ok || err != nil {
		return nil, err
	}

	return &node{
		tree: tree,
		id:   id,
		data: item.DataBody(),
		item: item,
	}, nil
}

func (n *node) Release() {
	n.item.Release()
}

func (n *node) IsLeaf() bool {
	n.item.RLock()
	defer n.item.RUnlock()

	return getLeaf(n.data)
}
