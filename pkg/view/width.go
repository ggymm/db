package view

import (
	"strings"
	"unicode/utf8"
)

var trie = newWidthTrie(0)

type Kind int

const (
	Neutral Kind = iota
	EastAsianAmbiguous
	EastAsianWide
	EastAsianNarrow
	EastAsianFullwidth
	EastAsianHalfwidth
)

type elem uint16

const (
	numTypeBits = 3
	typeShift   = 16 - numTypeBits
)

type Properties struct {
	elem elem
	last byte
}

func (e elem) kind() Kind {
	return Kind(e >> typeShift)
}

func (p Properties) Kind() Kind {
	return p.elem.kind()
}

func LookupRune(r rune) Properties {
	var buf [4]byte
	n := utf8.EncodeRune(buf[:], r)
	v, _ := trie.lookup(buf[:n])
	last := byte(r)
	if r >= utf8.RuneSelf {
		last = 0x80 + byte(r&0x3f)
	}
	return Properties{elem(v), last}
}

func runeWidth(s string) int {
	// 计算字符串宽度
	w := 0
	for _, r := range s {
		switch LookupRune(r).Kind() {
		case EastAsianFullwidth, EastAsianWide:
			w += 2
		case EastAsianHalfwidth, EastAsianNarrow,
			Neutral, EastAsianAmbiguous:
			w += 1
		}
	}
	return w
}

func runeFillRight(s string, num int) string {
	hasWide := false
	for _, r := range s {
		k := LookupRune(r).Kind()
		if k == EastAsianWide ||
			k == EastAsianFullwidth {
			hasWide = true
			break
		}
	}
	if hasWide {
		return strings.Repeat("\u0020", num-1) + "\u3000"
	} else {
		return strings.Repeat("\u0020", num)
	}
}
