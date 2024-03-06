package table

import (
	"fmt"
	"slices"

	"db/internal/tx"
	"db/pkg/bin"
	"db/pkg/sql"
	"db/pkg/str"
)

// 表
//
// 表内存中数据结构如下：
// +----------------+----------------+----------------+----------------+
// |     name       |    nextPtr     |     field1     |    field...    |
// +----------------+----------------+----------------+----------------+
// |     string     |     uint64     |     uint64     |    uint64      |
// +----------------+----------------+----------------+----------------+
//

type entry struct {
	raw   []byte
	value []any
	field []*field
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

func (t *table) raw(stmt *sql.InsertStmt) ([]entry, error) {
	es := make([]entry, len(t.fields))

	maps := make([]map[string]string, len(stmt.Value))
	for i, val := range stmt.Value {
		maps[i] = make(map[string]string)
		for j, f := range stmt.Field {
			maps[i][f] = val[j]
		}
	}

	for _, insert := range maps {
		e := entry{
			raw:   make([]byte, 0),
			value: make([]any, len(t.fields)),
			field: make([]*field, len(t.fields)),
		}
		e.raw = make([]byte, 0)
		for i, f := range t.fields {
			e.field[i] = f

			// 获取字段值
			val, ok := insert[f.name]
			switch {
			case ok:
				e.value[i] = val
			case len(f.defaultVal) != 0:
				e.value[i] = f.defaultVal
			case f.allowNull:
				e.value[i] = nil
			default:
				return nil, fmt.Errorf("field %s is not allowed to be null", f.name)
			}

			// 获取字段二进制值
			e.raw = append(e.raw, sql.FieldRaw(f.dataType, e.value[i])...)
		}
		es = append(es, e)
	}
	return es, nil
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
