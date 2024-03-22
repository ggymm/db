package view

import "testing"

func Test_CalcWidth(t *testing.T) {
	t.Log(runeWidth("中文"))
	t.Log(runeWidth("phone_number"))
	t.Log(runeWidth("更新时间10"))
}
