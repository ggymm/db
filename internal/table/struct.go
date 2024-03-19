package table

import (
	"db/internal/index"
	"db/internal/opt"
	"db/internal/tx"
	"db/pkg/bin"
	"db/pkg/sql"
	"db/pkg/str"
	"slices"
)

type entry struct {
	raw    []byte
	value  []any
	fields []*field
}

// table 内存结构
//
// +----------------+----------------+----------------+
// |     name       |     nextId     |     fields     |
// +----------------+----------------+----------------+
// |    string      |     uint64     |    uint64[]    |
// +----------------+----------------+----------------+
type table struct {
	tbm Manage

	itemId uint64

	tableName   string
	tableNext   uint64
	tableFields []*field
}

type newTable struct {
	txId uint64

	tableName string
	tableNext uint64

	index []*sql.CreateIndex
	field []*sql.CreateField
}

func readTable(tbm Manage, id uint64) *table {
	data, exist, err := tbm.VerManage().Read(tx.Super, id)
	if err != nil || !exist {
		panic(err)
	}

	// 读取 name
	name, pos := str.Deserialize(data)

	// 读取 next
	next := bin.Uint64(data[pos:])

	// 读取 fields
	pos += 8
	fields := make([]*field, 0)
	for pos < len(data) {
		f := bin.Uint64(data[pos:])
		fields = append(fields, readField(tbm, f))
		pos += 8
	}

	return &table{
		tbm:         tbm,
		itemId:      id,
		tableName:   name,
		tableNext:   next,
		tableFields: fields,
	}
}

func createTable(tbm *tableManage, info *newTable) (*table, error) {
	t := &table{
		tbm:       tbm,
		tableName: info.tableName,
		tableNext: info.tableNext,
	}

	// 读取 主键 和 索引
	pkList := make([]string, 0)
	indexList := make([]string, 0)
	for _, i := range info.index {
		if i.Pk {
			pkList = append(pkList, i.Field)
		}
		indexList = append(indexList, i.Field)
	}

	// 读取 field
	for _, f := range info.field {
		hasPk := slices.Contains(pkList, f.Name)
		hasIndex := slices.Contains(indexList, f.Name)

		// 如果是主键, 则不允许为空
		if hasPk {
			f.AllowNull = false
		}
		fd, err := createField(tbm, &newField{
			txId: info.txId,

			fieldName:  f.Name,
			fieldType:  f.Type.String(),
			fieldIndex: hasIndex,

			allowNull:  f.AllowNull,
			primaryKey: hasPk,
			defaultVal: f.DefaultVal,
		})
		if err != nil {
			return nil, err
		}
		t.tableFields = append(t.tableFields, fd)
	}

	// 持久化
	return t, t.persist(info.txId)
}

func (t *table) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(t.tableName)

	// next
	raw := bin.Uint64Raw(t.tableNext)
	data = append(data, raw...)

	// fields
	for _, f := range t.tableFields {
		raw = bin.Uint64Raw(f.itemId)
		data = append(data, raw...)
	}

	// 持久化
	t.itemId, err = t.tbm.VerManage().Insert(txId, data)
	return
}

// field 内存结构
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |    fieldName   |    fieldType   |   fieldIndex   |   allowNull    |   primaryKey   |   defaultVal   |
// +----------------+----------------+----------------+----------------+----------------+----------------+
// |     string     |     string     |     uint64     |     bool       |      bool      |     string     |
// +----------------+----------------+----------------+----------------+----------------+----------------+
//
// fieldName: 字段名
// fieldType: 字段类型
// fieldIndex: 字段索引Id
// allowNull: 是否允许为空
// primaryKey: 是否是主键
// defaultVal: 为空时的默认值
type field struct {
	tbm Manage
	idx index.Index

	itemId uint64

	fieldName  string // 字段名
	fieldType  string // 字段类型
	fieldIndex uint64 // 字段索引Id

	allowNull  bool   // 是否允许为空
	primaryKey bool   // 是否是主键
	defaultVal string // 为空时的默认值
}

type newField struct {
	txId uint64

	fieldName  string
	fieldType  string
	fieldIndex bool

	allowNull  bool
	primaryKey bool
	defaultVal string
}

func readField(tbm Manage, id uint64) *field {
	data, exist, err := tbm.VerManage().Read(tx.Super, id)
	if err != nil || !exist {
		panic(err)
	}
	f := &field{}
	var (
		pos   int
		shift int
	)

	// 读取 name
	f.fieldName, shift = str.Deserialize(data)

	// 读取 type
	pos += shift
	f.fieldType, shift = str.Deserialize(data[pos:])

	// 读取 index
	pos += shift
	f.fieldIndex = bin.Uint64(data[pos:])
	if f.fieldIndex != 0 {
		f.idx, err = index.NewIndex(tbm.DataManage(), &opt.Option{
			Open:   true,
			RootId: f.fieldIndex,
		})
		if err != nil {
			panic(err)
		}
	}
	return f
}

func createField(tbm Manage, info *newField) (*field, error) {
	f := &field{
		tbm: tbm,

		fieldName:  info.fieldName,
		fieldType:  info.fieldType,
		fieldIndex: 0,

		allowNull:  info.allowNull,
		primaryKey: info.primaryKey,
		defaultVal: info.defaultVal,
	}

	if info.fieldIndex {
		idx, err := index.NewIndex(tbm.DataManage(), &opt.Option{
			Open: false,
		})
		if err != nil {
			return nil, err
		}
		f.fieldIndex = idx.GetRootId()
	}
	return f, f.persist(info.txId)
}

func (f *field) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(f.fieldName)

	// type
	raw := str.Serialize(f.fieldType)
	data = append(data, raw...)

	// index
	raw = bin.Uint64Raw(f.fieldIndex)
	data = append(data, raw...)

	// 持久化
	f.itemId, err = f.tbm.VerManage().Insert(txId, data)
	return
}

func (f *field) isIndex() bool {
	return f.fieldIndex != 0
}

func (f *field) isPrimaryKey() bool {
	return f.primaryKey
}
