/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-29 15:50:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-29 16:25:59
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package consts

import (
	"errors"
	"fmt"
)

var (
	ErrProviderNotSupported         = errors.New("provider is not supported")                                                          // 提供商不支持
	ErrModelTypeNotSupported        = errors.New("model type is not supported")                                                        // 模型类型不支持
	ErrModelNotSupported            = errors.New("model is not supported")                                                             // 模型不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrCompletionNotStream          = errors.New("stream completion header is not stream")
	ErrMethodNotSupported           = errors.New("method is not supported") // 方法不支持
)

// WrapProviderNotSupported 包装提供商不支持错误
func WrapProviderNotSupported(provider Provider) (err error) {
	return fmt.Errorf("provider [%s] is not supported: %w", provider, ErrProviderNotSupported)
}

// WrapModelTypeNotSupported 包装模型类型不支持错误
func WrapModelTypeNotSupported(provider Provider, modelType ModelType) (err error) {
	return fmt.Errorf("provider [%s] does not support model type [%s]: %w",
		provider, modelType, ErrModelTypeNotSupported)
}

// WrapModelNotSupported 包装模型不支持错误
func WrapModelNotSupported(provider Provider, model string, modelType ModelType) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s]: %w",
		provider, model, modelType, ErrModelNotSupported)
}

// WrapMethodNotSupported 包装方法不支持错误
func WrapMethodNotSupported(provider Provider, model string, modelType ModelType, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s] with method [%s]: %w",
		provider, model, modelType, method, ErrMethodNotSupported)
}
