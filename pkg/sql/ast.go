package sql

import (
	"fmt"
	"strconv"
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
	row := make(map[string]string)
	for i, v := range s.Value {
		if len(v) != len(s.Field) {
			return nil, fmt.Errorf("插入列数与值数不匹配")
		}
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
	Negate()
	Prepare(fieldMapping map[string]int) error
	Filter(row *[]string) bool
}

type SelectWhereExpr struct {
	Negation bool
	Cnf      []SelectWhere
}

func (w *SelectWhereExpr) Negate() {
	w.Negation = !w.Negation
}

func (w *SelectWhereExpr) Prepare(fieldMapping map[string]int) error {
	for _, be := range w.Cnf {
		err := be.Prepare(fieldMapping)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *SelectWhereExpr) Filter(row *[]string) bool {
	filter := true
	for _, be := range w.Cnf {
		filter = filter && be.Filter(row)
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

func (w *SelectWhereField) Prepare(fieldMapping map[string]int) error {
	pos, exist := fieldMapping[w.Field]
	if !exist {
		_, err := strconv.Atoi(w.Field)
		if err != nil {
			return fmt.Errorf("查询列 %s 不存在", w.Field)
		}
		w.Pos = -1
		return nil
	}
	w.Pos = pos
	return nil
}

func (w *SelectWhereField) Filter(row *[]string) bool {
	var val string
	if w.Pos == -1 {
		val = w.Field
	} else {
		val = (*row)[w.Pos]
	}

	switch w.Operate {
	case EQ:
		return w.Value == val
	case NE:
		return w.Value != val
	case LT:
		return w.Value < val
	case GT:
		return w.Value > val
	case LE:
		return w.Value <= val
	case GE:
		return w.Value >= val
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
