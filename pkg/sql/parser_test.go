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

func TestParseSQL_Select(t *testing.T) {
	stmt, err := ParseSQL(test.SelectAllSQL)
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
