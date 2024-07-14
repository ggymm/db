package sql

import (
	_ "embed"

	"encoding/json"
	"testing"

	"github.com/ggymm/db/test"
)

func TestParseSQL_Tx(t *testing.T) {
	stmt, err := ParseSQL("BEGIN 0;")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)

	stmt, err = ParseSQL("COMMIT;")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ = json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)

	stmt, err = ParseSQL("ROLLBACK;")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ = json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_Create(t *testing.T) {
	stmt, err := ParseSQL(test.CreateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_Insert(t *testing.T) {
	stmt, err := ParseSQL(test.InsertSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_Update(t *testing.T) {
	stmt, err := ParseSQL(test.UpdateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_Delete(t *testing.T) {
	stmt, err := ParseSQL(test.DeleteSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_Select(t *testing.T) {
	stmt, err := ParseSQL(test.SelectSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)
}

func TestParseSQL_SelectWhere(t *testing.T) {
	stmt, err := ParseSQL(test.SelectWhereSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("%s", s)

	selectStmt := stmt.(*SelectStmt)
	// 应该首先将带索引的条件解析出来

	fm := map[string]int{
		"username": 0,
		"nickname": 1,
		"email":    2,
		"extras":   3,
	}
	data := []*[]string{
		{"名称1", "昵称1", "邮箱1", "1"},
		{"名称2", "昵称2", "邮箱2", "2"},
		{"名称3", "昵称3", "邮箱3", "3"},
		{"名称4", "昵称4", "邮箱4", "4"},
		{"名称5", "昵称5", "邮箱5", "5"},
	}

	indexes := make([]int, 0)
	for i := 0; i < len(data); i++ {
		exist := true
		for _, where := range selectStmt.Where {
			err = where.Prepare(fm)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if where.Filter(data[i]) && exist {
				exist = true
			} else {
				exist = false
			}
		}
		if exist {
			indexes = append(indexes, i)
		}
	}

	t.Logf("indexes: %+v", indexes)
}
