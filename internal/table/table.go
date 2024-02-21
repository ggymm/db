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

type table struct {
	tbm Manage

	Id     uint64
	Name   string
	Fields []*field

	NextId uint64
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
		tbm: tbm,

		Id:     id,
		Name:   name,
		Fields: fields,

		NextId: next,
	}
}

func createTable(tbm *tableManage, info *newTable) (*table, error) {
	t := &table{
		tbm:    tbm,
		Name:   info.Stmt.Name,
		NextId: info.NextId,
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
		t.Fields = append(t.Fields, fld)
	}

	// 持久化
	return t, t.persist(info.TxId)
}

func (t *table) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(t.Name)

	// next
	buf := make([]byte, 8)
	bin.PutUint64(buf, t.NextId)
	data = append(data, buf...)

	// fields
	for _, f := range t.Fields {
		bin.PutUint64(buf, f.Id)
		data = append(data, buf...)
	}

	fmt.Printf("table.persist: %v\n", data)

	// 持久化
	t.Id, err = t.tbm.VerManage().Insert(txId, data)
	return
}
