package table

import (
	"github.com/ggymm/db"
	"github.com/ggymm/db/index"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/tx"
	"slices"
)

const (
	Null byte = iota
	NotNull
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

	pos += shift
	t.Fields = make([]*field, 0)

	// 读取 fields
	id := uint64(0)
	for pos < len(data) {
		// 读取 field
		id, shift = decodeUint64(data[pos:])
		t.Fields = append(t.Fields, readField(tbm, id))

		pos += shift
	}
	return t
}

func (t *table) save(txId uint64) (err error) {
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

func (t *table) wrapRaw(row entry) ([]byte, error) {
	raw := make([]byte, 0)
	for _, f := range t.Fields {
		// 获取字段值
		v, exist := row[f.Name]
		if !exist || v == "" {
			if len(f.Default) != 0 {
				v = f.Default
				continue
			}
			if !f.Nullable {
				return nil, NewError(ErrNotAllowNull, f.Name)
			}
		}

		// 获取字段二进制值
		raw = append(raw, f.wrapRaw(v)...)
	}
	return raw, nil
}

func (t *table) wrapEntry(raw []byte, where []sql.SelectWhere) entry {
	pos := 0
	row := make(entry)
	for _, f := range t.Fields {
		val, shift := f.parseRaw(raw[pos:])
		row[f.Name] = val
		pos += shift
	}

	if where == nil || len(where) == 0 {
		return row
	}

	// 过滤条件
	match := true
	for _, w := range where {
		if w.Match(row) && match {
			match = true
		} else {
			match = false
		}
	}
	if !match {
		return nil
	}
	return row
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

func (f *field) save(txId uint64) (err error) {
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

func (f *field) wrapRaw(v any) []byte {
	if v == nil {
		return []byte{Null}
	}

	var raw []byte
	switch f.Type {
	case "INT32":
		raw = bin.Uint32Raw(v.(uint32))
	case "INT64":
		raw = bin.Uint64Raw(v.(uint64))
	case "VARCHAR":
		l := len(v.(string))
		raw = make([]byte, 4+l)

		// length
		raw[0] = byte(l)
		raw[1] = byte(l >> 8)
		raw[2] = byte(l >> 16)
		raw[3] = byte(l >> 24)

		// string
		copy(raw[4:], v.(string))
	}
	return slices.Insert(raw, 0, NotNull)
}

func (f *field) parseRaw(raw []byte) (any, int) {
	if raw[0] == Null {
		return nil, 1
	}
	var v any
	var shift int
	raw = raw[1:]
	switch f.Type {
	case "INT32":
		v = bin.Uint32(raw)
		shift = 4
	case "INT64":
		v = bin.Uint64(raw)
		shift = 8
	case "VARCHAR":
		l := int(raw[0]) |
			int(raw[1])<<8 |
			int(raw[2])<<16 |
			int(raw[3])<<24
		v = string(raw[4 : 4+l])
		shift = l + 4
	}
	return v, shift + 1
}

func (f *field) isIndex() bool {
	return f.TreeId != 0
}
