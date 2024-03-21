package sql

import (
	"testing"
)

func TestField(t *testing.T) {
	var (
		v     any
		shift int

		v1 = "100"
		v2 = "100"
		v3 = "100"
		v4 = "中文测试"
		v5 any
	)

	f1 := FieldRaw("INT32", v1)
	v, shift = FieldParse("INT32", f1)
	if shift == len(f1) {
		t.Log("FieldParse Uint32 success")
		t.Logf("v: %v, v1: %s, shift: %v", v, v1, shift)
	} else {
		t.Errorf("FieldParse Uint32 failed %v %v", v, shift)
	}

	f2 := FieldRaw("INT64", v2)
	v, shift = FieldParse("INT64", f2)
	if shift == len(f2) {
		t.Log("FieldParse Uint64 success")
		t.Logf("v: %v, v2: %s, shift: %v", v, v2, shift)
	} else {
		t.Errorf("FieldParse Uint64 failed %v %v", v, shift)
	}

	f3 := FieldRaw("VARCHAR", v3)
	v, shift = FieldParse("VARCHAR", f3)
	if v == v3 && shift == len(f3) {
		t.Logf("FieldParse String success")
	} else {
		t.Errorf("FieldParse String failed %v %v", v, shift)
	}

	f4 := FieldRaw("VARCHAR", v4)
	v, shift = FieldParse("VARCHAR", f4)
	if v == v4 && shift == len(f4) {
		t.Logf("FieldParse String success")
	} else {
		t.Errorf("FieldParse String failed %v %v", v, shift)
	}

	f5 := FieldRaw("VARCHAR", v5)
	v, shift = FieldParse("VARCHAR", f5)
	if v == v5 && shift == len(f5) {
		t.Logf("FieldParse String success")
	} else {
		t.Errorf("FieldParse String failed %v %v", v, shift)
	}
}
