package table

import "strings"

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

type view struct {
	header []string
	values [][]string
}

func newView() *view {
	return &view{
		header: []string{},
		values: [][]string{},
	}
}

func (v *view) setHeader(header []string) {
	v.header = header
}

func (v *view) setValues(values [][]string) {
	// 检查是否有 header
	if len(v.header) == 0 {
		panic("no header")
	}
	// 检查每条记录的长度是否和 header 一致
	for _, row := range values {
		if len(row) != len(v.header) {
			panic("invalid values")
		}
	}
	v.values = values
}

func (v *view) calcColumnLens() []int {
	widths := make([]int, len(v.header))
	for i, h := range v.header {
		if len(h) > widths[i] {
			widths[i] = len(h)
		}
	}
	for _, row := range v.values {
		for i, v := range row {
			if len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}
	// 加上 padding
	for i, w := range widths {
		widths[i] = w + padding*2
	}
	return widths
}

func (v *view) string(chars []string) string {
	str := ""
	res := make([]string, 0)

	// 计算每列的宽度
	lens := v.calcColumnLens()

	left := func() string {
		return strings.Repeat(" ", padding)
	}
	right := func(i int, len int) string {
		return strings.Repeat(" ", lens[i]-len-padding) + chars[10]
	}

	// 生成边框
	top, middle, bottom := border(lens, chars)

	// 输出 header
	res = append(res, top)
	str = chars[10]
	for i, val := range v.header {
		str += left() + val + right(i, len(val))
	}
	res = append(res, str)

	// 输出每条记录
	res = append(res, middle)
	for _, row := range v.values {
		str = chars[10]
		for i, val := range row {
			str += left() + val + right(i, len(val))
		}
		res = append(res, str)
	}
	res = append(res, bottom)
	return strings.Join(res, "\n")
}
