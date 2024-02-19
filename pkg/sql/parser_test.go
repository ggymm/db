package sql

import (
	"bufio"
	_ "embed"
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
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
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
			t.Logf("%d %+v", stmt.GetStmtType(), stmt)
		}
	}
}
