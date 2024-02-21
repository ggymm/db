package table

import (
	"fmt"
	"testing"
)

func TestView_Border(t *testing.T) {
	var (
		lens                = []int{10, 11, 12} // 3 columns
		top, middle, bottom string
	)

	top, middle, bottom = border(lens, asciiLine)
	t.Log(top)
	t.Log(middle)
	t.Log(bottom)

	top, middle, bottom = border(lens, singleLine)
	t.Log(top)
	t.Log(middle)
	t.Log(bottom)

	top, middle, bottom = border(lens, doubleLine)
	t.Log(top)
	t.Log(middle)
	t.Log(bottom)
}

func TestView_Print(t *testing.T) {
	v := newView()
	v.setHeader([]string{"Field", "Type", "Null", "Key", "Default", "Extra"})
	v.setValues([][]string{
		{"create_time", "datetime", "YES", "", "NULL", ""},
		{"create_id", "bigint(20)", "YES", "", "NULL", ""},
		{"update_time", "datetime", "YES", "", "NULL", ""},
		{"update_id", "bigint(20)", "YES", "", "NULL", ""},
		{"del_flag", "int(11)", "YES", "", "NULL", ""},
	})

	fmt.Println(v.string(asciiLine))
	fmt.Println(v.string(singleLine))
	fmt.Println(v.string(doubleLine))
}
