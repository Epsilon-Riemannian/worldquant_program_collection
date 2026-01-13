package small_program

import (
	"fmt"
	"time"
)

// --------------------------------- worldquantbrain 同前端页面 alpha 提交时间规则 ------------------------------
// ConvertToUTCPlus5 通用时间转换函数，加5小时后转换为UTC
// 输入格式支持：
// 1. "2025-12-16" → 转换为 "2025-12-16T05:00:00Z"
// 2. "2025-12-16 12:02:13" → 转换为 "2025-12-16T17:02:13Z"
func ConvertToUTCPlus5(input string) (time.Time, error) {
	// 直接解析为完整时间格式
	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析时间格式: %s", input)
	}

	// UTC-5时区
	estLoc := time.FixedZone("UTC-5", -5*60*60)
	sourceTime := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), 0,
		estLoc,
	)

	return sourceTime.UTC(), nil
}
