package file

import (
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
