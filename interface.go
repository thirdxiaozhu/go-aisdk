/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 15:01:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-11 16:36:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import "context"

// Provider AI服务提供商类型
type Provider string

const (
	OpenAI     Provider = "openai"     // OpenAI
	DeepSeek   Provider = "deepseek"   // DeepSeek
	Claude     Provider = "claude"     // Anthropic Claude
	Gemini     Provider = "gemini"     // Google Gemini
	Aliyunbl   Provider = "aliyunbl"   // 阿里云百炼
	Midjourney Provider = "midjourney" // Midjourney
	Vidu       Provider = "vidu"       // 生数科技
	Keling     Provider = "keling"     // 可灵AI
)

// ModelType 模型类型
type ModelType string

const (
	ChatModel  ModelType = "chat"  // 对话模型
	ImageModel ModelType = "image" // 图像生成模型
	VideoModel ModelType = "video" // 视频生成模型
	AudioModel ModelType = "audio" // 音频处理模型
	EmbedModel ModelType = "embed" // 嵌入模型
)

// ProviderService AI服务提供商的服务接口
type ProviderService interface {
	GetProvider() (provider Provider) // 获取提供商

	// 聊天相关
	CreateChatCompletion(ctx context.Context, request ChatRequest) (response ChatResponse, err error)
	// CreateChatCompletionStream(ctx context.Context, request ChatRequest) (response ChatResponseStream, err error)

	// TODO 图像相关

	// TODO 视频相关

	// TODO 音频相关
}
