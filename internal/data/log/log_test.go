package log

import (
	"path/filepath"
	"testing"

	"db/internal/app"

	"db/pkg/file"
)

func newOpt() *app.Option {
	base := app.RunPath()
	path := filepath.Join(base, "temp/log")

	if !file.IsEmpty(path) {
		return &app.Option{
			Open: true,
			Path: path,
		}
	} else {
		return &app.Option{
			Open: false,
			Path: path,
		}
	}
}

func TestNewLog(t *testing.T) {
	log := NewLog(newOpt())
	log.Log([]byte("hello world"))
}

func TestLogger_Next(t *testing.T) {
}
