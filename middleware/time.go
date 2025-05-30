/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 17:42:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 17:42:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package middleware

import (
	"fmt"
	"time"
)

// formatDuration 格式化时间，根据时长选择合适的单位
func formatDuration(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%.0fms", float64(duration.Nanoseconds())/1e6)
	} else if duration < time.Minute {
		return fmt.Sprintf("%.2fs", duration.Seconds())
	} else {
		return fmt.Sprintf("%.2fm", duration.Minutes())
	}
}
