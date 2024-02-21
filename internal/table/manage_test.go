package table

import (
	"fmt"
	"path/filepath"
	"testing"

	"db/internal/boot"
	"db/internal/data"
	"db/internal/opt"
	"db/internal/tx"
	"db/internal/ver"
	"db/pkg/bin"
	"db/pkg/sql"
	"db/pkg/utils"
	"db/test"
)

func testOpen() Manage {
	base := utils.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/table")

	b := boot.New(&opt.Option{
		Open: true,
		Name: name,
		Path: path,
	})

	txManage := tx.NewManager(&opt.Option{
		Open: true,
		Name: name,
		Path: path,
	})
	dataManage := data.NewManage(txManage, &opt.Option{
		Open:   true,
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	})

	return NewManage(b, ver.NewManage(txManage, dataManage), dataManage)
}

func testCreate() Manage {
	base := utils.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/table")

	b := boot.New(&opt.Option{
		Open: false,
		Name: name,
		Path: path,
	})

	// 初始化boot
	buf := make([]byte, 8)
	bin.PutUint64(buf, 0)
	b.Update(buf)

	txManage := tx.NewManager(&opt.Option{
		Open: false,
		Name: name,
		Path: path,
	})
	dataManage := data.NewManage(txManage, &opt.Option{
		Open:   false,
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	})

	return NewManage(b, ver.NewManage(txManage, dataManage), dataManage)
}

func TestTableManage_Show(t *testing.T) {
	tbm := testOpen()

	// 展示表
	fmt.Println(tbm.ShowTable())
}

func TestTableManage_Create(t *testing.T) {
	tbm := testCreate()
	stmts, err := sql.ParseSQL(test.CreateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	tbl := stmts[0].(*sql.CreateStmt)

	// 创建表
	txId := tbm.Begin(0)
	err = tbm.Create(txId, tbl)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = tbm.Commit(txId)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// 展示表
	fmt.Println(tbm.ShowTable())

	// 展示字段
	fmt.Println(tbm.ShowField(tbl.Name))
}
