package table

import (
	"encoding/json"
	"testing"

	"github.com/ggymm/db/pkg/sql"
	"github.com/ggymm/db/test"
)

func Test_fmtCond(t *testing.T) {
	src := []*Interval{
		{Min: 1, Max: 20},
		{Min: 30, Max: 40},
		{Min: 20, Max: 30},
		{Min: 20, Max: 50},
	}
	dst := NewExplain().format(src)
	t.Logf("%+v", dst)
}

func Test_mixCond(t *testing.T) {
	s0 := []*Interval{
		{Min: 1, Max: 20},
	}
	s1 := []*Interval{
		{Min: 30, Max: 40},
	}
	dst := NewExplain().compact(s0, s1)
	t.Logf("%+v", dst)
}

func Test_SelectIndex(t *testing.T) {
	list := []string{
		//"index/select1.sql",
		//"index/select2.sql",
		//"index/select3.sql",
		//"index/select4.sql",
		"index/select5.sql",
	}
	for _, item := range list {
		b, err := test.SelectIndexSQL.ReadFile(item)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		stmt, err := sql.ParseSQL(string(b))
		if err != nil {
			t.Fatalf("%+v", err)
		}
		s, _ := json.MarshalIndent(stmt, "", "  ")
		t.Logf("%s", s)

		selectStmt := stmt.(*sql.SelectStmt)
		res, err := NewExplain().Execute(&field{
			Name: "id",
			Type: sql.Int64.String(),
		}, selectStmt.Where)

		// 打印结果
		t.Logf("%s %+v", string(b), err)
		for _, r := range res {
			t.Logf("%+v", r)
		}
	}
}
