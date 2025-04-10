/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 20:13:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-10 13:30:16
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import "encoding/json"

// ChatMessage 聊天消息的通用接口
type ChatMessage interface {
	GetRole() (role string)             // 获取消息角色
	MarshalJSON() (b []byte, err error) // 序列化JSON
}

// ChatAudioFormatType 音频输出格式
type ChatAudioFormatType string

const (
	ChatAudioFormatTypeWAV   ChatAudioFormatType = "wav"
	ChatAudioFormatTypeMP3   ChatAudioFormatType = "mp3"
	ChatAudioFormatTypeFLAC  ChatAudioFormatType = "flac"
	ChatAudioFormatTypeOPUS  ChatAudioFormatType = "opus"
	ChatAudioFormatTypePCM16 ChatAudioFormatType = "pcm16"
)

// ChatAudioVoiceType 模型回应时使用的声音
type ChatAudioVoiceType string

const (
	ChatAudioVoiceTypeAlloy   ChatAudioVoiceType = "alloy"
	ChatAudioVoiceTypeAsh     ChatAudioVoiceType = "ash"
	ChatAudioVoiceTypeBallad  ChatAudioVoiceType = "ballad"
	ChatAudioVoiceTypeCoral   ChatAudioVoiceType = "coral"
	ChatAudioVoiceTypeEcho    ChatAudioVoiceType = "echo"
	ChatAudioVoiceTypeSage    ChatAudioVoiceType = "sage"
	ChatAudioVoiceTypeShimmer ChatAudioVoiceType = "shimmer"
)

// ChatAudio 音频输出
type ChatAudio struct {
	Format ChatAudioFormatType `json:"format"` // 输出音频格式
	Voice  ChatAudioVoiceType  `json:"voice"`  // 模型回应时使用的声音
}

// ChatPredictionType 预测内容的类型
type ChatPredictionType string

const (
	ChatPredictionTypeContent ChatPredictionType = "content"
)

// ChatPredictionContentPart 预测内容
type ChatPredictionContentPart struct {
	Type string `json:"type"` // 内容的类型
	Text string `json:"text"` // 文本内容
}

// ChatPrediction 预测输出配置
type ChatPrediction struct {
	Type    ChatPredictionType          `json:"type"`    // 预测内容的类型
	Content []ChatPredictionContentPart `json:"content"` // 预测内容
}

// ChatReasoningEffortType 推理努力程度
type ChatReasoningEffortType string

const (
	ChatReasoningEffortTypeLow    ChatReasoningEffortType = "low"
	ChatReasoningEffortTypeMedium ChatReasoningEffortType = "medium"
	ChatReasoningEffortTypeHigh   ChatReasoningEffortType = "high"
)

// ChatResponseFormatType 响应格式的类型
type ChatResponseFormatType string

const (
	ChatResponseFormatTypeText       ChatResponseFormatType = "text"
	ChatResponseFormatTypeJSONSchema ChatResponseFormatType = "json_schema"
	ChatResponseFormatTypeJSONObject ChatResponseFormatType = "json_object"
)

// ChatResponseFormatJSONSchema JSON Schema 配置
type ChatResponseFormatJSONSchema struct {
	Name        string         `json:"name"`                  // 响应格式名称，必须是 a-z、A-Z、0-9 或包含下划线和破折号，最大长度为 64
	Description string         `json:"description,omitempty"` // 响应格式的描述，用于指导模型如何响应
	Schema      json.Marshaler `json:"schema,omitempty"`      // 响应格式的 JSON Schema
	Strict      bool           `json:"strict,omitempty"`      // 是否启用严格模式，默认为 false
}

// ChatResponseFormat 响应格式
type ChatResponseFormat struct {
	Type       ChatResponseFormatType        `json:"type"`                  // 响应格式的类型
	JSONSchema *ChatResponseFormatJSONSchema `json:"json_schema,omitempty"` // JSON Schema 配置，仅当 Type 为 "json_schema" 时使用
}

// ChatStreamOptions 流式传输选项
type ChatStreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"` // 是否包含令牌使用统计信息
}

// ChatToolChoiceType 工具调用类型
type ChatToolChoiceType string

const (
	ChatToolChoiceTypeNone     ChatToolChoiceType = "none"     // 模型不会调用任何工具，而是生成消息
	ChatToolChoiceTypeAuto     ChatToolChoiceType = "auto"     // 模型可以在生成消息或调用一个或多个工具之间选择
	ChatToolChoiceTypeRequired ChatToolChoiceType = "required" // 模型必须调用一个或多个工具
)

// ChatToolChoiceFunction 工具调用函数
type ChatToolChoiceFunction struct {
	Name string `json:"name"` // 工具调用函数名称
}

// ChatToolChoice 模型是否调用工具
type ChatToolChoice struct {
	ToolChoiceType ChatToolChoiceType      // 工具调用类型
	Function       *ChatToolChoiceFunction // 工具调用函数
	Type           ToolType                // 工具类型
}

// MarshalJSON 序列化JSON
func (c ChatToolChoice) MarshalJSON() (b []byte, err error) {
	if c.ToolChoiceType != "" {
		return json.Marshal(c.ToolChoiceType)
	}
	return json.Marshal(struct {
		Function *ChatToolChoiceFunction `json:"function"`
		Type     ToolType                `json:"type"`
	}{
		Function: c.Function,
		Type:     c.Type,
	})
}

// ChatToolFunction 工具函数
type ChatToolFunction struct {
	Name        string         `json:"name"`                  // 函数名称，必须是 a-z, A-Z, 0-9 或者包含下划线和破折号，最大长度为 64
	Description string         `json:"description,omitempty"` // 函数描述，用于帮助模型决定何时以及如何调用函数
	Parameters  map[string]any `json:"parameters,omitempty"`  // 函数接受的参数，描述为一个 JSON Schema 对象
	Strict      bool           `json:"strict,omitempty"`      // 是否启用严格模式，默认为 false
}

// ChatTool 工具
type ChatTool struct {
	Function *ChatToolFunction `json:"function"` // 工具函数
	Type     ToolType          `json:"type"`     // 工具类型
}

// ChatSearchContextSize 搜索上下文大小
type ChatSearchContextSize string

const (
	ChatSearchContextSizeLow    ChatSearchContextSize = "low"
	ChatSearchContextSizeMedium ChatSearchContextSize = "medium"
	ChatSearchContextSizeHigh   ChatSearchContextSize = "high"
)

// ChatApproximateLocation 用户的大致位置参数
type ChatApproximateLocation struct {
	City     string `json:"city,omitempty"`     // 用户所在城市
	Country  string `json:"country,omitempty"`  // 用户所在国家的两字母 ISO 代码
	Region   string `json:"region,omitempty"`   // 用户所在地区
	Timezone string `json:"timezone,omitempty"` // 用户的 IANA 时区
}

// ChatApproximateLocationType 位置近似类型
type ChatApproximateLocationType string

const (
	ChatApproximateLocationTypeApproximate ChatApproximateLocationType = "approximate"
)

// ChatUserLocation 用户位置信息
type ChatUserLocation struct {
	Approximate *ChatApproximateLocation    `json:"approximate"` // 大致位置信息
	Type        ChatApproximateLocationType `json:"type"`        // 位置近似类型
}

// ChatWebSearchOptions 网络搜索选项
type ChatWebSearchOptions struct {
	SearchContextSize ChatSearchContextSize `json:"search_context_size,omitempty"` // 搜索上下文大小
	UserLocation      *ChatUserLocation     `json:"user_location,omitempty"`       // 用户位置信息
}

// ChatRequest 聊天请求
type ChatRequest struct {
	BaseRequest
	Messages            []ChatMessage           `json:"messages"`                        // 消息数组
	Audio               *ChatAudio              `json:"audio,omitempty"`                 // 音频输出
	FrequencyPenalty    float32                 `json:"frequency_penalty,omitempty"`     // 介于 -2.0 和 2.0 之间的数值。正值会根据文本中已有内容的出现频率对新 token 进行惩罚，从而降低模型逐字重复相同内容的可能性
	LogitBias           map[string]int          `json:"logit_bias,omitempty"`            // 修改指定标记在补全中出现的可能性
	LogProbs            bool                    `json:"logprobs,omitempty"`              // 是否返回所输出 token 的对数概率
	MaxCompletionTokens int                     `json:"max_completion_tokens,omitempty"` // 生成补全内容的最大令牌数上限，包括可见的输出令牌和推理令牌
	Metadata            map[string]string       `json:"metadata,omitempty"`              // 元数据
	Modalities          []string                `json:"modalities,omitempty"`            // 希望模型生成的输出类型
	N                   int                     `json:"n,omitempty"`                     // 为每个输入消息生成的聊天完成选项数量
	ParallelToolCalls   bool                    `json:"parallel_tool_calls,omitempty"`   // 是否在使用工具时启用并行函数调用
	Prediction          *ChatPrediction         `json:"prediction,omitempty"`            // 预测输出配置
	PresencePenalty     float32                 `json:"presence_penalty,omitempty"`      // 介于 -2.0 和 2.0 之间的数值。正值会根据新标记是否已在文本中出现过对其进行惩罚，从而增加模型讨论新话题的可能性
	ReasoningEffort     ChatReasoningEffortType `json:"reasoning_effort,omitempty"`      // 仅适用于 o 系列模型，约束推理模型的推理努力程度
	ResponseFormat      *ChatResponseFormat     `json:"response_format,omitempty"`       // 响应格式
	Seed                int                     `json:"seed,omitempty"`                  // 随机种子
	ServiceTier         string                  `json:"service_tier,omitempty"`          // 指定用于处理请求的延迟层级。此参数与订阅了规模层级服务的客户相关
	Stop                []string                `json:"stop,omitempty"`                  // 最多4个序列，当API遇到这些序列时将停止生成更多标记。返回的文本不会包含停止序列
	Store               bool                    `json:"store,omitempty"`                 // 是否存储此聊天完成请求的输出，用于我们的模型蒸馏或评估产品
	Stream              bool                    `json:"stream,omitempty"`                // 是否流式传输响应
	StreamOptions       *ChatStreamOptions      `json:"stream_options,omitempty"`        // 流式传输选项
	Temperature         float32                 `json:"temperature,omitempty"`           // 采样温度值，范围在0到2之间
	ToolChoice          *ChatToolChoice         `json:"tool_choice,omitempty"`           // 模型是否调用工具
	Tools               []ChatTool              `json:"tools,omitempty"`                 // 工具列表
	TopLogProbs         int                     `json:"top_logprobs,omitempty"`          // 一个介于0和20之间的整数，指定在每个标记位置返回的最可能标记的数量，每个标记都有相关的对数概率。如果使用此参数，必须将logprobs设置为true
	TopP                float32                 `json:"top_p,omitempty"`                 // 一种替代温度采样的方法，我们通常建议调整此参数或温度（temperature），但不要同时调整两者
	User                string                  `json:"user,omitempty"`                  // 代表你的终端用户的唯一标识符
	WebSearchOptions    *ChatWebSearchOptions   `json:"web_search_options,omitempty"`    // 网络搜索选项
}

// ChatResponse 聊天响应
type ChatResponse struct {
	BaseResponse
}
