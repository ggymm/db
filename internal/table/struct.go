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

type table struct {
	tbm Manage

	id     uint64
	name   string
	fields []*field

	nextId uint64
}

type newTable struct {
	TxId   uint64
	NextId uint64
	Stmt   *sql.CreateStmt
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
		tbm:    tbm,
		id:     id,
		name:   name,
		fields: fields,
		nextId: next,
	}
}

func createTable(tbm *tableManage, info *newTable) (*table, error) {
	t := &table{
		tbm:    tbm,
		name:   info.Stmt.Name,
		nextId: info.NextId,
	}

	// 读取 index
	index := make([]string, 0)
	for _, i := range info.Stmt.Table.Index {
		index = append(index, i.Field)
	}

	// 读取 field
	for _, f := range info.Stmt.Table.Field {
		fld, err := createField(tbm, &newField{
			TxId:    info.TxId,
			Name:    f.Name,
			Type:    f.Type.String(),
			Indexed: slices.Contains(index, f.Name),
		})
		if err != nil {
			return nil, err
		}
		t.fields = append(t.fields, fld)
	}

	// 持久化
	return t, t.persist(info.TxId)
}

func (t *table) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(t.name)

	// next
	buf := make([]byte, 8)
	bin.PutUint64(buf, t.nextId)
	data = append(data, buf...)

	// fields
	for _, f := range t.fields {
		bin.PutUint64(buf, f.id)
		data = append(data, buf...)
	}

	// 持久化
	t.id, err = t.tbm.VerManage().Insert(txId, data)
	return
}

type field struct {
	tbm Manage
	idx index.Index

	id         uint64
	name       string
	index      uint64
	dataType   string
	allowNull  bool
	defaultVal string
}

type newField struct {
	TxId    uint64
	Name    string
	Type    string
	Indexed bool
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
	f.name, shift = str.Deserialize(data)

	// 读取 type
	pos += shift
	f.dataType, shift = str.Deserialize(data[pos:])

	// 读取 index
	pos += shift
	f.index = bin.Uint64(data[pos:])
	if f.index != 0 {
		f.idx, err = index.NewIndex(tbm.DataManage(), &opt.Option{
			Open:   true,
			RootId: f.index,
		})
		if err != nil {
			panic(err)
		}
	}
	return f
}

func createField(tbm Manage, info *newField) (*field, error) {
	f := &field{
		tbm:      tbm,
		name:     info.Name,
		index:    0,
		dataType: info.Type,
	}

	if info.Indexed {
		idx, err := index.NewIndex(tbm.DataManage(), &opt.Option{
			Open: false,
		})
		if err != nil {
			return nil, err
		}
		f.index = idx.GetRootId()
	}
	return f, f.persist(info.TxId)
}

func (f *field) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(f.name)

	// type
	data = append(data, str.Serialize(f.dataType)...)

	// index
	buf := make([]byte, 8)
	bin.PutUint64(buf, f.index)
	data = append(data, buf...)

	// 持久化
	f.id, err = f.tbm.VerManage().Insert(txId, data)
	return
}

func (f *field) isIndexed() bool {
	return f.index != 0
}
