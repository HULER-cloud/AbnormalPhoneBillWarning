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
func parseFloat(value string) float64 {
	parsedFloat, _ := strconv.ParseFloat(value, 64)
	return parsedFloat
}

// 辅助函数：字符串转整数
func atoi(value string) int {
	parsedInt, _ := strconv.Atoi(value)
	return parsedInt
}
