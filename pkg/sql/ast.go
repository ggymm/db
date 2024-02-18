package sql

import (
	"fmt"
	"strconv"
)

type StmtType int

const (
	_ StmtType = iota
	Create
	Select
	Insert
	Update
	Delete
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
	GetStmtType() StmtType
}

type CreateStmt struct {
	Name   string
	Table  *CreateTable
	Option *CreateTableOption
}

func (*CreateStmt) GetStmtType() StmtType {
	return Create
}

type CreateTable struct {
	Pk    *CreateIndex
	Field []*CreateField
	Index []*CreateIndex
}

type CreateField struct {
	Name         string
	Type         FieldType
	AllowNull    bool
	DefaultValue string
}

type CreateIndex struct {
	Pk    bool
	Name  string
	Field []string
}

type CreateTableOption struct {
}

type SelectStmt struct {
	Field []*SelectField
	From  *SelectFrom
	Where []SelectWhere
	Order []*SelectOrder
	Limit *SelectLimit
}

func (*SelectStmt) GetStmtType() StmtType {
	return Select
}

type SelectFrom struct {
	Name string
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

type InsertStmt struct {
	Table  string
	Fields []string
	Values [][]string
}

func (*InsertStmt) GetStmtType() StmtType {
	return Insert
}

type UpdateStmt struct {
	Table string
	Value map[string]string
	Where []SelectWhere
}

func (*UpdateStmt) GetStmtType() StmtType {
	return Update
}

type DeleteStmt struct {
	Table string
	Where []SelectWhere
}

func (*DeleteStmt) GetStmtType() StmtType {
	return Delete
}