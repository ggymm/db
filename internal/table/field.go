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
