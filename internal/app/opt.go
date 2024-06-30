package app

import (
	"path/filepath"
)

type Option struct {
	Open   bool
	Name   string
	Path   string
	Memory int64

	RootId uint64 // 索引根节点的 id
}

func NewOption(path string) *Option {
	return new(Option)
}

func (o *Option) GetPath(suffix string) string {
	return filepath.Join(o.Path, o.Name+suffix)
}
