package util

import (
	"time"
)

func CurrentTimeFormant() string {
	str := time.Now().Format("2006-01-02 15:04:05")
	return str
}

func TimestampFormat(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}
