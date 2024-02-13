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

type Statement interface {
	GetStmtType() StmtType
}

type CreateStmt struct {
	TableName   string
	TableDef    *TableDef
	TableOption *TableOption
}

func (*CreateStmt) GetStmtType() StmtType {
	return Create
}

type TableDef struct {
	TableId   int
	PKAsRowId bool
	RowId     int
	Field     []*FieldDef
	Index     []*IndexDef
	Primary   *IndexDef
}

type FieldDef struct {
	FieldName    string
	FieldType    FieldType
	AllowNull    bool
	DefaultValue string
}

type IndexDef struct {
	IndexName  string
	Primary    bool
	IndexField []string
}

type TableOption struct {
}
