/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 17:39:52
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-10 13:35:40
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import "encoding/json"

// AssistantMessage 助手消息
type AssistantMessage struct {
	Audio             *ChatAssistantMsgAudio      `json:"audio,omitempty"`              // 音频
	Content           string                      `json:"content,omitempty"`            // 文本内容
	MultimodalContent []ChatAssistantMsgPart      `json:"multimodal_content,omitempty"` // 多模态内容
	Name              string                      `json:"name,omitempty"`               // 参与者名称
	Refusal           string                      `json:"refusal,omitempty"`            // 拒绝消息
	ToolCalls         []ChatAssistantMsgToolCalls `json:"tool_calls,omitempty"`         // 工具调用
	Prefix            bool                        `json:"prefix,omitempty"`             // (DeepSeek)设置此参数为 true，来强制模型在其回答中以此 assistant 消息中提供的前缀内容开始
	ReasoningContent  string                      `json:"reasoning_content,omitempty"`  // (DeepSeek)用于 deepseek-reasoner 模型在对话前缀续写功能下，作为最后一条 assistant 思维链内容的输入。使用此功能时，prefix 参数必须设置为 true
}

// GetRole 获取消息角色
func (m AssistantMessage) GetRole() (role string) { return "assistant" }

// MarshalJSON 序列化JSON
func (m AssistantMessage) MarshalJSON() (b []byte, err error) {
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
	return json.Marshal(temp)
}

// ChatAssistantMsgAudio 音频
type ChatAssistantMsgAudio struct {
	ID string `json:"id"` // 音频ID
}

// ChatAssistantMsgPartType 多模态内容类型
type ChatAssistantMsgPartType string

const (
	ChatAssistantMsgPartTypeText    ChatAssistantMsgPartType = "text"
	ChatAssistantMsgPartTypeRefusal ChatAssistantMsgPartType = "refusal"
)

// ChatAssistantMsgPart 多模态内容
type ChatAssistantMsgPart struct {
	Type    ChatAssistantMsgPartType `json:"type"`              // 内容类型
	Text    string                   `json:"text,omitempty"`    // 文本内容
	Refusal string                   `json:"refusal,omitempty"` // 拒绝消息
}

// ChatAssistantMsgToolFunction 函数
type ChatAssistantMsgToolFunction struct {
	Arguments string `json:"arguments"` // 函数参数
	Name      string `json:"name"`      // 函数名
}

// ChatAssistantMsgToolCalls 工具调用
type ChatAssistantMsgToolCalls struct {
	Function *ChatAssistantMsgToolFunction `json:"function"` // 函数调用
	ID       string                        `json:"id"`       // 工具ID
	Type     ToolType                      `json:"type"`     // 工具类型
}
