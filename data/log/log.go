package log

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"db"
	"db/pkg/bin"
	"db/pkg/file"
)

// 日志文件的读写
//
// 日志文件结构如下：
// +----------------+----------------+----------------+----------------+
// |    checksum    |      log1      |      logN      |    bad tail    |
// +----------------+----------------+----------------+----------------+
// 日志文件的前 4 byte 为 checksum，后面是日志数据
//
// 每条日志的结构如下：
// +----------------+----------------+----------------+
// |     size       |    checksum    |      data      |
// +----------------+----------------+----------------+
// |    4 byte      |     4 byte     |      size      |
// +----------------+----------------+----------------+
//
// 每次插入数据，更新日志文件的 checksum

var ErrBadLogFile = errors.New("bad log file")

const (
	seed     = 12321
	checkLen = 4

	offSize  = 0
	offCheck = offSize + 4
	offData  = offCheck + 4

	suffix = ".log"
)

type Log interface {
	Close()

	Log(data []byte)
	Next() ([]byte, bool)
	Rewind()
}

type logger struct {
	mu sync.Mutex

	pos      int64    // 迭代器的指针位置
	file     *os.File // 文件句柄
	size     int64    // 文件大小
	checksum uint32

	filepath string
}

func open(l *logger) {
	// 打开文件
	f, err := os.OpenFile(l.filepath, os.O_RDWR, file.Mode)
	if err != nil {
		panic(err)
	}

	// 读取 size 和 checksum
	stat, _ := f.Stat()
	size := stat.Size()
	if size < 4 {
		panic(ErrBadLogFile)
	}

	buf := make([]byte, checkLen)
	_, err = f.ReadAt(buf, 0)
	if err != nil {
		panic(err)
	}
	checksum := readUint32(buf)

	// 字段信息
	l.file = f
	l.size = size
	l.checksum = checksum
}

func create(l *logger) {
	p := l.filepath

	// 创建父文件夹
	dir := filepath.Dir(p)
	if !file.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 创建文件
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode)
	if err != nil {
		panic(err)
	}

	// 写入 checksum
	updateChecksum(f, 0)

	// 字段信息
	l.file = f
	l.checksum = 0
}

func readUint32(buf []byte) uint32 {
	return bin.Uint32(buf)
}

func writeUint32(buf []byte, uint32 uint32) {
	bin.PutUint32(buf, uint32)
}

func calcChecksum(res uint32, data []byte) uint32 {
	for _, b := range data {
		res = res*seed + uint32(b)
	}
	return res
}

func updateChecksum(file *os.File, checksum uint32) {
	buf := make([]byte, checkLen)
	writeUint32(buf, checksum)
	_, err := file.WriteAt(buf, 0)
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func NewLog(opt *db.Option) Log {
	l := new(logger)
	l.filepath = filepath.Join(opt.GetPath(suffix))

	if opt.Open {
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

// next 通过 pos 读取日志
// 返回日志数据，并且移动 pos 到下一条日志位置
// 注意，返回的日志数据包含 size 和 checksum
//
// 这里分别使用 log 和 data 做区分
// log 代表完整日志数据，即包含 size 和 checksum
// data 只代表日志数据，即不包含 size 和 checksum
func (l *logger) next() ([]byte, bool, error) {
	// 如果要读取的数据超过 size，返回 false
	if l.pos+offData >= l.size {
		return nil, false, nil
	}

	// 读取 size
	buf := make([]byte, offCheck-offSize)
	_, err := l.file.ReadAt(buf, l.pos+offSize)
	if err != nil {
		return nil, false, err
	}
	size := readUint32(buf)

	// 数据不完整
	if l.pos+offData+int64(size) > l.size {
		return nil, false, nil
	}

	// 读取日志数据
	log := make([]byte, offData+size)
	_, err = l.file.ReadAt(log, l.pos)
	if err != nil {
		return nil, false, err
	}

	// 校验 checksum 是否正确
	checksum1 := calcChecksum(0, log[offData:])
	checksum2 := readUint32(log[offCheck:offData])
	if checksum1 != checksum2 {
		return nil, false, nil
	}

	l.pos += int64(len(log))
	return log, true, nil
}

func (l *logger) check() error {
	l.Rewind()
	var checksum uint32
	for {
		log, next, err := l.next()
		if err != nil {
			return err
		}
		if next == false {
			break
		}
		checksum = calcChecksum(checksum, log)
	}
	if checksum != l.checksum {
		return ErrBadLogFile
	}
	err := l.file.Truncate(l.pos)
	if err != nil {
		return err
	}
	_, err = l.file.Seek(l.pos, 0)
	if err != nil {
		return err
	}
	l.Rewind()
	return nil
}

func (l *logger) Close() {
	err := l.file.Close()
	if err != nil {
		panic(err)
	}
}

// Log 写入日志
// 首先包装日志数据，然后写入文件
// 最后更新 checksum，并且同步文件内容到磁盘
func (l *logger) Log(data []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 包装日志数据
	log := make([]byte, offData+len(data))
	writeUint32(log[offSize:], uint32(len(data)))
	writeUint32(log[offCheck:], calcChecksum(0, data))
	copy(log[offData:], data)

	// 写入文件
	// 在写入 checksum 时，同步文件内容到磁盘
	_, err := l.file.Write(log)
	if err != nil {
		panic(err)
	}

	// 计算并写入新的 checksum
	l.checksum = calcChecksum(l.checksum, log)
	updateChecksum(l.file, l.checksum)
}

// Next 读取日志
// 使用迭代器模式，每次读取一条日志
// 注意，返回的日志数据不包含 size 和 checksum
func (l *logger) Next() ([]byte, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	log, next, err := l.next()
	if err != nil {
		panic(err)
	}

	if next == false {
		return nil, false
	}
	return log[offData:], false
}

func (l *logger) Rewind() {
	l.pos = checkLen
}
