/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 20:13:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-09 21:27:50
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

// ChatMessage 聊天消息的通用接口
type ChatMessage interface {
	GetRole() (role string)             // 获取消息角色
	MarshalJSON() (b []byte, err error) // 序列化JSON
}

// AudioOutputFormatType 音频输出格式
type AudioOutputFormatType string

const (
	AudioOutputFormatTypeWAV   AudioOutputFormatType = "wav"
	AudioOutputFormatTypeMP3   AudioOutputFormatType = "mp3"
	AudioOutputFormatTypeFLAC  AudioOutputFormatType = "flac"
	AudioOutputFormatTypeOPUS  AudioOutputFormatType = "opus"
	AudioOutputFormatTypePCM16 AudioOutputFormatType = "pcm16"
)

// AudioOutputVoiceType 模型回应时使用的声音
type AudioOutputVoiceType string

const (
	AudioOutputVoiceTypeAlloy   AudioOutputVoiceType = "alloy"
	AudioOutputVoiceTypeAsh     AudioOutputVoiceType = "ash"
	AudioOutputVoiceTypeBallad  AudioOutputVoiceType = "ballad"
	AudioOutputVoiceTypeCoral   AudioOutputVoiceType = "coral"
	AudioOutputVoiceTypeEcho    AudioOutputVoiceType = "echo"
	AudioOutputVoiceTypeSage    AudioOutputVoiceType = "sage"
	AudioOutputVoiceTypeShimmer AudioOutputVoiceType = "shimmer"
)

// AudioOutput 音频输出
type AudioOutput struct {
	Format AudioOutputFormatType `json:"format"` // 输出音频格式
	Voice  AudioOutputVoiceType  `json:"voice"`  // 模型回应时使用的声音
}

// PredictionType 预测内容的类型
type PredictionType string

const (
	PredictionTypeContent PredictionType = "content"
)

// PredictionContentPart 预测内容部分
type PredictionContentPart struct {
	Type string `json:"type"` // 内容部分的类型
	Text string `json:"text"` // 文本内容
}

// Prediction 预测输出配置
type Prediction struct {
	Type        PredictionType          `json:"type"`                   // 预测内容的类型
	Content     string                  `json:"content,omitempty"`      // 预测内容字符串
	ContentList []PredictionContentPart `json:"content_list,omitempty"` // 预测内容数组
}

// ChatRequest 聊天请求
type ChatRequest struct {
	BaseRequest
	Messages            []ChatMessage     `json:"messages"`                        // 消息数组
	Audio               *AudioOutput      `json:"audio,omitempty"`                 // 音频输出
	FrequencyPenalty    float32           `json:"frequency_penalty,omitempty"`     // 介于 -2.0 和 2.0 之间的数值。正值会根据文本中已有内容的出现频率对新 token 进行惩罚，从而降低模型逐字重复相同内容的可能性
	LogitBias           map[string]int    `json:"logit_bias,omitempty"`            // 修改指定标记在补全中出现的可能性
	LogProbs            bool              `json:"logprobs,omitempty"`              // 是否返回输出标记的对数概率
	MaxCompletionTokens int               `json:"max_completion_tokens,omitempty"` // 生成补全内容的最大令牌数上限，包括可见的输出令牌和推理令牌
	Metadata            map[string]string `json:"metadata,omitempty"`              // 元数据
	Modalities          []string          `json:"modalities,omitempty"`            // 希望模型生成的输出类型
	N                   int               `json:"n,omitempty"`                     // 为每个输入消息生成的聊天完成选项数量
	ParallelToolCalls   bool              `json:"parallel_tool_calls,omitempty"`   // 是否在使用工具时启用并行函数调用
	Prediction          *Prediction       `json:"prediction,omitempty"`            // 预测输出配置
	// TODO presence_penalty
}

// ChatResponse 聊天响应
type ChatResponse struct {
	BaseResponse
}
