/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 22:47:43
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
)

var (
	// 序列化聊天请求函数（OpenAI）
	marshalChatRequestByOpenAI = func(r ChatRequest) (b []byte, err error) {
		// 设置提供商
		for _, message := range r.Messages {
			message.SetProvider(r.Provider.String())
		}
		// 创建一个别名结构体
		type Alias ChatRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			User   string `json:"user,omitempty"` // 代表你的终端用户的唯一标识符
			Alias
		}{
			User:  r.UserInfo.UserID,
			Alias: Alias(r),
		}
		// 处理公共字段
		if r.WebSearchOptions != nil {
			temp.WebSearchOptions.ForcedSearch = false
			temp.WebSearchOptions.SearchStrategy = ""
		}
		// 移除不支持的字段
		temp.TopK = 0
		temp.EnableThinking = false
		temp.ThinkingBudget = 0
		temp.TranslationOptions = nil
		temp.XDashScopeDataInspection = ""
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}
	// 序列化聊天请求函数（DeepSeek）
	marshalChatRequestByDeepSeek = func(r ChatRequest) (b []byte, err error) {
		// 设置提供商
		for _, message := range r.Messages {
			message.SetProvider(r.Provider.String())
		}
		// 创建一个别名结构体
		type Alias ChatRequest
		temp := struct {
			UserID    string `json:"user_id,omitempty"`
			MaxTokens int    `json:"max_tokens,omitempty"`
			Alias
		}{
			MaxTokens: r.MaxCompletionTokens,
			Alias:     Alias(r),
		}
		// 处理公共字段
		if r.ResponseFormat != nil && r.ResponseFormat.JSONSchema != nil {
			temp.ResponseFormat = nil
		}
		if len(r.Tools) > 0 {
			tempTools := make([]ChatTool, 0, len(r.Tools))
			for _, v := range r.Tools {
				if v.Function != nil {
					v.Function.Strict = false
				}
				tempTools = append(tempTools, v)
			}
			temp.Tools = tempTools
		}
		// 移除不支持的字段
		temp.Audio = nil
		temp.LogitBias = nil
		temp.MaxCompletionTokens = 0
		temp.Metadata = nil
		temp.Modalities = nil
		temp.N = 0
		temp.ParallelToolCalls = false
		temp.Prediction = nil
		temp.ReasoningEffort = ""
		temp.Seed = 0
		temp.ServiceTier = ""
		temp.Store = false
		temp.WebSearchOptions = nil
		temp.TopK = 0
		temp.EnableThinking = false
		temp.ThinkingBudget = 0
		temp.TranslationOptions = nil
		temp.XDashScopeDataInspection = ""
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}
	// 序列化聊天请求函数（AliBL）
	marshalChatRequestByAliBL = func(r ChatRequest) (b []byte, err error) {
		// 设置提供商
		for _, message := range r.Messages {
			message.SetProvider(r.Provider.String())
		}
		// 创建一个别名结构体
		type Alias ChatRequest
		temp := struct {
			UserID        string                `json:"user_id,omitempty"`
			MaxTokens     int                   `json:"max_tokens,omitempty"`
			EnableSearch  bool                  `json:"enable_search,omitempty"`
			SearchOptions *ChatWebSearchOptions `json:"search_options,omitempty"`
			Alias
		}{
			MaxTokens:    r.MaxCompletionTokens,
			EnableSearch: r.WebSearchOptions != nil,
			Alias:        Alias(r),
		}
		// 处理公共字段
		if r.ResponseFormat != nil && r.ResponseFormat.JSONSchema != nil {
			temp.ResponseFormat = nil
		}
		if len(r.Tools) > 0 {
			tempTools := make([]ChatTool, 0, len(r.Tools))
			for _, v := range r.Tools {
				if v.Function != nil {
					v.Function.Strict = false
				}
				tempTools = append(tempTools, v)
			}
			temp.Tools = tempTools
		}
		if r.WebSearchOptions != nil {
			temp.SearchOptions = &ChatWebSearchOptions{
				ForcedSearch:   r.WebSearchOptions.ForcedSearch,
				SearchStrategy: r.WebSearchOptions.SearchStrategy,
			}
			temp.WebSearchOptions = nil
		}
		// 移除不支持的字段
		temp.FrequencyPenalty = 0
		temp.LogitBias = nil
		temp.MaxCompletionTokens = 0
		temp.Metadata = nil
		temp.Prediction = nil
		temp.ReasoningEffort = ""
		temp.ServiceTier = ""
		temp.Store = false
		temp.XDashScopeDataInspection = ""
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}
	// 序列化聊天请求函数（Ark）
	marshalChatRequestByArk = func(r ChatRequest) (b []byte, err error) {
		// 设置提供商
		for _, message := range r.Messages {
			message.SetProvider(r.Provider.String())
		}
		// 创建一个别名结构体
		type Alias ChatRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			Alias
		}{
			Alias: Alias(r),
		}
		// 处理公共字段
		if r.ResponseFormat != nil && r.ResponseFormat.JSONSchema != nil {
			temp.ResponseFormat = nil
		}
		if len(r.Tools) > 0 {
			tempTools := make([]ChatTool, 0, len(r.Tools))
			for _, v := range r.Tools {
				if v.Function != nil {
					v.Function.Strict = false
				}
				tempTools = append(tempTools, v)
			}
			temp.Tools = tempTools
		}
		// 移除不支持的字段
		temp.Metadata = nil
		temp.Prediction = nil
		temp.ReasoningEffort = ""
		temp.Store = false
		temp.XDashScopeDataInspection = ""
		// 序列化JSON
		temp.Provider = ""
		temp.EnableThinking = false
		return json.Marshal(temp)
	}
	// 策略映射
	chatRequestStrategies = map[consts.Provider]func(m ChatRequest) (b []byte, err error){
		consts.OpenAI:   marshalChatRequestByOpenAI,
		consts.DeepSeek: marshalChatRequestByDeepSeek,
		consts.AliBL:    marshalChatRequestByAliBL,
		consts.Ark:      marshalChatRequestByArk,
	}
)

// ChatMessage 聊天消息的通用接口
type ChatMessage interface {
	SetProvider(provider string)        // 设置提供商（序列化参数时，处理差异化参数）
	MarshalJSON() (b []byte, err error) // 序列化JSON
}

// ChatAudioFormatType 输出音频的格式
//
//	提供商支持: OpenAI | AliBL
type ChatAudioFormatType string

const (
	// 提供商支持: OpenAI | AliBL
	ChatAudioFormatTypeWAV ChatAudioFormatType = "wav"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeMP3 ChatAudioFormatType = "mp3"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeFLAC ChatAudioFormatType = "flac"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeOPUS ChatAudioFormatType = "opus"
	// 提供商支持: OpenAI
	ChatAudioFormatTypePCM16 ChatAudioFormatType = "pcm16"
)

// ChatAudioVoiceType 输出音频的音色
//
//	提供商支持: OpenAI | AliBL
type ChatAudioVoiceType string

const (
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeAlloy ChatAudioVoiceType = "alloy"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeAsh ChatAudioVoiceType = "ash"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeBallad ChatAudioVoiceType = "ballad"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeCoral ChatAudioVoiceType = "coral"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeEcho ChatAudioVoiceType = "echo"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeFable ChatAudioVoiceType = "fable"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeNova ChatAudioVoiceType = "nova"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeOnyx ChatAudioVoiceType = "onyx"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeSage ChatAudioVoiceType = "sage"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeShimmer ChatAudioVoiceType = "shimmer"
	// 提供商支持: AliBL
	ChatAudioVoiceTypeCherry ChatAudioVoiceType = "Cherry"
	// 提供商支持: AliBL
	ChatAudioVoiceTypeSerena ChatAudioVoiceType = "Serena"
	// 提供商支持: AliBL
	ChatAudioVoiceTypeEthan ChatAudioVoiceType = "Ethan"
	// 提供商支持: AliBL
	ChatAudioVoiceTypeChelsie ChatAudioVoiceType = "Chelsie"
)

// ChatAudioOutputArgs 音频输出参数
//
//	提供商支持: OpenAI | AliBL
type ChatAudioOutputArgs struct {
	// 输出音频的格式
	// 提供商支持: OpenAI | AliBL
	Format ChatAudioFormatType `json:"format,omitempty"`
	// 输出音频的音色
	// 提供商支持: OpenAI | AliBL
	Voice ChatAudioVoiceType `json:"voice,omitempty"`
}

// ChatModalitiesType 输出数据的模态
//
//	提供商支持: OpenAI | AliBL
type ChatModalitiesType string

const (
	// 提供商支持: OpenAI | AliBL
	ChatModalitiesTypeText ChatModalitiesType = "text"
	// 提供商支持: OpenAI | AliBL
	ChatModalitiesTypeAudio ChatModalitiesType = "audio"
)

// ChatPredictionType 预测内容的类型
//
//	提供商支持: OpenAI
type ChatPredictionType string

const (
	// 提供商支持: OpenAI
	ChatPredictionTypeContent ChatPredictionType = "content"
)

// ChatPredictionContentPart 预测内容
//
//	提供商支持: OpenAI
type ChatPredictionContentPart struct {
	Type string `json:"type,omitempty"` // 内容的类型
	Text string `json:"text,omitempty"` // 文本内容
}

// ChatPrediction 预测输出配置
//
//	提供商支持: OpenAI
type ChatPrediction struct {
	// 预测内容的类型
	// 提供商支持: OpenAI
	Type ChatPredictionType `json:"type,omitempty"`
	// 预测内容
	// 提供商支持: OpenAI
	Content []ChatPredictionContentPart `json:"content,omitempty"`
}

// ChatReasoningEffortType 推理努力程度
//
//	提供商支持: OpenAI
type ChatReasoningEffortType string

const (
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeLow ChatReasoningEffortType = "low"
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeMedium ChatReasoningEffortType = "medium"
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeHigh ChatReasoningEffortType = "high"
)

// ChatResponseFormatType 响应格式的类型
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatResponseFormatType string

const (
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatResponseFormatTypeText ChatResponseFormatType = "text"
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatResponseFormatTypeJSONObject ChatResponseFormatType = "json_object"
	// 提供商支持: OpenAI
	ChatResponseFormatTypeJSONSchema ChatResponseFormatType = "json_schema"
)

// ChatResponseFormatJSONSchema JSON Schema 配置
//
//	提供商支持: OpenAI
type ChatResponseFormatJSONSchema struct {
	// 响应格式名称，必须是 a-z、A-Z、0-9 或包含下划线和破折号，最大长度为 64
	// 提供商支持: OpenAI | Ark
	Name string `json:"name,omitempty"`
	// 响应格式的描述，用于指导模型如何响应
	// 提供商支持: OpenAI | Ark
	Description string `json:"description,omitempty"`
	// 响应格式的 JSON Schema
	// 提供商支持: OpenAI | Ark
	Schema json.Marshaler `json:"schema,omitempty"`
	// 是否启用严格模式，默认为 false
	// 提供商支持: OpenAI | Arkf
	Strict bool `json:"strict,omitempty"`
}

// ChatResponseFormat 响应格式
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatResponseFormat struct {
	// 响应格式的类型
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Type ChatResponseFormatType `json:"type,omitempty"`
	// JSON Schema 配置，仅当 Type 为 "json_schema" 时使用
	// 提供商支持: OpenAI | Ark
	JSONSchema *ChatResponseFormatJSONSchema `json:"json_schema,omitempty"`
}

// ChatStreamOptions 流式传输选项
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatStreamOptions struct {
	// 提供商支持: OpenAI | DeepSeek | AliBL
	IncludeUsage bool `json:"include_usage,omitempty"` // 是否包含令牌使用统计信息
}

// ChatToolChoiceType 工具调用类型
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatToolChoiceType string

const (
	// 模型不会调用任何 tool，而是生成一条消息
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatToolChoiceTypeNone ChatToolChoiceType = "none"
	// 模型可以选择生成一条消息或调用一个或多个 tool
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatToolChoiceTypeAuto ChatToolChoiceType = "auto"
	// 模型必须调用一个或多个 tool
	// 提供商支持: OpenAI | DeepSeek
	ChatToolChoiceTypeRequired ChatToolChoiceType = "required"
)

// ChatToolChoiceFunction 工具调用函数
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatToolChoiceFunction struct {
	// 工具调用函数名称
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Name string `json:"name,omitempty"`
}

// ToolType 工具类型
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ToolType string

const (
	// 函数
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolTypeFunction ToolType = "function"
)

// ChatToolChoice 指定工具调用的策略
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatToolChoice struct {
	// 工具调用类型
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolChoiceType ChatToolChoiceType
	// 工具调用函数
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Function *ChatToolChoiceFunction
	// 工具类型
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Type ToolType
}

// MarshalJSON 序列化JSON
func (c ChatToolChoice) MarshalJSON() (b []byte, err error) {
	if c.ToolChoiceType != "" {
		return json.Marshal(c.ToolChoiceType)
	}
	return json.Marshal(struct {
		Function *ChatToolChoiceFunction `json:"function,omitempty"`
		Type     ToolType                `json:"type,omitempty"`
	}{
		Function: c.Function,
		Type:     c.Type,
	})
}

// ChatToolFunction 工具函数
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatToolFunction struct {
	// 函数名称，必须是 a-z, A-Z, 0-9 或者包含下划线和破折号，最大长度为 64
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Name string `json:"name,omitempty"`
	// 函数描述，用于帮助模型决定何时以及如何调用函数
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Description string `json:"description,omitempty"`
	// 函数接受的参数，描述为一个 JSON Schema 对象
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Parameters map[string]any `json:"parameters,omitempty"`
	// 是否启用严格模式，默认为 false
	// 提供商支持: OpenAI
	Strict bool `json:"strict,omitempty"`
}

// ChatTool 工具
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatTool struct {
	// 工具类型
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Type ToolType `json:"type,omitempty"`
	// 工具函数
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Function *ChatToolFunction `json:"function,omitempty"`
}

// ChatSearchContextSize 搜索上下文大小
//
//	提供商支持: OpenAI
type ChatSearchContextSize string

const (
	// 提供商支持: OpenAI
	ChatSearchContextSizeLow ChatSearchContextSize = "low"
	// 提供商支持: OpenAI
	ChatSearchContextSizeMedium ChatSearchContextSize = "medium"
	// 提供商支持: OpenAI
	ChatSearchContextSizeHigh ChatSearchContextSize = "high"
)

// ChatApproximateLocation 用户的大致位置参数
//
//	提供商支持: OpenAI
type ChatApproximateLocation struct {
	// 用户所在城市
	// 提供商支持: OpenAI
	City string `json:"city,omitempty"`
	// 用户所在国家的两字母 ISO 代码
	// 提供商支持: OpenAI
	Country string `json:"country,omitempty"`
	// 用户所在地区
	// 提供商支持: OpenAI
	Region string `json:"region,omitempty"`
	// 用户的 IANA 时区
	// 提供商支持: OpenAI
	Timezone string `json:"timezone,omitempty"`
}

// ChatApproximateLocationType 位置近似类型
//
//	提供商支持: OpenAI
type ChatApproximateLocationType string

const (
	// 提供商支持: OpenAI
	ChatApproximateLocationTypeApproximate ChatApproximateLocationType = "approximate"
)

// ChatUserLocation 用户位置信息
//
//	提供商支持: OpenAI
type ChatUserLocation struct {
	// 大致位置信息
	// 提供商支持: OpenAI
	Approximate *ChatApproximateLocation `json:"approximate,omitempty"`
	// 位置近似类型
	// 提供商支持: OpenAI
	Type ChatApproximateLocationType `json:"type,omitempty"`
}

// ChatSearchStrategyType 搜索互联网信息的数量
//
//	提供商支持: AliBL
type ChatSearchStrategy string

const (
	// 在请求时搜索5条互联网信息
	// 提供商支持: AliBL
	ChatSearchStrategyTypeStandard ChatSearchStrategy = "standard"
	// 在请求时搜索10条互联网信息
	// 提供商支持: AliBL
	ChatSearchStrategyTypeDeep ChatSearchStrategy = "pro"
)

// ChatWebSearchOptions 网络搜索选项
//
//	提供商支持: OpenAI | AliBL
type ChatWebSearchOptions struct {
	// 搜索上下文大小
	// 提供商支持: OpenAI
	SearchContextSize ChatSearchContextSize `json:"search_context_size,omitempty"`
	// 用户位置信息
	// 提供商支持: OpenAI
	UserLocation *ChatUserLocation `json:"user_location,omitempty"`
	// 是否强制开启搜索
	// 提供商支持: AliBL
	ForcedSearch bool `json:"forced_search,omitempty"`
	// 搜索互联网信息的数量
	// 提供商支持: AliBL
	SearchStrategy ChatSearchStrategy `json:"search_strategy,omitempty"`
}

// ChatTranslationLanguageType 翻译支持的语言类型
//
//	提供商支持: AliBL
type ChatTranslationLanguageType string

const (
	// 自动检测语言
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeAuto ChatTranslationLanguageType = "auto"
	// 中文
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeChinese ChatTranslationLanguageType = "Chinese"
	// 英语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeEnglish ChatTranslationLanguageType = "English"
	// 日语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeJapanese ChatTranslationLanguageType = "Japanese"
	// 韩语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeKorean ChatTranslationLanguageType = "Korean"
	// 泰语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeThai ChatTranslationLanguageType = "Thai"
	// 法语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFrench ChatTranslationLanguageType = "French"
	// 德语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeGerman ChatTranslationLanguageType = "German"
	// 西班牙语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeSpanish ChatTranslationLanguageType = "Spanish"
	// 阿拉伯语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeArabic ChatTranslationLanguageType = "Arabic"
	// 印尼语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeIndonesian ChatTranslationLanguageType = "Indonesian"
	// 越南语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeVietnamese ChatTranslationLanguageType = "Vietnamese"
	// 巴西葡萄牙语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePortuguese ChatTranslationLanguageType = "Portuguese"
	// 意大利语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeItalian ChatTranslationLanguageType = "Italian"
	// 荷兰语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeDutch ChatTranslationLanguageType = "Dutch"
	// 俄语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeRussian ChatTranslationLanguageType = "Russian"
	// 高棉语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeKhmer ChatTranslationLanguageType = "Khmer"
	// 宿务语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCebuano ChatTranslationLanguageType = "Cebuano"
	// 菲律宾语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFilipino ChatTranslationLanguageType = "Filipino"
	// 捷克语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCzech ChatTranslationLanguageType = "Czech"
	// 波兰语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePolish ChatTranslationLanguageType = "Polish"
	// 波斯语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePersian ChatTranslationLanguageType = "Persian"
	// 希伯来语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHebrew ChatTranslationLanguageType = "Hebrew"
	// 土耳其语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeTurkish ChatTranslationLanguageType = "Turkish"
	// 印地语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHindi ChatTranslationLanguageType = "Hindi"
	// 孟加拉语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeBengali ChatTranslationLanguageType = "Bengali"
	// 乌尔都语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeUrdu ChatTranslationLanguageType = "Urdu"
	// 马来语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeMalay ChatTranslationLanguageType = "Malay"
	// 匈牙利语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHungarian ChatTranslationLanguageType = "Hungarian"
	// 瑞典语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeSwedish ChatTranslationLanguageType = "Swedish"
	// 芬兰语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFinnish ChatTranslationLanguageType = "Finnish"
	// 缅甸语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeBurmese ChatTranslationLanguageType = "Burmese"
	// 老挝语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeLao ChatTranslationLanguageType = "Lao"
	// 粤语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCantonese ChatTranslationLanguageType = "Cantonese"
	// 罗马尼亚语
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeRomanian ChatTranslationLanguageType = "Romanian"
)

// ChatTranslationTerm 翻译术语
//
//	提供商支持: AliBL
type ChatTranslationTerm struct {
	// 源语言的术语
	// 提供商支持: AliBL
	Source string `json:"source,omitempty"`
	// 目标语言的术语
	// 提供商支持: AliBL
	Target string `json:"target,omitempty"`
}

// ChatTranslationMemory 翻译记忆
//
//	提供商支持: AliBL
type ChatTranslationMemory struct {
	// 源语言的语句
	// 提供商支持: AliBL
	Source string `json:"source,omitempty"`
	// 目标语言的语句
	// 提供商支持: AliBL
	Target string `json:"target,omitempty"`
}

// ChatTranslationOptions 翻译选项
//
//	提供商支持: AliBL
type ChatTranslationOptions struct {
	// 源语言的英文全称，可以将source_lang设置为"auto"，模型会自动判断输入文本属于哪种语言
	// 提供商支持: AliBL
	SourceLang ChatTranslationLanguageType `json:"source_lang,omitempty"`
	// 目标语言的英文全称
	// 提供商支持: AliBL
	TargetLang ChatTranslationLanguageType `json:"target_lang,omitempty"`
	// 在使用术语干预翻译功能时需要设置的术语数组
	// 提供商支持: AliBL
	Terms []ChatTranslationTerm `json:"terms,omitempty"`
	// 在使用翻译记忆功能时需要设置的翻译记忆数组
	// 提供商支持: AliBL
	TmList []ChatTranslationMemory `json:"tm_list,omitempty"`
	// 在使用领域提示功能时需要设置的领域提示语句
	// 提供商支持: AliBL
	Domains string `json:"domains,omitempty"`
}

type ChatThinkingOptions struct {
	Type string `json:"type,omitempty"`
}

// ChatRequest 聊天请求
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ChatRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 消息数组
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Messages []ChatMessage `json:"messages,omitempty"`
	// 模型名称
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Model string `json:"model,omitempty"`
	// 输出音频的音色与格式
	// 提供商支持: OpenAI | AliBL
	Audio *ChatAudioOutputArgs `json:"audio,omitempty"`
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚，降低模型重复相同内容的可能性
	// 提供商支持: OpenAI | DeepSeek | Ark
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
	// 修改指定标记在补全中出现的可能性
	// 提供商支持: OpenAI | Ark
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// 是否返回输出 Token 的对数概率
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	LogProbs bool `json:"logprobs,omitempty"`
	// 生成补全内容的最大令牌数上限
	// 提供商支持: OpenAI | DeepSeek | AliBL
	MaxCompletionTokens int `json:"max_completion_tokens,omitempty"`
	// 元数据
	// 提供商支持: OpenAI
	Metadata map[string]string `json:"metadata,omitempty"`
	// 输出数据的模态
	// 提供商支持: OpenAI | AliBL
	Modalities []ChatModalitiesType `json:"modalities,omitempty"`
	// 生成响应的个数
	// 提供商支持: OpenAI | AliBL
	N int `json:"n,omitempty"`
	// 是否开启并行工具调用
	// 提供商支持: OpenAI | AliBL
	ParallelToolCalls bool `json:"parallel_tool_calls,omitempty"`
	// 预测输出配置
	// 提供商支持: OpenAI
	Prediction *ChatPrediction `json:"prediction,omitempty"`
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	PresencePenalty float32 `json:"presence_penalty,omitempty"`
	// 仅适用于 o 系列模型，约束推理模型的推理努力程度
	// 提供商支持: OpenAI
	ReasoningEffort ChatReasoningEffortType `json:"reasoning_effort,omitempty"`
	// 响应格式
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	ResponseFormat *ChatResponseFormat `json:"response_format,omitempty"`
	// 随机种子
	// 提供商支持: OpenAI | AliBL
	Seed int `json:"seed,omitempty"`
	// 指定用于处理请求的延迟层级。此参数与订阅了规模层级服务的客户相关
	// 提供商支持: OpenAI | Ark
	ServiceTier string `json:"service_tier,omitempty"`
	// 当API遇到这些序列时将停止生成更多标记。返回的文本不会包含停止序列
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Stop []string `json:"stop,omitempty"`
	// 是否存储此聊天完成请求的输出，用于我们的模型蒸馏或评估产品
	// 提供商支持: OpenAI
	Store bool `json:"store,omitempty"`
	// 是否流式传输响应
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Stream bool `json:"stream,omitempty"`
	// 流式传输选项
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	StreamOptions *ChatStreamOptions `json:"stream_options,omitempty"`
	// 采样温度值，范围在0到2之间
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Temperature float32 `json:"temperature,omitempty"`
	// 指定工具调用的策略
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolChoice *ChatToolChoice `json:"tool_choice,omitempty"`
	// 可供模型调用的工具数组
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Tools []ChatTool `json:"tools,omitempty"`
	// 一个介于0和20之间的整数，指定在每个标记位置返回的最可能标记的数量，每个标记都有相关的对数概率。如果使用此参数，必须将logprobs设置为true
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	TopLogProbs int `json:"top_logprobs,omitempty"`
	// 一种替代温度采样的方法，我们通常建议调整此参数或温度（temperature），但不要同时调整两者
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	TopP float32 `json:"top_p,omitempty"`
	// 网络搜索选项
	// 提供商支持: OpenAI | AliBL
	WebSearchOptions *ChatWebSearchOptions `json:"web_search_options,omitempty"`
	// 生成过程中采样候选集的大小
	// 提供商支持: AliBL
	TopK int `json:"top_k,omitempty"`
	// 是否开启思考模式
	// 提供商支持: AliBL
	EnableThinking bool `json:"enable_thinking,omitempty"`
	// 思考过程的最大长度，在enable_thinking为true时生效
	// 提供商支持: AliBL
	ThinkingBudget int `json:"thinking_budget,omitempty"`
	// 翻译选项
	// 提供商支持: AliBL
	TranslationOptions *ChatTranslationOptions `json:"translation_options,omitempty"`
	// 在 API 的内容安全能力基础上，是否进一步识别输入输出内容的违规信息
	// 可选值：'{"input":"cip","output":"cip"}'，表示同时检查输入和输出
	// 提供商支持: AliBL
	XDashScopeDataInspection string `json:"-"`
	// 是否开启思考模式
	// 提供商支持: Ark
	Thinking *ChatThinkingOptions `json:"thinking,omitempty"`
}

// MarshalJSON 序列化JSON
func (r ChatRequest) MarshalJSON() (b []byte, err error) {
	strategy, ok := chatRequestStrategies[r.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", r.Provider)
	}
	return strategy(r)
}

// ChatFinishReason 模型停止生成 token 的原因
type ChatFinishReason string

const (
	ChatFinishReasonNull                       ChatFinishReason = "null"
	ChatFinishReasonStop                       ChatFinishReason = "stop"                         // 模型自然停止生成，或遇到 stop 序列中列出的字符串
	ChatFinishReasonLength                     ChatFinishReason = "length"                       // 输出长度达到了模型上下文长度限制，或达到了 max_completion_tokens 的限制
	ChatFinishReasonContentFilter              ChatFinishReason = "content_filter"               // 输出内容因触发过滤策略而被过滤
	ChatFinishReasonToolCalls                  ChatFinishReason = "tool_calls"                   // 模型调用了工具
	ChatFinishReasonInsufficientSystemResource ChatFinishReason = "insufficient_system_resource" // 系统推理资源不足，生成被打断
)

// MarshalJSON 序列化JSON
func (r ChatFinishReason) MarshalJSON() (b []byte, err error) {
	if r == ChatFinishReasonNull || r == "" {
		return []byte("null"), nil
	}
	return []byte(`"` + string(r) + `"`), nil
}

// ChatTopLogProbs 输出概率 top N 的 token 的列表，以及它们的对数概率
type ChatTopLogProbs struct {
	Bytes   []byte  `json:"bytes,omitempty"`   // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
	LogProb float64 `json:"logprob,omitempty"` // 该 token 的对数概率
	Token   string  `json:"token,omitempty"`   // 输出的 token
}

// ChatLogProb 对数概率信息
type ChatLogProb struct {
	Bytes       []byte            `json:"bytes,omitempty"`        // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
	LogProb     float64           `json:"logprob,omitempty"`      // 该 token 的对数概率
	Token       string            `json:"token,omitempty"`        // 输出的 token
	TopLogProbs []ChatTopLogProbs `json:"top_logprobs,omitempty"` // 一个包含在该输出位置上，输出概率 top N 的 token 的列表，以及它们的对数概率
}

// ChatLogProbs 对数概率信息
type ChatLogProbs struct {
	Content []ChatLogProb `json:"content,omitempty"` // 一个包含输出 token 对数概率信息的列表
	Refusal []ChatLogProb `json:"refusal,omitempty"` // 一个带有对数概率信息的消息拒绝 token 列表
}

// ChatAnnotationType 消息的注释类型
type ChatAnnotationType string

const (
	ChatAnnotationTypeURLCitation ChatAnnotationType = "url_citation" // 网络搜索工具的 URL 引用
)

// ChatAnnotationURLCitation 网络搜索工具的 URL 引用
type ChatAnnotationURLCitation struct {
	EndIndex   int    `json:"end_index,omitempty"`   // 消息中URL引用的最后一个字符的索引
	StartIndex int    `json:"start_index,omitempty"` // 消息中URL引用的第一个字符的索引
	Title      string `json:"title,omitempty"`       // 网络资源的标题
	URL        string `json:"url,omitempty"`         // 网络资源的URL
}

// ChatAnnotation 消息的注释
type ChatAnnotation struct {
	Type        ChatAnnotationType         `json:"type,omitempty"`         // 注释的类型
	URLCitation *ChatAnnotationURLCitation `json:"url_citation,omitempty"` // 网络搜索工具的 URL 引用
}

// ChatAudioOutput 音频响应数据
type ChatAudioOutput struct {
	Data       string `json:"data,omitempty"`       // Base64 编码的音频字节，格式在请求中指定
	ExpiresAt  int64  `json:"expires_at,omitempty"` // 音频响应在服务器上不再可访问的 Unix 时间戳（秒）
	ID         string `json:"id,omitempty"`         // 音频响应的唯一标识符
	Transcript string `json:"transcript,omitempty"` // 模型生成的音频文本转录
}

// ToolCallsFunction 函数
type ToolCallsFunction struct {
	Arguments string `json:"arguments,omitempty"` // 函数参数
	Name      string `json:"name,omitempty"`      // 函数名
}

// ToolCalls 工具调用
type ToolCalls struct {
	Index    int                `json:"index,omitempty"`    // 索引
	Function *ToolCallsFunction `json:"function,omitempty"` // 函数调用
	ID       string             `json:"id,omitempty"`       // 工具ID
	Type     ToolType           `json:"type,omitempty"`     // 工具类型
}

// ChatCompletionMessage 模型生成的 completion 消息
type ChatCompletionMessage struct {
	Content          string           `json:"content,omitempty"`           // 文本内容
	ReasoningContent string           `json:"reasoning_content,omitempty"` // 推理内容
	Refusal          string           `json:"refusal,omitempty"`           // 拒绝消息
	Role             string           `json:"role,omitempty"`              // 角色
	Annotations      []ChatAnnotation `json:"annotations,omitempty"`       // 消息的注释，在适用情况下提供，例如使用网络搜索工具时
	Audio            *ChatAudioOutput `json:"audio,omitempty"`             // 音频响应数据
	ToolCalls        []ToolCalls      `json:"tool_calls,omitempty"`        // 工具调用
}

// ChatChoice 模型生成的 completion
type ChatChoice struct {
	FinishReason ChatFinishReason       `json:"finish_reason,omitempty"` // 模型停止生成 token 的原因
	Index        int                    `json:"index,omitempty"`         // 该 completion 在模型生成的 completion 的选择列表中的索引
	LogProbs     *ChatLogProbs          `json:"logprobs,omitempty"`      // 该 choice 的对数概率信息
	Message      *ChatCompletionMessage `json:"message,omitempty"`       // 模型生成的 completion 消息
	Delta        *ChatCompletionMessage `json:"delta,omitempty"`         // 流式传输的增量信息
}

// CompletionTokensDetails completion tokens 的详细信息
type CompletionTokensDetails struct {
	TextTokens               int `json:"text_tokens,omitempty"`                // 模型生成的文本 token 数量
	AudioTokens              int `json:"audio_tokens,omitempty"`               // 模型生成的音频输入 token 数量
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`           // 模型用于推理而生成的 token 数量
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"` // 使用预测输出时，预测中出现在完成结果中的 token 数量
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"` // 使用预测输出时，预测中未出现在完成结果中的 token 数量
}

// PromptTokensDetails prompt tokens 的详细信息
type PromptTokensDetails struct {
	AudioTokens  int `json:"audio_tokens,omitempty"`  // prompt 中存在的音频输入 token 数
	CachedTokens int `json:"cached_tokens,omitempty"` // prompt 中存在的缓存 token 数
	TextTokens   int `json:"text_tokens,omitempty"`   // prompt 中存在的文本 token 数
	ImageTokens  int `json:"image_tokens,omitempty"`  // prompt 中存在的图像 token 数
	VideoTokens  int `json:"video_tokens,omitempty"`  // prompt 中存在的视频 token 数
}

// ChatUsage 该对话补全请求的用量信息
type ChatUsage struct {
	CompletionTokens        int                      `json:"completion_tokens,omitempty"`         // 模型 completion 产生的 token 数
	PromptTokens            int                      `json:"prompt_tokens,omitempty"`             // 用户 prompt 所包含的 token 数
	PromptCacheHitTokens    int                      `json:"prompt_cache_hit_tokens,omitempty"`   // 用户 prompt 中，命中上下文缓存的 token 数
	PromptCacheMissTokens   int                      `json:"prompt_cache_miss_tokens,omitempty"`  // 用户 prompt 中，未命中上下文缓存的 token 数
	TotalTokens             int                      `json:"total_tokens,omitempty"`              // 该请求中，所有 token 的数量（prompt + completion）
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // completion tokens 的详细信息
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // prompt tokens 的详细信息
}

// ChatBaseResponse 聊天响应基础信息
type ChatBaseResponse struct {
	Choices           []ChatChoice            `json:"choices,omitempty"`            // 模型生成的 completion 的选择列表
	Created           int64                   `json:"created,omitempty"`            // 创建聊天完成时的 Unix 时间戳（以秒为单位）
	ID                string                  `json:"id,omitempty"`                 // 该对话的唯一标识符
	Model             string                  `json:"model,omitempty"`              // 生成该 completion 的模型名
	Object            string                  `json:"object,omitempty"`             // 对象的类型
	ServiceTier       string                  `json:"service_tier,omitempty"`       // 用于处理请求的服务层级
	SystemFingerprint string                  `json:"system_fingerprint,omitempty"` // 此指纹表示模型运行的后端配置。可以与 seed 请求参数一起使用，以了解何时进行了可能影响确定性的后端更改
	Usage             *ChatUsage              `json:"usage,omitempty"`              // 该对话补全请求的用量信息
	StreamStats       *httpclient.StreamStats `json:"stream_stats,omitempty"`       // 流式传输统计信息
}

// SetStreamStats 设置流式传输统计信息
func (c *ChatBaseResponse) SetStreamStats(stats httpclient.StreamStats) {
	c.StreamStats = &stats
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ChatBaseResponse
	httpclient.HttpHeader
}

// ChatResponseStream 流式传输的聊天响应
type ChatResponseStream struct {
	*httpclient.StreamReader[ChatBaseResponse]
}
