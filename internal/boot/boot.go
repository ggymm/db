package boot

import (
	"io"
	"os"
	"path/filepath"

	"db/internal/opt"
	"db/pkg/utils"
)

const (
	Suffix    = ".boot"
	SuffixTmp = ".boot_tmp"
)

type Boot interface {
	Load() []byte
	Update(data []byte)
}

type boot struct {
	f    *os.File
	path string
}

func New(opt *opt.Option) Boot {
	_ = os.Remove(opt.Path + SuffixTmp)

	var (
		err error

		f    *os.File
		path = opt.GetPath(Suffix)
	)
	if opt.Open {
		// 读取文件
		f, err = os.OpenFile(path, os.O_RDWR, 0o666)
		if err != nil {
			panic(err)
		}
	} else {
		// 创建父文件夹
		dir := filepath.Dir(path)
		if !utils.IsExist(dir) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		// 创建文件
		f, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			panic(err)
		}
	}
	return &boot{
		f:    f,
		path: opt.GetPath(""), // 不带后缀的路径
	}
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
	tmpFile, err := os.OpenFile(b.path+SuffixTmp, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o666)
	if err != nil {
		panic(err)
	}

	_, err = tmpFile.Write(data)
	if err != nil {
		panic(err)
	}
	err = tmpFile.Sync()
	if err != nil {
		panic(err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	// 文件重命名（保证原子性）
	_ = b.f.Close()
	err = os.Rename(b.path+SuffixTmp, b.path+Suffix)
	if err != nil {
		panic(err)
	}

	// 重新打开文件
	b.f, err = os.OpenFile(b.path+Suffix, os.O_RDWR, 0o666)
	if err != nil {
		panic(err)
	}
}
