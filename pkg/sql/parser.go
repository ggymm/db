package sql

import "strings"

var mapping map[string]int

func init() {
	singleCharToken := []string{">", "<", "=", ",", ";", "(", ")"}
	mapping = make(map[string]int, len(yyTokenLiteralStrings)+len(singleCharToken))

	for k, v := range yyTokenLiteralStrings {
		mapping[v] = k
		mapping[strings.ToLower(v)] = k
	}

	for _, v := range singleCharToken {
		mapping[v] = int(int8(v[0]))
	}
}

type Lexer struct {
	sql    string
	stmts  []Statement
	offset int
	errs   []string
}

func (l *Lexer) Error(s string) {
	l.errs = append(l.errs, s)
}

func (l *Lexer) Lex(val *yySymType) int {
	start := l.offset
	finish := len(l.sql)
	if l.offset >= finish {
		return 0
	}

	prevQuote := false
	prevBacktick := false
	prevSingleQuotes := false
	prevDoubleQuotes := false

	for i := l.offset; i < len(l.sql); i++ {
		switch l.sql[i] {
		case '\\':
			continue
		case '\'':
			prevSingleQuotes = !prevSingleQuotes
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		case '"':
			prevDoubleQuotes = !prevDoubleQuotes
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		case '`':
			prevBacktick = !prevBacktick
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		}

		if !prevQuote {
			switch l.sql[i] {
			case ' ', '\n', '\t':
				finish = i
			case '<', '>', '!':
				finish = i
				if start == finish {
					if i+1 < len(l.sql) && l.sql[i+1] == '=' {
						finish += 2
					} else {
						finish++
					}
				}
			case '=', ',', ';', '(', ')':
				finish = i
				if start == finish {
					finish++
				}
			}
		}

		if i == finish {
			break
		}
	}

	l.offset = finish
	for l.offset < len(l.sql) {
		char := l.sql[l.offset]
		if char == ' ' || char == '\n' || char == '\t' {
			l.offset++
		} else {
			break
		}
	}

	token := l.sql[start:finish]
	val.str = token

	num, ok := mapping[token]
	if ok {
		return num
	} else {
		return VARIABLE
	}
}

func (l *Lexer) Reduced(_, state int, val *yySymType) bool {
	if state == 2 {
		l.stmts = val.stmtList
	}
	return false
}
