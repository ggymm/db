%{
package sql
%}

%union {
	str string
	strList []string
	boolean bool

	stmt Statement
	stmtList []Statement

	createStmt *CreateStmt
	tableDef *TableDef
	fieldDef *FieldDef
	indexDef *IndexDef
	tableOption *TableOption
	fieldType FieldType
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
%type <createStmt> CreateStmt
%type <tableDef> TableDef

%type <fieldDef> FieldDef

%type <indexDef> IndexDef
%type <indexDef> PrimaryDef
%type <tableOption> TableOption

%type <fieldType> FieldType
%type <boolean> AllowNull
%type <str> DefaultValue


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
CreateStmt:
	"CREATE" "TABLE" Expr '(' TableDef ')' TableOption ';'
	{
		$$ = &CreateStmt{
			TableName: $3,
			TableDef: $5,
			TableOption: $7,
		}
	}

TableDef:
	FieldDef
	{
		$$ = &TableDef{
			Field: []*FieldDef{$1},
			Index: []*IndexDef{},
		}
	}
	| IndexDef
	{
		$$ = &TableDef{
			Field: []*FieldDef{},
			Index: []*IndexDef{$1},
		}
	}
	| PrimaryDef
	{
		$$ = &TableDef{
			Field: []*FieldDef{},
			Index: []*IndexDef{},
			Primary: $1,
		}
	}
	| TableDef ',' FieldDef
	{
		$$.Field = append($$.Field, $3)
	}
	| TableDef ',' IndexDef
	{
		$$.Index = append($$.Index, $3)
	}
	| TableDef ',' PrimaryDef
	{
		if $$.Primary == nil {
			$$.Primary = $3
		} else {
			 __yyfmt__.Printf("重复定义主键 %v %v ", $$.Primary, $3)
			goto ret1
		}
	}

FieldDef:
	Expr FieldType AllowNull DefaultValue
	{
		$$ = &FieldDef{
			FieldName: $1,
			FieldType: $2,
			AllowNull: $3,
			DefaultValue: $4,
		}
	}

IndexDef:
	"INDEX" Expr '(' VaribleList ')'
	{
		$$ = &IndexDef{
			IndexName: $2,
			IndexField: $4,
		}
	}

PrimaryDef:
	"PRIMARY" "KEY" '(' VaribleList ')'
	{
		$$ = &IndexDef{
			Primary: true,
			IndexField: $4,
		}
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

TableOption:
	{
		$$ = nil
	}


%%