package rlp

import (
	"reflect"
	"unsafe"
)

// byteArrayBytes 返回字节数组v的一个切片。
func byteArrayBytes(v reflect.Value, length int) []byte {
	var s []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	hdr.Data = v.UnsafeAddr()
	hdr.Cap = length
	hdr.Len = length
	return s
}
