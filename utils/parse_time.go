package utils

import (
	"log"
	"time"
)

func ParseQueryTime(timeStr string) time.Time {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		log.Println(err)
		return time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	}

	// 创建 time.Date
	finalTime := time.Date(1, 1, 1, parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	return finalTime
}
