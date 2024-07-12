package sql

import (
	"fmt"
)

func ParseSQL(sql string) (Statement, error) {
	ss, err := ParseMultiSQL(sql)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, fmt.Errorf("解析sql为空")
	}
	return ss[0], nil
}

func ParseMultiSQL(sql string) ([]Statement, error) {
	lex := &Lexer{
		sql: sql,
	}

	ret := yyParse(lex)
	if ret != 0 {
		return nil, fmt.Errorf("解析sql异常 %+v", lex.errs)
	}
	return lex.stmts, nil
}

func TrimQuote(str string) (string, error) {
	end := len(str) - 1
	switch str[0] {
	case '\'':
		if str[end] != '\'' {
			return "", fmt.Errorf("%s 缺少单引号", str)
		}
		return str[1:end], nil
	case '"':
		if str[end] != '"' {
			return "", fmt.Errorf("%s 缺少双引号", str)
		}
		return str[1:end], nil
	case '`':
		if str[end] != '`' {
			return "", fmt.Errorf("%s 缺少反引号", str)
		}
		return str[1:end], nil
	}
	return str, nil
}
