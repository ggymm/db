package main

import (
	"bufio"
	"db/internal/boot"
	"db/internal/data"
	"db/internal/opt"
	"db/internal/table"
	"db/internal/tx"
	"db/internal/ver"
	"db/pkg/sql"
	"db/pkg/utils"
	"db/test"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	// 创建基础数据库
}

func openTbm() table.Manage {
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
	return table.NewManage(b, ver.NewManage(tm, dm), dm)
}

// 同步数据到磁盘
func closeTbm(tbm table.Manage) {
	tbm.DataManage().TxManage().Close()
	tbm.DataManage().Close()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("db> ")
		input, err := reader.ReadString(';')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		// 用户输入exit时退出
		if input == "select" {
			tbm := openTbm()

			// 解析查询表语句
			stmt, err := sql.ParseSQL(test.SelectAllSQL)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error parsing SQL:", err)
				continue
			}

			txId := tbm.Begin(0)
			entries, err := tbm.Select(txId, stmt.(*sql.SelectStmt))
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error selecting data:", err)
				continue
			}
			err = tbm.Commit(txId)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error committing transaction:", err)
				continue
			}

			// 展示字段
			fmt.Println(tbm.ShowResult(stmt.TableName(), entries))

			// 释放资源
			closeTbm(tbm)
			break
		}

		// 在这里处理其他命令
		fmt.Printf("You entered: %s\n", input)
	}
}
