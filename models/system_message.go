/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 22:51:21
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

// SystemMessage 系统消息
type SystemMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Content string `json:"content,omitempty" providers:"openai,deepseek,alibl"`
	// 消息角色
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Role string `json:"role,omitempty" providers:"openai,deepseek,alibl" default:"system"`
	// 参与者名称
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty" providers:"openai,deepseek"`
}

// SetProvider 设置提供商
func (m *SystemMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m SystemMessage) MarshalJSON() (b []byte, err error) {
	return NewSerializer(m.provider).Serialize(m)
}
