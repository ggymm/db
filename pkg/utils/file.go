package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func IsEmpty(path string) bool {
	// 首先判断目录是否存在
	if !IsExist(path) {
		return true
	}

	// 判断目录是否为空
	files, err := os.ReadDir(path)
	if err != nil {
		return true
	}
	return len(files) == 0
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	buf := bufio.NewReader(f)
	lines := make([]string, 0)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		lines = append(lines, string(l))
	}
	return lines, nil
}
