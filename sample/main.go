package main

import (
	_ "embed"

	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ggymm/db"
	"github.com/ggymm/db/boot"
	"github.com/ggymm/db/data"
	"github.com/ggymm/db/pkg/file"
	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/table"
	"github.com/ggymm/db/tx"
	"github.com/ggymm/db/ver"
)

var (
	tm tx.Manage
	dm data.Manage

	tbm table.Manage
)

//go:embed sample_data.sql
var sampleData string

//go:embed sample_struct.sql
var sampleStruct string

// 初始化目录
// 创建基础数据库
func init() {
	// 创建基础数据库
	name := "sample"
	base := db.RunPath()
	path := filepath.Join(base, name)

	if !file.IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	opt := &db.Option{
		Path:   path,
		Memory: (1 << 20) * 64,
	}

	// 判断目录是否为空
	if !file.IsEmpty(path) {
		opt.Open = true
	} else {
		opt.Open = false
	}

	tm = tx.NewManager(opt)
	dm = data.NewManage(tm, opt)
	tbm = table.NewManage(boot.New(opt), ver.NewManage(tm, dm), dm)

	// 初始化表
	if !opt.Open {
		stmt, err := sql.ParseSQL(sampleStruct)
		if err != nil {
			panic(err)
		}
		err = tbm.Create(tx.Super, stmt.(*sql.CreateStmt))
		if err != nil {
			panic(err)
		}

		// 初始化数据
		stmt, err = sql.ParseSQL(sampleData)
		if err != nil {
			panic(err)
		}
		txId := tbm.Begin(0)
		err = tbm.Insert(txId, stmt.(*sql.InsertStmt))
		if err != nil {
			panic(err)
		}
		err = tbm.Commit(txId)
		if err != nil {
			panic(err)
		}
	}
}

// 同步数据到磁盘
func exit() {
	tm.Close()
	dm.Close()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var (
		err error

		in   string
		stmt sql.Statement
	)

	for {
		fmt.Print("db> ")
		in, err = reader.ReadString(';')
		if err != nil {
			println("Error reading input:", err)
			continue
		}

		in = strings.TrimSpace(in)
		if in == "exit;" {
			exit()
			break
		}
		if in == "test;" {
			// 测试
			stmt, err = sql.ParseSQL("select * from user;")
			txId := tbm.Begin(0)
			entries, err := tbm.Select(txId, stmt.(*sql.SelectStmt))
			if err != nil {
				panic(err)
			}
			err = tbm.Commit(txId)
			if err != nil {
				panic(err)
			}
			println(tbm.ShowResult(stmt.TableName(), entries))
			continue
		}

		// 解析 sql 语句
		stmt, err = sql.ParseSQL(in)
		if err != nil {
			println("Error parsing sql:", err)
			continue
		}

		switch stmt.StmtType() {
		case sql.Select:
		default:
			println("Error Unsupported statement type:", stmt.StmtType())
		}
	}
}
