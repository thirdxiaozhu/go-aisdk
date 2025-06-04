/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 17:42:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-03 21:11:18
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"fmt"
	"time"
)

// FormatDuration 格式化时间，根据时长选择合适的单位
func FormatDuration(duration time.Duration) (timeStr string) {
	absDuration := duration.Abs()
	switch {
	case absDuration < time.Second:
		timeStr = fmt.Sprintf("%.0fms", float64(absDuration.Nanoseconds())/1e6) // 毫秒
	case absDuration < time.Minute:
		timeStr = fmt.Sprintf("%gs", absDuration.Seconds()) // 秒
	case absDuration < time.Hour:
		timeStr = fmt.Sprintf("%gm", absDuration.Minutes()) // 分钟
	case absDuration < time.Hour*24:
		timeStr = fmt.Sprintf("%gh", absDuration.Hours()) // 小时
	case absDuration < time.Hour*24*30:
		timeStr = fmt.Sprintf("%gD", absDuration.Hours()/24) // 天
	case absDuration < time.Hour*24*30*12:
		timeStr = fmt.Sprintf("%gM", absDuration.Hours()/24/30) // 月
	default:
		timeStr = fmt.Sprintf("%gY", absDuration.Hours()/24/30/12) // 年
	}

	if duration < 0 {
		timeStr = fmt.Sprintf("-%s", timeStr)
	}
	return
}
