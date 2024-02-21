package opt

import "path/filepath"

type Option struct {
	Open   bool
	Name   string
	Path   string
	Memory int64

	RootId uint64 // 索引根节点的 id
}

func (o *Option) GetPath(suffix string) string {
	return filepath.Join(o.Path, o.Name+suffix)
}
