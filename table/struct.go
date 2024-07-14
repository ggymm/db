package table

import (
	"github.com/ggymm/db"
	"github.com/ggymm/db/index"
	"github.com/ggymm/db/tx"
)

type entry map[string]any

// table 结构
//
// +----------------+----------------+----------------+
// |     Name       |      Next      |     Fields     |
// +----------------+----------------+----------------+
// |    string      |     uint64     |    uint64[]    |
// +----------------+----------------+----------------+
// Name: 表名
// Next: 下一张表的 itemId
// Fields: 表字段 itemId 列表
type table struct {
	tbm    Manage
	itemId uint64

	Name   string
	Next   uint64
	Fields []*field
}

func readTable(tbm Manage, itemId uint64) *table {
	data, exist, err := tbm.VerManage().Read(tx.Super, itemId)
	if err != nil || !exist {
		panic(err)
	}

	t := &table{
		tbm:    tbm,
		itemId: itemId,
	}
	var (
		pos   int
		shift int
	)

	// 读取 name
	t.Name, shift = decodeString(data)

	// 读取 next
	pos += shift
	t.Next, shift = decodeUint64(data[pos:])

	t.Fields = make([]*field, 0)
	// 读取 fields
	id := uint64(0)
	for pos < len(data) {
		pos += shift

		// 读取 field
		id, shift = decodeUint64(data[pos:])
		t.Fields = append(t.Fields, readField(tbm, id))
	}
	return t
}

func (t *table) persist(txId uint64) (err error) {
	// name
	data := encodeString(t.Name)

	// next
	raw := encodeUint64(t.Next)
	data = append(data, raw...)

	// fields
	for _, f := range t.Fields {
		raw = encodeUint64(f.itemId)
		data = append(data, raw...)
	}

	// 持久化
	t.itemId, err = t.tbm.VerManage().Write(txId, data)
	return
}

// field 字段信息
//
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |	  name      |	  type       |	   default    |	   allowNull   |	primaryKey  |	  index      |
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |	 string     |	 string      |	   string     |	     bool      |	  bool      |	  uint64     |
// +----------------+----------------+----------------+----------------+----------------+----------------+
//
// Name: 名称
// Type: 类型
// TreeId: 索引根节点 itemId
// Default: 默认值
// Nullable: 是否允许为空
// PrimaryKey: 是否是主键
type field struct {
	tbm    Manage
	index  index.Index
	itemId uint64

	Name       string
	Type       string
	TreeId     uint64
	Default    string
	Nullable   bool
	PrimaryKey bool
}

func readField(tbm Manage, itemId uint64) *field {
	data, exist, err := tbm.VerManage().Read(tx.Super, itemId)
	if err != nil || !exist {
		panic(err)
	}

	f := &field{
		tbm:    tbm,
		itemId: itemId,
	}
	var (
		pos   int
		shift int
	)

	// name
	f.Name, shift = decodeString(data)

	// type
	pos += shift
	f.Type, shift = decodeString(data[pos:])

	// treeId
	pos += shift
	f.TreeId, shift = decodeUint64(data[pos:])

	// default
	pos += shift
	f.Default, shift = decodeString(data[pos:])

	// nullable
	pos += shift
	f.Nullable = data[pos] == 1

	// primaryKey
	pos++
	f.PrimaryKey = data[pos] == 1

	// 读取索引
	if f.TreeId != 0 {
		f.index, err = index.NewIndex(tbm.DataManage(), &db.Option{
			Open:   true,
			RootId: f.TreeId,
		})
		if err != nil {
			panic(err)
		}
	}
	return f
}

// persist 将该field持久化
func (f *field) persist(txId uint64) (err error) {
	// name
	data := encodeString(f.Name)

	// type
	data = append(data, encodeString(f.Type)...)

	// treeId
	data = append(data, encodeUint64(f.TreeId)...)

	// default
	data = append(data, encodeString(f.Default)...)

	// nullable
	if f.Nullable {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	// primaryKey
	if f.PrimaryKey {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	// 保存到磁盘
	f.itemId, err = f.tbm.VerManage().Write(txId, data)
	return
}

func (f *field) isIndex() bool {
	return f.TreeId != 0
}
