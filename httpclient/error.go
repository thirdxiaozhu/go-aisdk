/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:56:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-06 01:06:37
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages") // 流式传输发送了太多空消息
)

// RequestError 请求错误
type RequestError struct {
	HTTPStatus     string // HTTP 状态描述
	HTTPStatusCode int    // HTTP 状态码
	Err            error  // 错误信息
	Body           []byte // 响应体
}

// Error 实现 error 接口的方法
func (e *RequestError) Error() (s string) {
	return fmt.Sprintf(
		"error, status code: %d, status: %s, message: %s, body: %s",
		e.HTTPStatusCode, e.HTTPStatus, e.Err, e.Body,
	)
}

// Unwrap 解包错误
func (e *RequestError) Unwrap() (err error) {
	return e.Err
}
