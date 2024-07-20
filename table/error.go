package table

import (
	"errors"
	"fmt"
)

var (
	ErrNoSuchTable       = NewError("no such table")
	ErrNoPrimaryKey      = NewError("no primary key")
	ErrMustHaveCondition = NewError("must have condition")

	ErrInsertNotMatch = errors.New("mismatch between number of fields and values")
)

const (
	ErrNotAllowNull = "field %s is not allowed to be null"
)

func NewError(msg string, args ...any) error {
	if len(args) == 0 {
		return errors.New(msg)
	}
	return fmt.Errorf(msg, args...)
}
