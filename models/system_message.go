/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:16:45
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/liusuxian/go-aisdk/internal/utils"

type ChatSystemMsgPartType string

const (
	// 提供商支持: Ark
	ChatSystemMsgPartTypeText ChatSystemMsgPartType = "text"
)

// ChatSystemMsgImageURLDetail 图像质量
type ChatSystemMsgImageURLDetail string

const (
	// 提供商支持: Ark
	ChatSystemMsgImageURLDetailHigh ChatSystemMsgImageURLDetail = "high"
	// 提供商支持: Ark
	ChatSystemMsgImageURLDetailLow ChatSystemMsgImageURLDetail = "low"
	// 提供商支持: Ark
	ChatSystemMsgImageURLDetailAuto ChatSystemMsgImageURLDetail = "auto"
)

// ChatSystemMsgImageURL 图像URL
type ChatSystemMsgImageURL struct {
	// 图像URL，支持url和base64编码
	// 提供商支持: Ark
	URL string `json:"url,omitempty" providers:"ark" mapping:"alibl:image"`
	// 图像质量
	// 提供商支持: Ark
	Detail ChatSystemMsgImageURLDetail `json:"detail,omitempty" providers:"ark"`
}

// ChatSystemMsgPart 多模态内容
type ChatSystemMsgPart struct {
	// 内容类型
	// 提供商支持: Ark
	Type ChatSystemMsgPartType `json:"type,omitempty" providers:"ark"`
	// 文本内容
	// 提供商支持: Ark
	Text string `json:"text,omitempty" providers:"ark"`
	// 图像URL
	// 提供商支持: Ark
	ImageURL *ChatSystemMsgImageURL `json:"image_url,omitempty" providers:"ark"`
}

// SystemMessage 系统消息
type SystemMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 文本内容
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Content string `json:"content,omitempty" providers:"openai,deepseek,alibl,ark"`
	// 多模态内容
	// 提供商支持: Ark
	MultimodalContent []ChatSystemMsgPart `json:"multimodal_content,omitempty" providers:"ark" copyto:"Content"`
	// 消息角色
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Role string `json:"role,omitempty" providers:"openai,deepseek,alibl,ark" default:"system"`
	// 参与者名称
	//
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty" providers:"openai,deepseek"`
}

// SetProvider 设置提供商
func (m *SystemMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m SystemMessage) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(m.provider).Serialize(m)
}
