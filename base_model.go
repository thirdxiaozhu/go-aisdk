/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 16:37:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-10 14:10:40
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

// BaseRequest 基础请求结构
type BaseRequest struct {
	Model     string    `json:"model"`      // 模型名称
	Provider  Provider  `json:"provider"`   // 提供商
	ModelType ModelType `json:"model_type"` // 模型类型
}

// GetModel 获取模型名称
func (b BaseRequest) GetModel() (model string) {
	return b.Model
}

// GetProvider 获取提供商
func (b BaseRequest) GetProvider() (provider Provider) {
	return b.Provider
}

// GetModelType 获取模型类型
func (b BaseRequest) GetModelType() (modelType ModelType) {
	return b.ModelType
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Error error `json:"error,omitempty"` // 错误
}

// GetError 获取错误
func (b BaseResponse) GetError() (err error) {
	return b.Error
}

// ToolType 工具类型
type ToolType string

const (
	ToolTypeFunction ToolType = "function" // 函数
)
