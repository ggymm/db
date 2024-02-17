package sql

type StmtType int

const (
	_ StmtType = iota
	Create
	Select
	Insert
)

type FieldType int

const (
	_ FieldType = iota
	Int
	Int64
	Varchar
)

var typeMapping = map[string]FieldType{
	"INT":     Int,
	"INT64":   Int64,
	"VARCHAR": Varchar,
}

type CompareOp int

const (
	EQ CompareOp = iota // =
	NE                  // !=
	LT                  // <
	GT                  // >
	LE                  // <=
	GE                  // >=
)

func (o *CompareOp) Negate() {
	switch *o {
	case EQ:
		*o = NE
	case NE:
		*o = EQ
	case LT:
		*o = GE
	case GT:
		*o = LE
	case LE:
		*o = GT
	case GE:
		*o = LT
	}
}

type Statement interface {
	GetStmtType() StmtType
}

type CreateStmt struct {
	TableName   string
	Table       *CreateTable
	TableOption *CreateTableOption
}

func (*CreateStmt) GetStmtType() StmtType {
	return Create
}

type CreateTable struct {
	Field   []*CreateField
	Index   []*CreateIndex
	Primary *CreateIndex
}

type CreateField struct {
	FieldName    string
	FieldType    FieldType
	AllowNull    bool
	DefaultValue string
}

type CreateIndex struct {
	IndexName  string
	Primary    bool
	IndexField []string
}

type CreateTableOption struct {
}

type SelectStmt struct {
	Field []string
	Table string
	Where interface{}
	Order interface{}
	Limit interface{}
}

func (*SelectStmt) GetStmtType() StmtType {
	return Select
}

type InsertStmt struct {
}

func (*InsertStmt) GetStmtType() StmtType {
	return Insert
}
