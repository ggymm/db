package log

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"db/pkg/utils"
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
	Seed     = 12321
	Checksum = 4

	OffSize     = 0
	OffChecksum = OffSize + 4
	OffData     = OffChecksum + 4

	Suffix = ".log"
)

type Logger interface {
	Close()

	Log(data []byte)
	Truncate(size int64) error

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

	buf := make([]byte, Checksum)
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
	buf := make([]byte, Checksum)
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

func wrapLog(data []byte) []byte {
	buf := make([]byte, OffData+len(data))
	// 写入 size
	binary.LittleEndian.PutUint32(buf[OffSize:], uint32(len(data)))
	// 写入 checksum
	binary.LittleEndian.PutUint32(buf[OffChecksum:], calcChecksum(0, data))
	return buf
}

func calcChecksum(res uint32, data []byte) uint32 {
	for _, b := range data {
		res = res*Seed + uint32(b)
	}
	return res
}

func NewLog(filename string) Logger {
	l := new(logger)
	l.filename = filename + Suffix

	if utils.IsExist(l.filename) {
		open(l)
	} else {
		create(l)
	}

	// 校验日志文件
	err := l.check()
	if err != nil {
		panic(err)
	}
	return l
}

func (l *logger) next() ([]byte, bool, error) {
	// 如果要读取的数据超过 filesize，返回 false
	if l.pos+OffData >= l.filesize {
		return nil, false, nil
	}

	// 读取 size
	buf := make([]byte, OffChecksum-OffSize)
	_, err := l.file.ReadAt(buf, l.pos+OffSize)
	if err != nil {
		return nil, false, err
	}
	size := binary.LittleEndian.Uint32(buf)

	// 数据不完整
	if l.pos+OffData+int64(size) > l.filesize {
		return nil, false, nil
	}

	// 读取数据
	buf = make([]byte, OffData+size)
	_, err = l.file.ReadAt(buf, l.pos)
	if err != nil {
		return nil, false, err
	}

	// 校验 checksum 是否正确
	var (
		cs1 uint32
		cs2 uint32
	)
	cs1 = calcChecksum(0, buf[OffData:])
	binary.LittleEndian.PutUint32(buf[OffChecksum:], cs2)
	if cs1 != cs2 {
		return nil, false, nil
	}

	l.pos += int64(len(buf))
	return buf, true, nil
}

func (l *logger) check() error {
	return nil
}

func (l *logger) Close() {
	err := l.file.Close()
	if err != nil {
		panic(err)
	}
}

func (l *logger) Log(data []byte) {
	l.lock.Lock()
	defer l.lock.Unlock()

	log := wrapLog(data)
	_, err := l.file.Write(log)
	if err != nil {
		panic(err)
	}

	// 计算新的 checksum
	l.checksum = calcChecksum(l.checksum, data)
	buf := make([]byte, Checksum)
	binary.LittleEndian.PutUint32(buf, l.checksum)
	_, err = l.file.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}
	err = l.file.Sync()
	if err != nil {
		panic(err)
	}
}

func (l *logger) Truncate(size int64) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.file.Truncate(size)
}

func (l *logger) Next() ([]byte, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	data, next, err := l.next()
	if err != nil {
		panic(err)
	}

	if next == false {
		return nil, false
	}
	return data[OffData:], false
}

func (l *logger) Rewind() {
	l.pos = Checksum
}
