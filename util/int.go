package util

import (
	"encoding/binary"
)

// Int64ToBytes 将 int64 转换成 []byte
func Int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
