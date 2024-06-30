package boot

import (
	"path/filepath"
	"testing"

	"db/internal/app"
)

func TestNew(t *testing.T) {
	base := app.RunPath()
	path := filepath.Join(base, "temp/boot")

	opt := &app.Option{
		Open: false,
		Name: "boot",
		Path: path,
	}

	b := New(opt)
	t.Logf("%+v", b)
}

func TestBoot_Handle(t *testing.T) {
	base := app.RunPath()
	path := filepath.Join(base, "temp/boot")

	opt := &app.Option{
		Open: false,
		Name: "boot",
		Path: path,
	}

	b := New(opt)
	t.Logf("%+v", b)

	data := []byte("hello")
	b.Update(data)

	data = b.Load()
	t.Logf("%s", data)
}
