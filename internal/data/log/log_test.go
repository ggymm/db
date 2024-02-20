package log

import (
	"path/filepath"
	"testing"

	"db/internal/opt"
	"db/pkg/utils"
)

func newOpt() *opt.Option {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/log")

	if !utils.IsEmpty(path) {
		return &opt.Option{
			Open: true,
			Path: path,
		}
	} else {
		return &opt.Option{
			Open: false,
			Path: path,
		}
	}
}

func TestNewLog(t *testing.T) {
	log := NewLog(newOpt())
	log.Log(utils.RandBytes(60))
}

func TestLogger_Next(t *testing.T) {
}
