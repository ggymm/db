package db

import (
	"testing"
)

func Test_NewOption(t *testing.T) {
	opt := NewOption("C:\\Users\\19679\\code\\basic\\db\\temp\\table\\a")
	t.Logf("%+v", opt)
}
