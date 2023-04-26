package common

import "time"

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TimestampToTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
