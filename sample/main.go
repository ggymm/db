package main

import (
	_ "embed"

	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ggymm/db"
	"github.com/ggymm/db/boot"
	"github.com/ggymm/db/data"
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

// 创建基础数据库
func init() {
	opt := db.NewOption("temp")
	opt.Memory = (1 << 20) * 64

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

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("db> ")
		in, err := reader.ReadString(';')
		if err != nil {
			fmt.Println("Error reading input:", err.Error())
			continue
		}

		in = strings.TrimSpace(in)
		if in == "exit;" {
			exit()
		}

		in = strings.Replace(in, "\r", " ", -1)
		in = strings.Replace(in, "\n", " ", -1)

		if in == "show tables;" {
			fmt.Println(tbm.ShowTable())
			continue
		}

		// 解析 sql 语句
		stmt, err := sql.ParseSQL(in)
		if err != nil {
			fmt.Println("Error parsing input sql:", err.Error())
			continue
		}

		tid := tbm.Begin(1) // 可重复读
		switch stmt.StmtType() {
		case sql.Create:
			err = tbm.Create(tid, stmt.(*sql.CreateStmt))
		case sql.Insert:
			err = tbm.Insert(tid, stmt.(*sql.InsertStmt))
		case sql.Update:
			err = tbm.Update(tid, stmt.(*sql.UpdateStmt))
		case sql.Delete:
			err = tbm.Delete(tid, stmt.(*sql.DeleteStmt))
		case sql.Select:
			var entries []table.Entry
			entries, err = tbm.Select(tid, stmt.(*sql.SelectStmt))
			fmt.Println(tbm.ShowResult(stmt.TableName(), entries))
		default:
			println("Error Unsupported stmt type:", stmt.StmtType())
		}
		if err != nil {
			tbm.Rollback(tid)
			fmt.Println("Error exec sql:", err.Error())
			continue
		}
		err = tbm.Commit(tid)
		if err != nil {
			fmt.Println("Error exec sql:", err.Error())
			continue
		}
		fmt.Println("ok")
	}
}
