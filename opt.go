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

type Option struct {
	Open   bool
	Path   string
	Memory int64

	RootId uint64 // 索引根节点的 id
}

func NewOption(p ...string) *Option {
	var (
		open bool
		path string
	)

	if len(p) <= 0 {
		path = ""
	} else {
		for _, v := range p {
			path = filepath.Join(path, v)
		}
	}
	if path == "" {
		open = false
		path = "data"
	} else {
		// 判断目录是否存在
		_, err := os.Stat(path)
		if err == nil {
			// 判断目录是否为空
			files, err1 := os.ReadDir(path)
			if err1 != nil {
				open = false
			} else {
				open = len(files) > 0
			}
		} else {
			open = false
		}
	}
	return &Option{
		Open: open,
		Path: path,
	}
}

func (o *Option) GetPath(name string) string {
	return filepath.Join(o.Path, name)
}
