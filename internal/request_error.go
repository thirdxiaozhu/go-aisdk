/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 20:26:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-07 20:27:00
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import "fmt"

// RequestError 请求错误
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

// Error
func (re *RequestError) Error() (text string) {
	return fmt.Sprintf("error, status code: %d, message: %v", re.HTTPStatusCode, re.Err)
}

// Unwrap
func (re *RequestError) Unwrap() (err error) {
	return re.Err
}
