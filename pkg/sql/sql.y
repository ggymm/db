%{
package sql
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
}

%token <str>
	CREATE "CREATE"
	TABLE "TABLE"
	KEY "KEY"
	NOT	"NOT"
	NULL "NULL"
	INDEX "INDEX"
	DEFAULT "DEFAULT"
	PRIMARY "PRIMARY"

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
%type <compareOperate> CompareOperate

%type <selectStmt> SelectStmt


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
	"CREATE" "TABLE" Expr '(' createTable ')' createTableOption ';'
	{
		$$ = &CreateStmt{
			TableName: $3,
			Table: $5,
			TableOption: $7,
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
			Field: []*CreateField{},
			Index: []*CreateIndex{},
			Primary: $1,
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
		if $$.Primary == nil {
			$$.Primary = $3
		} else {
			 __yyfmt__.Printf("重复定义主键 %v %v ", $$.Primary, $3)
			goto ret1
		}
	}

CreateField:
	Expr FieldType AllowNull DefaultValue
	{
		$$ = &CreateField{
			FieldName: $1,
			FieldType: $2,
			AllowNull: $3,
			DefaultValue: $4,
		}
	}

CreateIndex:
	"INDEX" Expr '(' VaribleList ')'
	{
		$$ = &CreateIndex{
			IndexName: $2,
			IndexField: $4,
		}
	}

CreatePrimary:
	"PRIMARY" "KEY" '(' VaribleList ')'
	{
		$$ = &CreateIndex{
			Primary: true,
			IndexField: $4,
		}
	}

CreateTableOption:
	{
		$$ = nil
	}

// 语法规则（查询表）
CompareOp:
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
    "SELECT" SelectFieldList SelectStmtLimit ';'
    {
        $$ = &SelectStmt{
        	Filed: $2,
        	Limit: $3,
        }
    }
    |  "SELECT" SelectFieldList SelectStmtFrom SelectStmtWhere SelectStmtOrder SelectStmtLimit ';'
    {
        $$ = &SelectStmt{
        	Filed: $2,
        	Table: $3,
        	Where: $4,
        	Order: $5,
        	Limit: $6,
        }
    }

%%