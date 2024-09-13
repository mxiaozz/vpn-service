package util

import (
	"fmt"
	"time"
)

func DiffDays(date1, date2 time.Time) int {
	diff := int(date2.Sub(date1).Milliseconds() / (1000 * 3600 * 24))
	if diff < 0 {
		return -diff
	}
	return diff
}

// 时间差（天/小时/分钟）
func TimeDistance(endTime, startTime time.Time) string {
	var nd int64 = 1000 * 24 * 60 * 60
	var nh int64 = 1000 * 60 * 60
	var nm int64 = 1000 * 60

	// 获得两个时间的毫秒时间差异
	diff := endTime.Sub(startTime).Milliseconds()
	// 计算差多少天
	day := diff / nd
	// 计算差多少小时
	hour := diff % nd / nh
	// 计算差多少分钟
	min := diff % nd % nh / nm

	return fmt.Sprintf("%d天%d时%d分", day, hour, min)
}
