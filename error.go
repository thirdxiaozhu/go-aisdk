/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-05 19:18:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-05 19:23:55
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"errors"
	"fmt"
)

// SDKError SDK错误
type SDKError struct {
	RequestID string // 请求ID
	Err       error  // 原始错误
}

// Error 错误信息
func (e *SDKError) Error() (s string) {
	return fmt.Sprintf("request_id: %s, error: %v", e.RequestID, e.Err)
}

// Unwrap 解包错误
func (e *SDKError) Unwrap() (err error) {
	return e.Err
}

// GetRequestID 获取请求ID
func GetRequestID(err error) (requestId string) {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.RequestID
	}
	return ""
}
