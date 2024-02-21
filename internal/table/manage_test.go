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

func TestTableManage_Show(t *testing.T) {
	base := utils.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/table")

	b := boot.New(&opt.Option{
		Open: true,
		Name: name,
		Path: path,
	})

	tm := tx.NewManager(&opt.Option{
		Open: true,
		Name: name,
		Path: path,
	})
	dm := data.NewManage(tm, &opt.Option{
		Open:   true,
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	})
	tbm := NewManage(b, ver.NewManage(tm, dm), dm)

	// 展示表
	fmt.Println(tbm.ShowTable())
}

func TestTableManage_Create(t *testing.T) {
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

	tm := tx.NewManager(&opt.Option{
		Open: false,
		Name: name,
		Path: path,
	})
	dm := data.NewManage(tm, &opt.Option{
		Open:   false,
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	})
	tbm := NewManage(b, ver.NewManage(tm, dm), dm)

	// 解析创建表语句
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

	// 释放资源
	// 同步数据到磁盘
	tm.Close()
	dm.Close()
}
