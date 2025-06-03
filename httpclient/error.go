/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:56:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-03 11:43:08
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient

import (
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
)

var (
	ErrProviderNotSupported         = errors.New("provider is not supported")                                                          // 提供商不支持
	ErrModelTypeNotSupported        = errors.New("model type is not supported")                                                        // 模型类型不支持
	ErrModelNotSupported            = errors.New("model is not supported")                                                             // 模型不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrMethodNotSupported           = errors.New("method is not supported")                                                            // 方法不支持
	ErrTooManyEmptyStreamMessages   = errors.New("stream has sent too many empty messages")                                            // 流式传输发送了太多空消息
)

// WrapProviderNotSupported 包装提供商不支持错误
func WrapProviderNotSupported(provider consts.Provider) (err error) {
	return fmt.Errorf("provider [%s] is not supported: %w", provider, ErrProviderNotSupported)
}

// WrapModelTypeNotSupported 包装模型类型不支持错误
func WrapModelTypeNotSupported(provider consts.Provider, modelType consts.ModelType) (err error) {
	return fmt.Errorf("provider [%s] does not support model type [%s]: %w",
		provider, modelType, ErrModelTypeNotSupported)
}

// WrapModelNotSupported 包装模型不支持错误
func WrapModelNotSupported(provider consts.Provider, model string, modelType consts.ModelType) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s]: %w",
		provider, model, modelType, ErrModelNotSupported)
}

// WrapMethodNotSupported 包装方法不支持错误
func WrapMethodNotSupported(provider consts.Provider, model string, modelType consts.ModelType, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s] with method [%s]: %w",
		provider, model, modelType, method, ErrMethodNotSupported)
}

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
