package main

import (
	_ "embed"

	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"db/internal/boot"
	"db/internal/data"
	"db/internal/opt"
	"db/internal/table"
	"db/internal/tx"
	"db/internal/ver"
	"db/pkg/sql"
	"db/pkg/utils"
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
	base := utils.RunPath()
	path := filepath.Join(base, name)

	if !utils.IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 判断目录是否为空
	cfg := &opt.Option{
		Name:   name,
		Path:   path,
		Memory: (1 << 20) * 64,
	}
	if !utils.IsEmpty(path) {
		cfg.Open = true
	} else {
		cfg.Open = false
	}

	tm = tx.NewManager(cfg)
	dm = data.NewManage(tm, cfg)

	tbm = table.NewManage(boot.New(cfg), ver.NewManage(tm, dm), dm)

	// 初始化表
	if !cfg.Open {
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
		err = tbm.Insert(tx.Super, stmt.(*sql.InsertStmt))
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

func printErr(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
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
			printErr("Error reading input:", err)
			continue
		}

		in = strings.TrimSpace(in)
		if in == "exit" {
			exit()
			break
		}
		if in == "select" {
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
			fmt.Println(tbm.ShowResult(stmt.TableName(), entries))
		}

		// 解析 sql 语句
		stmt, err = sql.ParseSQL(in)
		if err != nil {
			printErr("Error parsing sql:", err)
			continue
		}

		switch stmt.StmtType() {
		case sql.Select:
		default:
			printErr("Unsupported statement type:", stmt.StmtType())
		}

	}
}
