package boot

import (
	"testing"

	"db"
)

func TestNew(t *testing.T) {
	opt := db.NewOption(db.RunPath(), "temp/boot")

	b := New(opt)
	t.Logf("%+v", b)
}

func TestBoot_Handle(t *testing.T) {
	opt := db.NewOption(db.RunPath(), "temp/boot")

	b := New(opt)
	t.Logf("%+v", b)

	data := []byte("hello")
	b.Update(data)

	data = b.Load()
	t.Logf("%s", data)
}
