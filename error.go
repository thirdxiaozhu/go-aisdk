/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-16 14:41:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 21:07:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk/httpclient"
)

var (
	ErrFailedToCreateConfigManager  = errors.New("failed to create config manager")                                                    // 创建配置管理器失败
	ErrFailedToCreateFlakeInstance  = errors.New("failed to create flake instance")                                                    // 创建分布式唯一ID生成器失败
	ErrProviderNotSupported         = errors.New("provider is not supported")                                                          // 提供商不支持
	ErrModelTypeNotSupported        = errors.New("model type is not supported")                                                        // 模型类型不支持
	ErrModelNotSupported            = errors.New("model is not supported")                                                             // 模型不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrMethodNotSupported           = errors.New("method is not supported")                                                            // 方法不支持
)

// SDKError SDK错误
type SDKError struct {
	RequestID string // 请求ID
	Err       error  // 原始错误
}

// Error 错误信息
func (e *SDKError) Error() (errStr string) {
	return fmt.Sprintf("request_id: %s, error: %v", e.RequestID, e.Err)
}

// RequestID 获取请求ID
func RequestID(err error) (requestId string) {
	if err == nil {
		return ""
	}

	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.RequestID
	}

	return ""
}

// Unwrap 解包错误
func Unwrap(err error) (originalError error) {
	if err == nil {
		return nil
	}
	// 解包 SDKError
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		if sdkErr.Err != nil {
			return sdkErr.Err
		}
		return err // 如果内部错误为 nil，返回 SDKError 本身
	}
	// 解包 RequestError
	var requestError *httpclient.RequestError
	if errors.As(err, &requestError) {
		if requestError.Err != nil {
			return requestError.Err
		}
		return err // 如果内部错误为 nil，返回 RequestError 本身
	}
	// 其他类型的错误
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err // 已经是最底层错误，返回原错误
	}
	return unwrapped
}

// Cause 错误根因
func Cause(err error) (causeError error) {
	return doCause(err)
}

// doCause 递归获取错误根因
func doCause(err error) (causeError error) {
	if err == nil {
		return nil
	}
	// 解包错误
	unwrapped := Unwrap(err)
	if unwrapped == nil {
		return err // 已经到达最底层错误，返回当前错误
	}
	// 防止无限递归：如果解包后的错误与原错误相同，直接返回
	if unwrapped == err {
		return err
	}
	return doCause(unwrapped)
}

// wrapFailedToCreateConfigManager 包装创建配置管理器失败错误
func wrapFailedToCreateConfigManager(text string) (err error) {
	return fmt.Errorf("%s: %w", text, ErrFailedToCreateConfigManager)
}

// wrapFailedToCreateFlakeInstance 包装创建分布式唯一ID生成器失败错误
func wrapFailedToCreateFlakeInstance(text string) (err error) {
	return fmt.Errorf("%s: %w", text, ErrFailedToCreateFlakeInstance)
}

// wrapProviderNotSupported 包装提供商不支持错误
func wrapProviderNotSupported(provider fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] is not supported: %w", provider.String(), ErrProviderNotSupported)
}

// wrapModelTypeNotSupported 包装模型类型不支持错误
func wrapModelTypeNotSupported(provider, modelType fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] does not support model type [%s]: %w",
		provider.String(), modelType.String(), ErrModelTypeNotSupported)
}

// wrapModelNotSupported 包装模型不支持错误
func wrapModelNotSupported(provider fmt.Stringer, model string, modelType fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s]: %w",
		provider.String(), model, modelType.String(), ErrModelNotSupported)
}

// wrapMethodNotSupported 包装方法不支持错误
func wrapMethodNotSupported(provider fmt.Stringer, model string, modelType fmt.Stringer, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s] with method [%s]: %w",
		provider.String(), model, modelType.String(), method, ErrMethodNotSupported)
}

// wrapMethodNotSupportedByProvider 包装方法不支持错误
func wrapMethodNotSupportedByProvider(provider fmt.Stringer, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] method: %w", provider.String(), method, ErrMethodNotSupported)
}
