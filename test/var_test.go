package test

import (
	"fmt"
	"testing"
)

func Test_Fmt(t *testing.T) {
	s := make([]any, 0)
	for i := 0; i < 13; i++ {
		s = append(s, fmt.Sprintf("%d", 1))
	}
	t.Log(fmt.Sprintf(InsertSQL, s...))
}
