package util

import (
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// StringToBytes 将 string 转换成 []byte
func StringToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{str, len(str)},
	))
}

// StringIsDigit 判断 str 是否为字符串形式的数字
func StringIsDigit(str string) bool {
	var t int
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) && x != int32(46) {
			return false
		}
		if x == int32(46) {
			t++
		}
	}
	if t >= 2 {
		return false
	}
	return true
}

func StringIs0ToN(str string, num int) bool {
	for i := 0; i < num; i++ {
		if strings.Compare(str, strconv.FormatInt(int64(i), 10)) == 0 {
			return true
		}
	}
	return false
}

func StringIs1ToN(str string, num int) bool {
	for i := 0; i < num; i++ {
		if strings.Compare(str, strconv.FormatInt(int64(i+1), 10)) == 0 {
			return true
		}
	}
	return false
}
