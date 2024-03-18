package view

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
	v := NewTable()
	v.SetHead([]string{"Field", "Type", "Null", "Key", "Default", "Extra"})
	v.SetBody([][]string{
		{"create_time", "datetime", "YES", "", "NULL", ""},
		{"create_id", "bigint(20)", "YES", "", "NULL", ""},
		{"update_time", "datetime", "YES", "", "NULL", ""},
		{"update_id", "bigint(20)", "YES", "", "NULL", ""},
		{"del_flag", "int(11)", "YES", "", "NULL", ""},
	})

	v.SetChars(asciiLine)
	fmt.Println(v.String())
	v.SetChars(singleLine)
	fmt.Println(v.String())
	v.SetChars(doubleLine)
	fmt.Println(v.String())
}
