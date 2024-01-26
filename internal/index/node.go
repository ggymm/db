package index

import (
	"db/internal/data"
	"encoding/binary"
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
//  val: 8 bytes（uint64）, 数据的 id（itemId）（如果非叶子节点，则是子节点的 id）

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

func (n *node) getLeaf() bool {
	return n.data[offLeaf] == 1
}

func (n *node) setLeaf(leaf bool) {
	if leaf {
		n.data[offLeaf] = 1
	} else {
		n.data[offLeaf] = 0
	}
}

func (n *node) getKeysNum() int {
	return int(bin.Uint16(n.data[offKeysNum:]))
}

func (n *node) setKeysNum(num int) {
	bin.PutUint16(n.data[offKeysNum:], uint16(num))
}

func (n *node) getSibling() uint64 {
	return bin.Uint64(n.data[offSibling:])
}

func (n *node) setSibling(id uint64) {
	bin.PutUint64(n.data[offSibling:], id)
}

func (n *node) getKey(i int) uint64 {
	return bin.Uint64(n.data[headerLen+i*16:])
}

func (n *node) setKey(i int, key uint64) {
	bin.PutUint64(n.data[headerLen+i*16:], key)
}

func (n *node) getVal(i int) uint64 {
	return bin.Uint64(n.data[headerLen+i*16+8:])
}

func (n *node) setVal(i int, val uint64) {
	bin.PutUint64(n.data[headerLen+i*16+8:], val)
}

func (n *node) readData(i uint64, data []byte) {
	copy(data[headerLen:], n.data[headerLen+i*16:])
}
