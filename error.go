/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-16 14:41:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 14:43:56
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"errors"
	"fmt"
)

var (
	ErrFailedToCreateFlakeInstance  = errors.New("failed to create flake instance")                                                    // 创建分布式唯一ID生成器失败
	ErrProviderNotSupported         = errors.New("provider is not supported")                                                          // 提供商不支持
	ErrModelTypeNotSupported        = errors.New("model type is not supported")                                                        // 模型类型不支持
	ErrModelNotSupported            = errors.New("model is not supported")                                                             // 模型不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrMethodNotSupported           = errors.New("method is not supported")                                                            // 方法不支持
	ErrEmptyAPIKeyList              = errors.New("api key list is empty")                                                              // API密钥列表为空
	ErrNoAPIKeyAvailable            = errors.New("no api key available")                                                               // 没有可用的API密钥
	ErrAPIKeyNotFound               = errors.New("api key not found")                                                                  // API密钥不存在
	ErrAPIKeyAlreadyExists          = errors.New("api key already exists")                                                             // API密钥已存在
	ErrWeightMustBeGreaterThan0     = errors.New("weight must be greater than 0")                                                      // 权重必须大于0
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

	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Err
	}

	return err
}

// Cause 错误根因
func Cause(err error) (causeError error) {
	originalError := Unwrap(err)
	return doCause(originalError)
}

// doCause 递归获取错误根因
func doCause(err error) (causeError error) {
	if err == nil {
		return nil
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err // 已经到达最底层错误，返回当前错误
	}

	return doCause(unwrapped)
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
