package log

import (
	"testing"

	"github.com/ggymm/db"
)

func TestNewLog(t *testing.T) {
	log := NewLog(db.NewOption(db.RunPath(), "temp/log"))
	log.Log([]byte("hello world"))
}

func TestLogger_Next(t *testing.T) {
}
