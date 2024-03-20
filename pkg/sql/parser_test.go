package sql

import (
	_ "embed"
	"encoding/json"
	"testing"

	"db/test"
)

func TestParseSQL_DDL(t *testing.T) {
	stmt, err := ParseSQL(test.CreateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Log(s)
}

func TestParseSQL_Select(t *testing.T) {
	stmt, err := ParseSQL(test.SelectSQL)
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
