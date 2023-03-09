package util

import (
	"encoding/binary"
	"unsafe"
)

// BytesToInt64 将 []byte 转换成 int64
func BytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

// BytesToString 将 []byte 转换成 string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
