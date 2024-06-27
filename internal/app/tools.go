package app

import (
	"strconv"
	"time"
)

// 辅助函数
// 辅助函数：解析时间
func parseTime(value string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, value)
	return parsedTime
}

// 辅助函数：解析浮点数
func parseFloat(value string) float32 {
	parsedFloat, _ := strconv.ParseFloat(value, 32)
	return float32(parsedFloat)
}

// 辅助函数：字符串转整数
func atoi(value string) uint {
	parsedInt, _ := strconv.Atoi(value)
	return uint(parsedInt)
}
