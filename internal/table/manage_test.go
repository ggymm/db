package table

import (
	"db/pkg/view"
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

func openTbm() Manage {
	base := utils.RunPath()
	name := "test"
	path := filepath.Join(base, "temp/table")

	b := boot.New(&opt.Option{
		Open: false,
		Name: name,
		Path: path,
	})

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
	return NewManage(b, ver.NewManage(tm, dm), dm)
}

// 同步数据到磁盘
func closeTbm(tbm Manage) {
	tbm.DataManage().TxManage().Close()
	tbm.DataManage().Close()
}

func TestTableManage_Show(t *testing.T) {
	tbm := openTbm()

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
	raw := bin.Uint64Raw(0)
	b.Update(raw)

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
	stmt, err := sql.ParseSQL(test.CreateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// 创建表
	txId := tbm.Begin(0)
	err = tbm.Create(txId, stmt.(*sql.CreateStmt))
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
	fmt.Println(tbm.ShowField(stmt.TableName()))

	// 释放资源
	closeTbm(tbm)
}

func TestTableManage_Insert(t *testing.T) {
	tbm := openTbm()

	// 解析创建表语句
	stmt, err := sql.ParseSQL(test.InsertSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// 插入数据
	txId := tbm.Begin(0)
	err = tbm.Insert(txId, stmt.(*sql.InsertStmt))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = tbm.Commit(txId)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// 释放资源
	closeTbm(tbm)
}

func TestTableManage_Select(t *testing.T) {
	tbm := openTbm()

	// 解析创建表语句
	stmt, err := sql.ParseSQL(test.SelectSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	txId := tbm.Begin(0)
	entries, err := tbm.Select(txId, stmt.(*sql.SelectStmt))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = tbm.Commit(txId)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// 打印数据
	thead := make([]string, 0)
	tbody := make([][]string, 0)
	for i, ent := range entries {
		if i == 0 {
			for k := range ent {
				thead = append(thead, k)
			}
		}
		row := make([]string, 0)
		for _, v := range ent {
			row = append(row, fmt.Sprintf("%v", v))
		}
		tbody = append(tbody, row)
	}

	// 表格形式输出
	vt := view.NewTable()
	vt.SetHead(thead)
	vt.SetBody(tbody)

	// 打印表格
	fmt.Println(vt.String())

	// 释放资源
	closeTbm(tbm)
}
