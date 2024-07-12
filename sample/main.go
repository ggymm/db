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
	tm  tx.Manage
	dm  data.Manage
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
}

// 同步数据到磁盘
func exit() {
	tm.Close()
	dm.Close()

	// 退出
	os.Exit(0)
}

func output(s ...string) {
	fmt.Println(s)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var (
		err error

		in   string
		stmt sql.Statement
	)

	for {
		output("db> ")
		in, err = reader.ReadString(';')
		if err != nil {
			output("Error reading input:", err.Error())
			continue
		}

		in = strings.TrimSpace(in)
		if in == "exit;" {
			exit()
		}

		// 解析 sql 语句
		stmt, err = sql.ParseSQL(in)
		if err != nil {
			output("Error parsing input sql:", err.Error())
			continue
		}

		switch stmt.StmtType() {
		case sql.Select:
		default:
			println("Error Unsupported stmt type:", stmt.StmtType())
		}
	}
}
