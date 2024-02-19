package sql

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"strings"
	"testing"
)

//go:embed test_ddl.sql
var testDDLSQL string

//go:embed test_dml.sql
var testDMLSQL string

func TestParseSQL_DDL(t *testing.T) {
	stmts, err := ParseSQL(testDDLSQL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, stmt := range stmts {
		s, _ := json.MarshalIndent(stmt, "", "  ")
		t.Log(s)
	}
}

func TestParseSQL_DML(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(testDMLSQL))
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
