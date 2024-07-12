package sql

import (
	"fmt"
	"strings"
)

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

func ParseSQL(sql string) (Statement, error) {
	lex := &Lexer{
		sql: sql,
	}

	r := yyParse(lex)
	if r != 0 {
		return nil, fmt.Errorf("parse sql error %v", lex.errs)
	}

	if len(lex.stmts) == 0 {
		return nil, fmt.Errorf("parse sql error")
	}
	return lex.stmts[0], nil
}

func TrimQuote(str string) (string, error) {
	end := len(str) - 1
	switch str[0] {
	case '`':
		if str[end] != '`' {
			return "", fmt.Errorf("%s missing back quote", str)
		}
		return str[1:end], nil
	case '\'':
		if str[end] != '\'' {
			return "", fmt.Errorf("%s missing single quote", str)
		}
		return str[1:end], nil
	case '"':
		if str[end] != '"' {
			return "", fmt.Errorf("%s missing double quote", str)
		}
		return str[1:end], nil
	}
	return str, nil
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

// Lex 是词法分析器的实现，用于识别并提取 SQL 语句中的 token
// 此函数会被 yacc 生成的代码调用
func (l *Lexer) Lex(val *yySymType) int {
	// start 和 finish 用于标记 token 的起始和结束位置
	start := l.offset
	finish := len(l.sql)
	if l.offset >= finish {
		return 0
	}

	prevQuote := false        // 用于标记是否在引号内
	prevBacktick := false     // 用于标记是否在反引号内
	prevSingleQuotes := false // 用于标记是否在单引号内
	prevDoubleQuotes := false // 用于标记是否在双引号内
	for i := l.offset; i < len(l.sql); i++ {
		switch l.sql[i] {
		case '\\':
			continue
		case '`':
			prevBacktick = !prevBacktick
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		case '\'':
			prevSingleQuotes = !prevSingleQuotes
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		case '"':
			prevDoubleQuotes = !prevDoubleQuotes
			prevQuote = prevBacktick || prevSingleQuotes || prevDoubleQuotes
		}

		// 如果不在引号内
		// 那么根据当前字符判断是否结束 token
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
