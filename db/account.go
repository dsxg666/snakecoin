package db

import (
	"path/filepath"
	"strings"
)

// AccountIsInDB 判断账户是否是合法账户，高级版
// 返回1：表示是并且不是自己
// 返回0：表示是并且就是自己
// 返回-1：不是
func AccountIsInDB(from, to string) int {
	if AccountIsExist(to) && strings.Compare(from, to) != 0 {
		return 1
	} else if strings.Compare(from, to) == 0 {
		return 0
	} else {
		return -1
	}
}

// AccountIsExist 判断账户是否是合法账户，低级版
func AccountIsExist(to string) bool {
	files, _ := filepath.Glob(KeystoreDataPath + "/*")
	for i := 0; i < len(files); i++ {
		if strings.Compare(files[i][14:], to) == 0 {
			return true
		}
	}
	return false
}
