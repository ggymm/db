package table

import (
	"db/internal/index"
	"db/internal/opt"
	"db/internal/tx"
	"db/pkg/bin"
	"db/pkg/str"
)

// 字段
//
// 字段内存中数据结构如下：
// +----------------+----------------+----------------+
// |     name       |     type       |     index      |
// +----------------+----------------+----------------+
// |    string      |    string      |     uint64     |
// +----------------+----------------+----------------+
//

type field struct {
	tbm Manage
	idx index.Index

	Id    uint64
	Name  string
	Type  string
	Index uint64
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
	f.Name, shift = str.Deserialize(data)

	// 读取 type
	pos += shift
	f.Type, shift = str.Deserialize(data[pos:])

	// 读取 index
	pos += shift
	f.Index = bin.Uint64(data[pos:])
	if f.Index != 0 {
		f.idx, err = index.NewIndex(tbm.DataManage(), &opt.Option{
			Open:   true,
			RootId: f.Index,
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

		Name:  info.Name,
		Type:  info.Type,
		Index: 0,
	}

	if info.Indexed {
		idx, err := index.NewIndex(tbm.DataManage(), &opt.Option{
			Open: false,
		})
		if err != nil {
			return nil, err
		}
		f.Index = idx.GetRootId()
	}
	return f, f.persist(info.TxId)
}

func (f *field) persist(txId uint64) (err error) {
	// name
	data := str.Serialize(f.Name)

	// type
	data = append(data, str.Serialize(f.Type)...)

	// index
	buf := make([]byte, 8)
	bin.PutUint64(buf, f.Index)
	data = append(data, buf...)

	// 持久化
	f.Id, err = f.tbm.VerManage().Insert(txId, data)
	return
}

func (f *field) String() string {
	return "{" + f.Name + ": " + f.Type + "} "
}
