package rlp

import (
	"reflect"
	"unsafe"
)

// byteArrayBytes returns a slice of the byte array v.
func byteArrayBytes(v reflect.Value, length int) []byte {
	var s []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	hdr.Data = v.UnsafeAddr()
	hdr.Cap = length
	hdr.Len = length
	return s
}
