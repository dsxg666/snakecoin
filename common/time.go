package common

import "time"

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TimestampToTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
