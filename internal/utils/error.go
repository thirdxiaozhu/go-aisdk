/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:26:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 18:24:13
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"errors"
	"fmt"
)

var (
	ErrMethodNotSupported          = errors.New("method is not supported")                 // 方法不支持
	ErrTooManyEmptyStreamMessages  = errors.New("stream has sent too many empty messages") // 流式传输发送了太多空消息
	ErrStreamReturnIntervalTimeout = errors.New("stream return interval timeout")          // 流式传输返回间隔超时
)

// WrapMethodNotSupported 包装方法不支持错误
func WrapMethodNotSupported(provider, modelType fmt.Stringer, model, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s] with method [%s]: %w",
		provider.String(), model, modelType.String(), method, ErrMethodNotSupported)
}

// WrapMethodNotSupportedByProvider 包装方法不支持错误
func WrapMethodNotSupportedByProvider(provider fmt.Stringer, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] method: %w", provider.String(), method, ErrMethodNotSupported)
}
