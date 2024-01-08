package log

import (
	"db/pkg/utils"
	"path/filepath"
	"testing"
)

func TestNewLog(t *testing.T) {
	base := utils.RunPath()
	filename := filepath.Join(base, "temp/dataLog/test", "log")

	log := NewLog(filename)
	log.Log([]byte("test"))
}

func TestLogger_Next(t *testing.T) {

}
