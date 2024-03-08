package sql

import (
	_ "embed"
	"encoding/json"
	"testing"

	"db/test"
)

func TestParseSQL_DDL(t *testing.T) {
	stmts, err := ParseSQL(test.CreateSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, stmt := range stmts {
		s, _ := json.MarshalIndent(stmt, "", "  ")
		t.Log(s)
	}
}

func TestParseSQL_Select(t *testing.T) {
	stmts, err := ParseSQL(test.SelectSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	stmt := stmts[0].(*SelectStmt)
	s, _ := json.MarshalIndent(stmt, "", "  ")
	t.Logf("stmtType: %d\n%s", stmt.Type(), s)
}
