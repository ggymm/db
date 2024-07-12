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

	beginStmt *BeginStmt
	commitStmt *CommitStmt
	rollbackStmt *RollbackStmt

	createStmt *CreateStmt
	createTable *CreateTable
	createField *CreateField
	createIndex *CreateIndex
	createTableOption *CreateTableOption

	insertStmt *InsertStmt

    updateStmt *UpdateStmt
    updateValue map[string]string
    
    deleteStmt *DeleteStmt

	selectStmt *SelectStmt
	selectFieldList []*SelectField
	selectWhereList []SelectWhere
	selectOrderList []*SelectOrder
	selectLimit *SelectLimit
}

%token <str>
	// 关键字（事务）
	BEGIN "BEGIN"
	COMMIT "COMMIT"
	ROLLBACK "ROLLBACK"
	// 关键字（创建表）
	CREATE "CREATE"
	TABLE "TABLE"
	KEY "KEY"
	NOT	"NOT"
	NULL "NULL"
	INDEX "INDEX"
	DEFAULT "DEFAULT"
	PRIMARY "PRIMARY"
	// 关键字（插入数据）
	INSERT "INSERT"
	INTO "INTO"
	VALUE "VALUE"
	// 关键字（更新数据）
	UPDATE "UPDATE"
	SET "SET"
	// 关键字（删除数据）
	DELETE "DELETE"
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
%type <str> DefaultVal
%type <boolean> AllowNull
%type <fieldType> FieldType

%type <beginStmt> BeginStmt
%type <commitStmt> CommitStmt
%type <rollbackStmt> RollbackStmt

%type <createStmt> CreateStmt
%type <createTable> CreateTable
%type <createField> CreateField
%type <createIndex> CreateIndex
%type <createIndex> CreatePrimary
%type <createTableOption> CreateTableOption

// 语法定义（插入数据）
%type <insertStmt> InsertStmt
%type <strList> InsertField InsertFieldList
%type <strList> InsertValue InsertValueList

// 语法定义（更新数据）
%type <updateStmt> UpdateStmt
%type <updateValue> UpdateValue

// 语法定义（删除数据）
%type <deleteStmt> DeleteStmt

// 语法定义（查询表）
%type <boolean> Ascend
%type <compareOperate> CompareOperate

%type <selectStmt> SelectStmt
%type <selectFieldList> SelectFieldList
%type <selectWhereList> SelectWhere SelectWhereList
%type <selectOrderList> SelectOrder SelectOrderList
%type <selectLimit> SelectLimit

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
	BeginStmt
	{
		$$ = Statement($1)
	}
	| CommitStmt
	{
		$$ = Statement($1)
	}
	| RollbackStmt
	{
		$$ = Statement($1)
	}
	| CreateStmt
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
	| UpdateStmt
	{
		$$ = Statement($1)
	}
	| DeleteStmt
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

DefaultVal:
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

BeginStmt:
    "BEGIN" Expr ';'
    {
        $$ = &BeginStmt{ $2 }
    }

CommitStmt:
    "COMMIT" ';'
    {
        $$ = &CommitStmt{}
    }

RollbackStmt:
    "ROLLBACK" ';'
    {
        $$ = &RollbackStmt{}
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
	Expr FieldType AllowNull DefaultVal
	{
		$$ = &CreateField{
			Name: $1,
			Type: $2,
			AllowNull: $3,
			DefaultVal: $4,
		}
	}

CreateIndex:
	"INDEX" Expr '(' Expr ')'
	{
		$$ = &CreateIndex{
			Name: $2,
			Field: $4,
		}
	}

CreatePrimary:
	"PRIMARY" "KEY" '(' Expr ')'
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

// 语法规则（插入数据）
InsertStmt:
	"INSERT" "INTO" Expr InsertField InsertValue ';'
	{
		$$ = &InsertStmt{
			Table: $3,
			Field: $4,
			Value: $5,
		}
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
	"VALUE" '(' InsertValueList ')'
	{
		$$ = $3
	}

InsertValueList:
	{
		$$ = nil
	}
	| VaribleList

// 语法规则（更新数据）
UpdateStmt:
	"UPDATE" Expr "SET" UpdateValue SelectWhere ';'
	{
		$$ = &UpdateStmt{
			Table: $2,
			Value: $4,
			Where: $5,
		}
	}

UpdateValue:
	Expr '=' Expr
	{
		$$ = map[string]string{
			$1: $3,
		}
	}
	| UpdateValue ',' Expr '=' Expr
	{
		$$[$3] = $5
	}

// 语法规则（删除数据）
DeleteStmt:
	"DELETE" "FROM" Expr SelectWhere ';'
	{
		$$ = &DeleteStmt{
			Table: $3,
			Where: $4,
		}
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
    |  "SELECT" SelectFieldList "FROM" Expr SelectWhere SelectOrder SelectLimit ';'
    {
        $$ = &SelectStmt{
        	Table: $4,
        	Field: $2,
        	Where: $5,
        	Order: $6,
        	Limit: $7,
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
%%