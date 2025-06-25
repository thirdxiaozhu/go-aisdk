/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 10:57:24
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
	// 序列化工具消息函数（OpenAI）
	marshalToolMessageByOpenAI = func(m ToolMessage) (b []byte, err error) {
		type Alias ToolMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "tool",
			Alias: Alias(m),
		}
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 策略映射
	toolMessageStrategies = map[consts.Provider]func(m ToolMessage) (b []byte, err error){
		consts.OpenAI:   marshalToolMessageByOpenAI,
		consts.DeepSeek: marshalToolMessageByOpenAI,
		consts.AliBL: marshalToolMessageByOpenAI,
	}
)

// ToolMessage 工具消息
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type ToolMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Content string `json:"content,omitempty"`
	// 工具调用ID
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// GetRole 获取消息角色
func (m ToolMessage) GetRole() (role string) { return "tool" }

// SetProvider 设置提供商
func (m *ToolMessage) SetProvider(provider consts.Provider) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m ToolMessage) MarshalJSON() (b []byte, err error) {
	strategy, ok := toolMessageStrategies[m.provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}
	return strategy(m)
}
