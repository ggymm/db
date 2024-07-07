package log

import (
	"path/filepath"
	"testing"

	"db"
	"db/pkg/file"
)

func newOpt() *db.Option {
	base := db.RunPath()
	path := filepath.Join(base, "temp/log")

	if !file.IsEmpty(path) {
		return &db.Option{
			Open: true,
			Path: path,
		}
	} else {
		return &db.Option{
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
