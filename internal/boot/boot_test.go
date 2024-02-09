package boot

import (
	"db/internal/opt"
	"db/pkg/utils"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/boot")

	cfg := &opt.Option{
		Open: false,
		Path: path,
		Name: "boot",
	}

	b := New(cfg)
	t.Logf("%+v", b)
}

func TestBoot_Handle(t *testing.T) {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/boot")

	cfg := &opt.Option{
		Open: false,
		Path: path,
		Name: "boot",
	}

	b := New(cfg)
	t.Logf("%+v", b)

	data := []byte("hello")
	b.Update(data)

	data = b.Load()
	t.Logf("%s", data)
}
