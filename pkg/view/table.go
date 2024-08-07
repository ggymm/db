package view

import (
	"strings"
)

var (
	padding = 1

	asciiLine = []string{
		"+", "+", "+",
		"|", "+", "|",
		"+", "+", "+",
		"-",
		"|",
	}
	singleLine = []string{
		"┌", "┬", "┐",
		"├", "┼", "┤",
		"└", "┴", "┘",
		"─",
		"│",
	}
	doubleLine = []string{
		"╔", "╦", "╗",
		"╠", "╬", "╣",
		"╚", "╩", "╝",
		"═",
		"║",
	}
)

func border(lens []int, chars []string) (string, string, string) {
	top := chars[0]
	middle := chars[3]
	bottom := chars[6]
	for i, l := range lens {
		top += strings.Repeat(chars[9], l)
		middle += strings.Repeat(chars[9], l)
		bottom += strings.Repeat(chars[9], l)

		if i == len(lens)-1 {
			top += chars[2]
			middle += chars[5]
			bottom += chars[8]
		} else {
			top += chars[1]
			middle += chars[4]
			bottom += chars[7]
		}
	}
	return top, middle, bottom
}

type Table struct {
	head []string
	body [][]string

	chars []string
}

func NewTable() *Table {
	return &Table{
		chars: singleLine,

		head: []string{},
		body: [][]string{},
	}
}

func (v *Table) calcLens() []int {
	lens := make([]int, len(v.head))
	for i, h := range v.head {
		w := runeWidth(h)
		if w > lens[i] {
			lens[i] = w
		}
	}
	for _, row := range v.body {
		for i, r := range row {
			w := runeWidth(r)
			if w > lens[i] {
				lens[i] = w
			}
		}
	}
	// 加上 padding
	for i, l := range lens {
		lens[i] = l + padding*2
	}
	return lens
}

func (v *Table) SetChars(chars []string) {
	v.chars = chars
}

func (v *Table) SetHead(head []string) {
	v.head = head
}

func (v *Table) SetBody(body [][]string) {
	// 检查是否有 head
	if len(v.head) == 0 {
		panic("no head")
	}
	// 检查每条记录的长度是否和 header 一致
	for _, row := range body {
		if len(row) != len(v.head) {
			panic("invalid body")
		}
	}
	v.body = body
}

func (v *Table) String() string {
	str := ""
	res := make([]string, 0)

	// 计算每列的宽度
	lens := v.calcLens()
	chars := v.chars

	// 左右填充
	left := func() string {
		return strings.Repeat("\u0020", padding)
	}
	right := func(i int, str string) string {
		w := runeWidth(str)
		return strings.Repeat("\u0020", lens[i]-w-padding) + chars[10]
		// return runeFillRight(str, lens[i]-runeWidth(str)-padding) + chars[10]
	}

	// 生成边框
	top, middle, bottom := border(lens, chars)

	// 输出标题
	res = append(res, top)
	str = chars[10]
	for i, h := range v.head {
		str += left() + h + right(i, h)
	}
	res = append(res, str)

	// 输出每条记录
	res = append(res, middle)
	for _, row := range v.body {
		str = chars[10]
		for i, r := range row {
			str += left() + r + right(i, r)
		}
		res = append(res, str)
	}
	return strings.Join(append(res, bottom), "\n")
}
