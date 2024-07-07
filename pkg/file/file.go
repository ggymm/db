package file

import (
	"bufio"
	"errors"
	"io"
	"os"
)

const Mode = 0o666

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
