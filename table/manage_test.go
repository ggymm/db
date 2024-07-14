package table_test

import (
	"fmt"
	"testing"

	"github.com/ggymm/db"
	"github.com/ggymm/db/boot"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/table"
	"github.com/ggymm/db/test"
	"github.com/ggymm/db/tx"
	"github.com/ggymm/db/ver"
)

var (
	tm tx.Manage
	dm data.Manage
)

func openTbm() table.Manage {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/table")
	opt.Memory = (1 << 20) * 64

	b := boot.New(opt)
	tm = tx.NewManager(opt)
	dm = data.NewManage(tm, opt)
	return table.NewManage(b, ver.NewManage(tm, dm), dm)
}

// 同步数据到磁盘
func closeTbm() {
	tm.Close()
	dm.Close()
}

func TestTableManage_Show(t *testing.T) {
	tbm := openTbm()

	// 展示表
	fmt.Println(tbm.ShowTable())
}

func TestTableManage_Create(t *testing.T) {
	tbm := openTbm()

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

	// 展示字段
	fmt.Println(tbm.ShowField(stmt.TableName()))

	// 释放资源
	closeTbm()
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
	closeTbm()
}

func TestTableManage_Select(t *testing.T) {
	tbm := openTbm()

	// 解析查询表语句
	stmt, err := sql.ParseSQL(test.SelectAllSQL)
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

	// 展示字段
	fmt.Println(tbm.ShowResult(stmt.TableName(), entries))

	// 释放资源
	closeTbm()
}
