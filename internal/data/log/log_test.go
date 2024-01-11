package log

import (
	"path/filepath"
	"testing"

	"db/internal/ops"
	"db/pkg/utils"
)

func NewOps() *ops.Option {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/log")

	if !utils.IsEmpty(path) {
		return &ops.Option{
			Open: true,
			Path: path,
		}
	} else {
		return &ops.Option{
			Open: false,
			Path: path,
		}
	}
}

func TestNewLog(t *testing.T) {
	log := NewLog(NewOps())
	log.Log(utils.RandBytes(60))
}

func TestLogger_Next(t *testing.T) {

}
