package data

import (
	"db/pkg/utils"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// 日志文件的读写
// 日志文件结构如下：
// [Checksum] [Log1] [Log2] [Log3] ... [BadTail]
//
// 每条日志的结构如下：
// [Size] uint32 （4byte）
// [Checksum] uint32 （4byte）
// [Data] size
//
// 每次插入数据，更新日志文件的 Checksum

var (
	ErrBadLogFile = errors.New("bad log file")
)

const (
	Seed = 12321

	Size     = 0
	Checksum = Size + 4
	Data     = Checksum + 4

	Suffix = ".log"
)

type Logger interface {
	Close()

	Log(data []byte)
	Truncate(pos int64) error

	Next() ([]byte, bool)
	Rewind()
}

type logger struct {
	lock sync.Mutex

	pos      int64
	file     *os.File
	filesize int64
	checksum uint32

	filename string
}

func open(l *logger) {
	// 打开文件
	file, err := os.OpenFile(l.filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// 读取 filesize 和 checksum
	stat, _ := file.Stat()
	filesize := stat.Size()
	if filesize < 4 {
		panic(ErrBadLogFile)
	}

	buf := make([]byte, 4)
	_, err = file.ReadAt(buf, 0)
	if err != nil {
		panic(err)
	}
	checksum := binary.LittleEndian.Uint32(buf)

	// 字段信息
	l.file = file
	l.filesize = filesize
	l.checksum = checksum
}

func create(l *logger) {
	filename := l.filename

	// 创建父文件夹
	dir := filepath.Dir(filename)
	if !utils.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	// 写入 checksum
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 0)
	_, err = file.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}

	// 字段信息
	l.file = file
	l.checksum = 0
}

func NewLog(filename string) Logger {
	l := new(logger)
	l.filename = filename + Suffix

	if utils.IsExist(l.filename) {
		open(l)
	} else {
		create(l)
	}
	return l
}

func (l *logger) init() {
}

func (l *logger) Close() {
}

func (l *logger) Log(data []byte) {
}

func (l *logger) Truncate(pos int64) error {
	return nil
}

func (l *logger) Next() ([]byte, bool) {
	return nil, false
}

func (l *logger) Rewind() {
}
