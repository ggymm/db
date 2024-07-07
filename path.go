package db

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func RunPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	base := filepath.Base(exe)
	if !strings.HasPrefix(base, "___") {
		return filepath.Dir(exe)
	} else {
		var path string
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			path = filepath.Dir(filename)
		}
		return path
	}
}
