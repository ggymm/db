package log

import (
	"db/internal/ops"
	"path/filepath"
	"testing"

	"db/pkg/utils"
)

func NewOps() *ops.Option {
	base := utils.RunPath()
	path := filepath.Join(base, "temp/dataLog/test")

	if utils.IsExist(path) {
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
	log.Log([]byte("test"))
}

func TestLogger_Next(t *testing.T) {

}
