package utils

import (
	"errors"
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
