package table

import (
	"github.com/ggymm/db"
	"github.com/ggymm/db/index"
	"github.com/ggymm/db/tx"
)

type entry map[string]any

// table 内存结构
//
// +----------------+----------------+----------------+
// |     name       |     nextId     |     fields     |
// +----------------+----------------+----------------+
// |    string      |     uint64     |    uint64[]    |
// +----------------+----------------+----------------+
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

	// 读取 name
	name, pos := decodeString(data)

	// 读取 next
	next := decodeUint64(data[pos:])

	// 读取 fields
	pos += 8
	fields := make([]*field, 0)
	for pos < len(data) {
		f := decodeUint64(data[pos:])
		fields = append(fields, readField(tbm, f))
		pos += 8
	}

	return &table{
		tbm:    tbm,
		itemId: itemId,

		Name:   name,
		Next:   next,
		Fields: fields,
	}
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

// field 内存结构
//
// Name: 名称
// Type: 类型
// IndexId: 索引根节点 itemId
// Default: 默认值
// allowNull: 是否允许为空
// primaryKey: 是否是主键
type field struct {
	tbm    Manage
	index  index.Index
	itemId uint64

	Name       string // 名称
	Type       string // 类型
	IndexId    uint64 // 索引
	AllowNull  bool   // 是否允许为空
	DefaultVal string // 默认值
	PrimaryKey bool   // 是否是主键
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

	// default
	pos += shift
	f.DefaultVal, shift = decodeString(data[pos:])

	// allowNull
	pos += shift
	f.AllowNull = data[pos] == 1

	// primaryKey
	pos++
	f.PrimaryKey = data[pos] == 1

	// index
	pos++
	f.IndexId = decodeUint64(data[pos:])

	// 读取索引
	if f.IndexId != 0 {
		f.index, err = index.NewIndex(tbm.DataManage(), &db.Option{
			Open:   true,
			RootId: f.IndexId,
		})
		if err != nil {
			panic(err)
		}
	}
	return f
}

// persist 将该field持久化
//
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |	  name      |	  type       |	   default    |	   allowNull   |	primaryKey  |	  index      |
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |	 string     |	 string      |	   string     |	     bool      |	  bool      |	  uint64     |
// +----------------+----------------+----------------+----------------+----------------+----------------+
func (f *field) persist(txId uint64) (err error) {
	// name
	data := encodeString(f.Name)

	// type
	data = append(data, encodeString(f.Type)...)

	// default
	data = append(data, encodeString(f.DefaultVal)...)

	// allowNull
	if f.AllowNull {
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

	// index
	data = append(data, encodeUint64(f.IndexId)...)

	// 保存到磁盘
	f.itemId, err = f.tbm.VerManage().Write(txId, data)
	return
}

func (f *field) isIndex() bool {
	return f.IndexId != 0
}
