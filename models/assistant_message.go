/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 18:24:38
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
)

var (
	// 序列化助手消息函数（OpenAI）
	marshalAssistantMessageByOpenAI = func(m AssistantMessage) (b []byte, err error) {
		type Alias AssistantMessage
		temp := struct {
			Role    string `json:"role"`
			Content any    `json:"content"`
			Alias
		}{
			Role:  "assistant",
			Alias: Alias(m),
		}
		// 根据内容类型设置 content 字段
		if len(m.MultimodalContent) > 0 {
			temp.Content = m.MultimodalContent
		} else {
			temp.Content = m.Content
		}
		// 移除 multimodal_content 字段
		temp.MultimodalContent = nil
		// 移除不支持的字段
		temp.Prefix = false
		temp.ReasoningContent = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化助手消息函数（DeepSeek）
	marshalAssistantMessageByDeepSeek = func(m AssistantMessage) (b []byte, err error) {
		type Alias AssistantMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "assistant",
			Alias: Alias(m),
		}
		// 移除不支持的字段
		temp.Audio = nil
		temp.MultimodalContent = nil
		temp.Refusal = ""
		temp.ToolCalls = nil
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化助手消息函数（Aliyunbl）
	marshalAssistantMessageByAliyunbl = func(m AssistantMessage) (b []byte, err error) {
		type Alias AssistantMessage
		temp := struct {
			Role    string `json:"role"`
			Partial bool   `json:"partial,omitempty"`
			Alias
		}{
			Role:  "assistant",
			Alias: Alias(m),
		}
		// 设置前缀续写
		temp.Partial = m.Prefix
		// 移除不支持的字段
		temp.Audio = nil
		temp.MultimodalContent = nil
		temp.Name = ""
		temp.Refusal = ""
		temp.Prefix = false
		temp.ReasoningContent = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 策略映射
	assistantMessageStrategies = map[consts.Provider]func(m AssistantMessage) (b []byte, err error){
		consts.OpenAI:   marshalAssistantMessageByOpenAI,
		consts.DeepSeek: marshalAssistantMessageByDeepSeek,
		consts.Aliyunbl: marshalAssistantMessageByAliyunbl,
	}
)

// ChatAssistantMsgAudio 音频
//
//	提供商支持: OpenAI
type ChatAssistantMsgAudio struct {
	// 音频ID
	// 提供商支持: OpenAI
	ID string `json:"id,omitempty"`
}

// ChatAssistantMsgPartType 多模态内容类型
//
//	提供商支持: OpenAI
type ChatAssistantMsgPartType string

const (
	// 提供商支持: OpenAI
	ChatAssistantMsgPartTypeText ChatAssistantMsgPartType = "text"
	// 提供商支持: OpenAI
	ChatAssistantMsgPartTypeRefusal ChatAssistantMsgPartType = "refusal"
)

// ChatAssistantMsgPart 多模态内容
//
//	提供商支持: OpenAI
type ChatAssistantMsgPart struct {
	// 内容类型
	// 提供商支持: OpenAI
	Type ChatAssistantMsgPartType `json:"type,omitempty"`
	// 文本内容
	// 提供商支持: OpenAI
	Text string `json:"text,omitempty"`
	// 拒绝消息
	// 提供商支持: OpenAI
	Refusal string `json:"refusal,omitempty"`
}

// AssistantMessage 助手消息
//
//	提供商支持: OpenAI | DeepSeek | Aliyunbl
type AssistantMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 音频
	// 提供商支持: OpenAI
	Audio *ChatAssistantMsgAudio `json:"audio,omitempty"`
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | Aliyunbl
	Content string `json:"content,omitempty"`
	// 多模态内容
	// 提供商支持: OpenAI
	MultimodalContent []ChatAssistantMsgPart `json:"-"`
	// 参与者名称
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty"`
	// 拒绝消息
	// 提供商支持: OpenAI
	Refusal string `json:"refusal,omitempty"`
	// 工具调用
	// 提供商支持: OpenAI | Aliyunbl
	ToolCalls []ToolCalls `json:"tool_calls,omitempty"`
	// 设置此参数为 true，来强制模型在其回答中以此 assistant 消息中提供的前缀内容开始
	// 提供商支持: DeepSeek | Aliyunbl
	Prefix bool `json:"prefix,omitempty"`
	// 用于模型在对话前缀续写功能下，作为最后一条 assistant 思维链内容的输入。使用此功能时，prefix 参数必须设置为 true
	// 提供商支持: DeepSeek
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// GetRole 获取消息角色
func (m AssistantMessage) GetRole() (role string) { return "assistant" }

// SetProvider 设置提供商
func (m *AssistantMessage) SetProvider(provider consts.Provider) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m AssistantMessage) MarshalJSON() (b []byte, err error) {
	strategy, ok := assistantMessageStrategies[m.provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}
	return strategy(m)
}
