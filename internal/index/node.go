package index

// node
// B+Tree 的节点结构
const (
	offFlag    = 0
	offKeyNum  = offFlag + 1
	offSibling = offKeyNum + 2

	headerLen = offSibling + 8
)
