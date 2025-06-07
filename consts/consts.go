/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:57:28
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 18:57:30
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package consts

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
	Ark        Provider = "ark"
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
