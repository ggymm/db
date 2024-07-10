package table

import (
	"fmt"
	"testing"

	"db"
	"db/boot"
	"db/data"
	"db/pkg/sql"
	"db/test"
	"db/tx"
	"db/ver"
)

func openTbm() Manage {
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/table")
	opt.Memory = (1 << 20) * 64

	b := boot.New(opt)
	tm := tx.NewManager(opt)
	dm := data.NewManage(tm, opt)
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
	abs := db.RunPath()
	opt := db.NewOption(abs, "temp/table")
	opt.Memory = (1 << 20) * 64

	b := boot.New(opt)
	tm := tx.NewManager(opt)
	dm := data.NewManage(tm, opt)
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
	closeTbm(tbm)
}
