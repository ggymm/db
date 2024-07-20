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

	entry := []map[string]any{
		{"id": "1", "username": "名称1", "nickname": "昵称1", "email": "邮箱1", "extras": "1"},
		{"id": "2", "username": "名称2", "nickname": "昵称2", "email": "邮箱2", "extras": "2"},
		{"id": "3", "username": "名称3", "nickname": "昵称3", "email": "邮箱3", "extras": "3"},
		{"id": "4", "username": "名称4", "nickname": "昵称4", "email": "邮箱4", "extras": "4"},
		{"id": "5", "username": "名称5", "nickname": "昵称5", "email": "邮箱5", "extras": "5"},
	}

	indexes := make([]int, 0)
	selectStmt := stmt.(*SelectStmt)

	for i := 0; i < len(entry); i++ {
		exist := true
		for _, where := range selectStmt.Where {
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if where.Match(entry[i]) && exist {
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
