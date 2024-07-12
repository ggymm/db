package boot

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ggymm/db"
	"github.com/ggymm/db/pkg/bin"
	"github.com/ggymm/db/pkg/file"
)

const (
	name = "BOOT"
	temp = "BOOT_TMP"
)

type Boot interface {
	Load() []byte
	Update(data []byte)
}

type boot struct {
	f    *os.File
	path string
}

func New(opt *db.Option) Boot {
	var (
		err error

		b  = new(boot)
		f  = filepath.Join(opt.Path, name)
		ft = filepath.Join(opt.Path, temp)
	)

	_ = os.Remove(ft)
	b.path = opt.Path
	if opt.Open {
		// 读取文件
		b.f, err = os.OpenFile(f, os.O_RDWR, file.Mode)
		if err != nil {
			panic(err)
		}
	} else {
		// 创建父文件夹
		if !file.IsExist(b.path) {
			err = os.MkdirAll(b.path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		// 创建文件
		b.f, err = os.OpenFile(f, os.O_RDWR|os.O_TRUNC|os.O_CREATE, file.Mode)
		if err != nil {
			panic(err)
		}

		// 初始化文件
		b.Update(bin.Uint64Raw(0))
	}
	return b
}

func (b *boot) Load() []byte {
	var (
		err error
		buf []byte
	)

	// 重置文件指针
	_, err = b.f.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	// 读取文件内容
	buf, err = io.ReadAll(b.f)
	if err != nil {
		panic(err)
	}
	return buf
}

func (b *boot) Update(data []byte) {
	f := filepath.Join(b.path, name)
	ft := filepath.Join(b.path, temp)
	tmp, err := os.OpenFile(ft, os.O_RDWR|os.O_TRUNC|os.O_CREATE, file.Mode)
	if err != nil {
		panic(err)
	}

	_, err = tmp.Write(data)
	if err != nil {
		panic(err)
	}
	err = tmp.Sync()
	if err != nil {
		panic(err)
	}
	err = tmp.Close()
	if err != nil {
		panic(err)
	}

	// 文件重命名（保证原子性）
	_ = b.f.Close()
	err = os.Rename(ft, f)
	if err != nil {
		panic(err)
	}

	// 重新打开文件
	b.f, err = os.OpenFile(f, os.O_RDWR, file.Mode)
	if err != nil {
		panic(err)
	}
}
