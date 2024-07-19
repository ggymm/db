package sql

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

type Type int

const (
	_ Type = iota
	Begin
	Commit
	Rollback

	Create
	Select
	Insert
	Update
	Delete
)

type FieldType int

const (
	_ FieldType = iota
	Int32
	Int64
	Varchar
)

var typeMapping = map[string]FieldType{
	"INT32":   Int32,
	"INT64":   Int64,
	"VARCHAR": Varchar,
}

func (t FieldType) String() string {
	switch t {
	case Int32:
		return "INT32"
	case Int64:
		return "INT64"
	case Varchar:
		return "VARCHAR"
	}
	return ""
}

type CompareOperate int

const (
	EQ CompareOperate = iota // =
	NE                       // !=
	LT                       // <
	GT                       // >
	LE                       // <=
	GE                       // >=
)

func (o *CompareOperate) Negate() {
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
	StmtType() Type
	TableName() string
}

type BeginStmt struct {
	Level string
}

func (*BeginStmt) StmtType() Type {
	return Begin
}

func (*BeginStmt) TableName() string {
	return ""
}

type CommitStmt struct {
}

func (*CommitStmt) StmtType() Type {
	return Commit
}

func (*CommitStmt) TableName() string {
	return ""
}

type RollbackStmt struct {
}

func (*RollbackStmt) StmtType() Type {
	return Rollback
}

func (*RollbackStmt) TableName() string {
	return ""
}

type CreateStmt struct {
	Name   string
	Table  *CreateTable
	Option *CreateTableOption
}

func (*CreateStmt) StmtType() Type {
	return Create
}

func (s *CreateStmt) TableName() string {
	return s.Name
}

type CreateTable struct {
	Pk    *CreateIndex
	Field []*CreateField
	Index []*CreateIndex
}

type CreateField struct {
	Name     string
	Type     FieldType
	Default  string
	Nullable bool
}

type CreateIndex struct {
	Pk    bool
	Name  string
	Field string
}

type CreateTableOption struct{}

type InsertStmt struct {
	Table string
	Field []string
	Value []string
}

func (*InsertStmt) StmtType() Type {
	return Insert
}

func (s *InsertStmt) TableName() string {
	return s.Table
}

func (s *InsertStmt) FormatData() (map[string]string, error) {
	if len(s.Value) != len(s.Field) {
		return nil, fmt.Errorf("插入列数与值数不匹配")
	}
	row := make(map[string]string)
	for i, v := range s.Value {
		row[s.Field[i]] = v
	}
	return row, nil
}

type UpdateStmt struct {
	Table string
	Value map[string]string
	Where []SelectWhere
}

func (*UpdateStmt) StmtType() Type {
	return Update
}

func (s *UpdateStmt) TableName() string {
	return s.Table
}

type DeleteStmt struct {
	Table string
	Where []SelectWhere
}

func (*DeleteStmt) StmtType() Type {
	return Delete
}

func (s *DeleteStmt) TableName() string {
	return s.Table
}

type SelectStmt struct {
	Table string
	Field []*SelectField
	Where []SelectWhere
	Order []*SelectOrder
	Limit *SelectLimit
}

func (*SelectStmt) StmtType() Type {
	return Select
}

func (s *SelectStmt) TableName() string {
	return s.Table
}

type SelectField struct {
	Name  string
	Alias string
}

type SelectWhere interface {
	// Negate 取反
	Negate()

	// Match 判断是否符合条件
	Match(row map[string]any) bool
}

type SelectWhereExpr struct {
	Negation bool
	Cnf      []SelectWhere
}

func (w *SelectWhereExpr) Negate() {
	w.Negation = !w.Negation
}

func (w *SelectWhereExpr) Match(row map[string]any) bool {
	filter := true
	for _, be := range w.Cnf {
		filter = filter && be.Match(row)
	}

	if w.Negation {
		return !filter
	}
	return filter
}

type SelectWhereField struct {
	Pos     int
	Field   string
	Value   string
	Operate CompareOperate
}

func (w *SelectWhereField) Negate() {
	w.Operate.Negate()
}

func (w *SelectWhereField) Match(row map[string]any) bool {
	val, ok := row[w.Field]
	if !ok {
		return false
	}

	var r int
	switch v := val.(type) {
	case uint32:
		dst, err := strconv.ParseUint(w.Value, 10, 32)
		if err != nil {
			return false
		}
		r = cmp.Compare(v, uint32(dst))
	case uint64:
		dst, err := strconv.ParseUint(w.Value, 10, 64)
		if err != nil {
			return false
		}
		r = cmp.Compare(v, dst)
	case string:
		r = strings.Compare(v, w.Value)
	default:
		return false
	}
	switch w.Operate {
	case EQ:
		return r == 0
	case NE:
		return r != 0
	case LT:
		return r < 0
	case GT:
		return r > 0
	case LE:
		return r <= 0
	case GE:
		return r >= 0
	}
	return false
}

type SelectOrder struct {
	Asc   bool
	Field string
}

type SelectLimit struct {
	Limit  int
	Offset int
}
