/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-29 15:47:18
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/liusuxian/go-aisdk/consts"

// ModelInfo 模型信息结构
type ModelInfo struct {
	Provider  consts.Provider  `json:"provider"`   // 提供商
	ModelType consts.ModelType `json:"model_type"` // 模型类型
	Model     string           `json:"model"`      // 模型名称
}

// ToolType 工具类型
type ToolType string

const (
	ToolTypeFunction ToolType = "function" // 函数
)

// ToolCallsFunction 函数
type ToolCallsFunction struct {
	Arguments string `json:"arguments"` // 函数参数
	Name      string `json:"name"`      // 函数名
}

// ToolCalls 工具调用
type ToolCalls struct {
	Function *ToolCallsFunction `json:"function"` // 函数调用
	ID       string             `json:"id"`       // 工具ID
	Type     ToolType           `json:"type"`     // 工具类型
}

// CompletionTokensDetails completion tokens 的详细信息
type CompletionTokensDetails struct {
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"` // 使用预测输出时，预测中出现在 completion 中的 token 数量
	AudioTokens              int `json:"audio_tokens,omitempty"`               // 模型生成的音频输入 token 数
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`           // 模型生成的用于推理的 token 数
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"` // 使用预测输出时，预测中未出现在 completion 中的 token 数量
}

// PromptTokensDetails prompt tokens 的详细信息
type PromptTokensDetails struct {
	AudioTokens  int `json:"audio_tokens"`  // prompt 中存在的音频输入 token 数
	CachedTokens int `json:"cached_tokens"` // prompt 中存在的缓存 token 数
}

// Usage 该对话补全请求的用量信息
type Usage struct {
	CompletionTokens        int                     `json:"completion_tokens"`                  // 模型 completion 产生的 token 数
	PromptTokens            int                     `json:"prompt_tokens"`                      // 用户 prompt 所包含的 token 数
	PromptCacheHitTokens    int                     `json:"prompt_cache_hit_tokens,omitempty"`  // 用户 prompt 中，命中上下文缓存的 token 数
	PromptCacheMissTokens   int                     `json:"prompt_cache_miss_tokens,omitempty"` // 用户 prompt 中，未命中上下文缓存的 token 数
	TotalTokens             int                     `json:"total_tokens"`                       // 该请求中，所有 token 的数量（prompt + completion）
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`          // completion tokens 的详细信息
	PromptTokensDetails     *PromptTokensDetails    `json:"prompt_tokens_details,omitempty"`    // prompt tokens 的详细信息
}
