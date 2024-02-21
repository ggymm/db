package sql

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"strings"
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

func TestParseSQL_DML(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(test.SQLList))
	for scanner.Scan() {
		stmts, err := ParseSQL(scanner.Text())
		if err != nil {
			t.Fatalf("%+v", err)
		}
		for _, stmt := range stmts {
			s, _ := json.MarshalIndent(stmt, "", "  ")
			t.Logf("stmtType: %d\n%s", stmt.Type(), s)
		}
	}
}
