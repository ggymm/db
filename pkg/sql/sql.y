%{
package sql

import (
	"strconv"
)

%}

%union {
	str string
	strList []string
	boolean bool
	fieldType FieldType
	compareOperate CompareOperate

	stmt Statement
	stmtList []Statement

	createStmt *CreateStmt
	createTable *CreateTable
	createField *CreateField
	createIndex *CreateIndex
	createTableOption *CreateTableOption

	selectStmt *SelectStmt
	selectFieldList []*SelectField
	selectFrom *SelectFrom
	selectWhereList []SelectWhere
	selectOrderList []*SelectOrder
	selectLimit *SelectLimit

	insertStmt *InsertStmt
    valueList [][]string
}

%token <str>
	// 关键字（创建表）
	CREATE "CREATE"
	TABLE "TABLE"
	KEY "KEY"
	NOT	"NOT"
	NULL "NULL"
	INDEX "INDEX"
	DEFAULT "DEFAULT"
	PRIMARY "PRIMARY"
	// 关键字（查询表）
	SELECT "SELECT"
	FROM "FROM"
	WHERE "WHERE"
	OR "OR"
	AND "AND"
	ORDER "ORDER"
	BY "BY"
	ASC "ASC"
	DESC "DESC"
	LIMIT "LIMIT"
	OFFSET "OFFSET"
	// 关键字（插入表）
	INSERT "INSERT"
	INTO "INTO"
	VALUE "VALUE"
	VALUES "VALUES"

%token <str>
	COMP_NE "!="
	COMP_LE "<="
	COMP_GE ">="

%token <str> VARIABLE

%type <str> Expr
%type <strList> VaribleList

%type <stmt> Stmt
%type <stmtList> StmtList

// 语法定义（创建表）
%type <str> DefaultValue
%type <boolean> AllowNull
%type <fieldType> FieldType

%type <createStmt> CreateStmt
%type <createTable> CreateTable
%type <createField> CreateField
%type <createIndex> CreateIndex
%type <createIndex> CreatePrimary
%type <createTableOption> CreateTableOption

// 语法定义（查询表）
%type <boolean> Ascend
%type <compareOperate> CompareOperate

%type <selectStmt> SelectStmt
%type <selectFrom> SelectFrom
%type <selectFieldList> SelectFieldList
%type <selectWhereList> SelectWhere SelectWhereList
%type <selectOrderList> SelectOrder SelectOrderList
%type <selectLimit> SelectLimit

// 语法定义（插入表）
%type <insertStmt> InsertStmt
%type <str> InsertTable
%type <strList> InsertField InsertFieldList
%type <valueList> InsertValue InsertValueList


%left OR
%left AND
%left '+' '-'
%left '*' '/'

%start start

%%

start:
	StmtList

Expr:
	VARIABLE
	{
		str, err := TrimQuote($1)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		$$ = str
	}

VaribleList:
	Expr
	{
		$$ = []string{ $1 }
	}
	| VaribleList ',' Expr
	{
		$$ = append($1, $3)
	}

Stmt:
	CreateStmt
	{
		$$ = Statement($1)
	}
	| SelectStmt
	{
		$$ = Statement($1)
	}
	| InsertStmt
	{
		$$ = Statement($1)
	}

StmtList:
	Stmt
	{
		$$ = append($$, $1)
	}
	| StmtList Stmt
	{
		$$ = append($$, $2)
	}

// 语法规则（创建表）
AllowNull:
	{
		$$ = true
	}
	| "NULL"
	{
		$$ = true
	}
	| "NOT" "NULL"
	{
		$$ = false
	}

DefaultValue:
	{
		$$ = ""
	}
	| "DEFAULT"
	{
		$$ = ""
	}
	| "DEFAULT" "NULL"
	{
		$$ = ""
	}
	| "DEFAULT" Expr
   	{
	   	$$ = $2
   	}

FieldType:
	Expr
	{
		t, ok := typeMapping[$1]
		if ok {
			$$ = t
		} else {
			__yyfmt__.Printf("不支持的数据类型 %s",$1)
			goto ret1
		}
	}

CreateStmt:
	"CREATE" "TABLE" Expr '(' CreateTable ')' CreateTableOption ';'
	{
		$$ = &CreateStmt{
			Name: $3,
			Table: $5,
			Option: $7,
		}
	}

CreateTable:
	CreateField
	{
		$$ = &CreateTable{
			Field: []*CreateField{$1},
			Index: []*CreateIndex{},
		}
	}
	| CreateIndex
	{
		$$ = &CreateTable{
			Field: []*CreateField{},
			Index: []*CreateIndex{$1},
		}
	}
	| CreatePrimary
	{
		$$ = &CreateTable{
			Pk: $1,
			Field: []*CreateField{},
			Index: []*CreateIndex{},
		}
	}
	| CreateTable ',' CreateField
	{
		$$.Field = append($$.Field, $3)
	}
	| CreateTable ',' CreateIndex
	{
		$$.Index = append($$.Index, $3)
	}
	| CreateTable ',' CreatePrimary
	{
		if $$.Pk == nil {
			$$.Pk = $3
		} else {
			 __yyfmt__.Printf("重复定义主键 %v %v ", $$.Pk, $3)
			goto ret1
		}
	}

CreateField:
	Expr FieldType AllowNull DefaultValue
	{
		$$ = &CreateField{
			Name: $1,
			Type: $2,
			AllowNull: $3,
			DefaultValue: $4,
		}
	}

CreateIndex:
	"INDEX" Expr '(' VaribleList ')'
	{
		$$ = &CreateIndex{
			Name: $2,
			Field: $4,
		}
	}

CreatePrimary:
	"PRIMARY" "KEY" '(' VaribleList ')'
	{
		$$ = &CreateIndex{
			Pk: true,
			Field: $4,
		}
	}

CreateTableOption:
	{
		$$ = nil
	}

// 语法规则（查询表）
Ascend:
    {
        $$ = true
    }
    | "ASC"
    {
        $$ = true
    }
    | "DESC" {
        $$ = false
    }

CompareOperate:
    '='
    {
        $$ = EQ
    }
    | '<'
    {
        $$ = LT
    }
    | '>'
    {
        $$ = GT
    }
    | "<="
    {
        $$ = LE
    }
    | ">="
    {
        $$ = GE
    }
    | "!="
    {
        $$ = NE
    }

SelectStmt:
    "SELECT" SelectFieldList SelectLimit ';'
    {
        $$ = &SelectStmt{
        	Field: $2,
        	Limit: $3,
        }
    }
    |  "SELECT" SelectFieldList SelectFrom SelectWhere SelectOrder SelectLimit ';'
    {
        $$ = &SelectStmt{
        	From: $3,
        	Field: $2,
        	Where: $4,
        	Order: $5,
        	Limit: $6,
        }
    }

SelectFieldList:
   Expr
   {
	   $$ = []*SelectField{
		   &SelectField{
			   Name: $1,
		   },
	   }
   }
   | SelectFieldList ',' Expr
   {
	   $$ = append($1, &SelectField{
		   Name: $3,
	   })
   }

SelectFrom:
	"FROM" Expr
	{
		$$ = &SelectFrom{
			Name: $2,
		}
	}

SelectWhere:
	{
        $$ = nil
    }
    | "WHERE" SelectWhereList
    {
        $$ = $2
    }

SelectWhereList:
	Expr CompareOperate Expr
	{
		$$ = []SelectWhere{
			&SelectWhereField{
				Field: $1,
				Value: $3,
				Operate: $2,
			},
		}
	}
	// A OR B == !(!A AND !B)
	| SelectWhereList OR Expr CompareOperate Expr %prec OR
	{
		$4.Negate()
		field := &SelectWhereField{
			Field: $3,
			Value: $5,
			Operate: $4,
		}
		if len($$) == 1 {
			$$[0].Negate()
			$$ = append($$, field)
			$$ = []SelectWhere{
				&SelectWhereExpr{
					Negation: true,
					Cnf: $$,
				},
			}
		} else {
			$$ = []SelectWhere{
				&SelectWhereExpr{
					Negation: true,
					Cnf: []SelectWhere{
						&SelectWhereExpr{
							Negation: true,
							Cnf: $$,
						},
						field,
					},
				},
			}
		}
	}
	// A AND B
	| SelectWhereList AND Expr CompareOperate Expr %prec AND
	{
		$$ = append($$, &SelectWhereField{
			Field: $3,
			Value: $5,
			Operate: $4,
		})
	}
	// A OR (B...) == !(!A AND !(B...))
	| SelectWhereList OR '(' SelectWhereList ')' %prec OR
	{
		expr := &SelectWhereExpr{
			Negation: true,
			Cnf: $4,
		}
		if len($$) == 1 {
			$$[0].Negate()
			$$ = append($$, expr)
			$$ = []SelectWhere{
				&SelectWhereExpr{
					Negation: true,
					Cnf: $$,
				},
			}
		} else {
			$$ = []SelectWhere{
				&SelectWhereExpr{
					Negation: true,
					Cnf: []SelectWhere{
						&SelectWhereExpr{
							Negation: true,
							Cnf: $$,
						},
						expr,
					},
				},
			}
		}
	}
	// A AND (B...) == A AND B...
	| SelectWhereList AND '(' SelectWhereList ')' %prec AND
	{
		$$ = append($$, $4...)
	}

SelectOrder:
	{
		$$ = nil
	}
	| "ORDER" "BY" SelectOrderList
	{
		$$ = $3
	}

SelectOrderList:
	Expr Ascend
	{
		$$ = []*SelectOrder{
			&SelectOrder{
				Asc: $2,
				Field: $1,
			},
		}
	}
	| SelectOrderList ',' Expr Ascend
	{
		$$ = append($1, &SelectOrder{
			Asc: $4,
			Field: $3,
		})
	}

SelectLimit:
	{
		$$ = nil
	}
	| "LIMIT" VARIABLE
	{
		limit, err := strconv.Atoi($2)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		$$ = &SelectLimit{
			Limit: limit,
		}
	}
	| "LIMIT" VARIABLE ',' VARIABLE
	{
		limit, err := strconv.Atoi($2)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		offset, err := strconv.Atoi($4)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		$$ = &SelectLimit{
			Limit: limit,
			Offset: offset,
		}
	}
	| "LIMIT" VARIABLE "OFFSET" VARIABLE
	{
		limit, err := strconv.Atoi($2)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		offset, err := strconv.Atoi($4)
		if err != nil {
			yylex.Error(err.Error())
			goto ret1
		}
		$$ = &SelectLimit{
			Limit: limit,
			Offset: offset,
		}
	}

// 语法规则（插入表）
InsertStmt:
	"INSERT" InsertTable InsertField InsertValue ';'
	{
		$$ = &InsertStmt{
			Table: $2,
			Fields: $3,
			Values: $4,
		}
	}

InsertTable:
	Expr
	{
		$$ = $1
	}
	| "INTO" Expr
	{
		$$ = $2
	}

InsertField:
	'(' InsertFieldList ')'
	{
		$$ = $2
	}

InsertFieldList:
	{
		$$ = nil
	}
	| VaribleList

InsertValue:
	"VALUE" InsertValueList
	{
		$$ = $2
	}
	| "VALUES" InsertValueList
	{
		$$ = $2
	}

InsertValueList:
	'(' VaribleList ')'
	{
		$$ = [][]string{$2}
	}
	| InsertValueList ',' '(' VaribleList ')'
	{
		$$ = append($1, $4)
	}

%%