/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 10:57:18
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
	// 序列化系统消息函数（OpenAI）
	marshalSystemMessageByOpenAI = func(m SystemMessage) (b []byte, err error) {
		type Alias SystemMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "system",
			Alias: Alias(m),
		}
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化系统消息函数（AliBL）
	marshalSystemMessageByAliBL = func(m SystemMessage) (b []byte, err error) {
		type Alias SystemMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "system",
			Alias: Alias(m),
		}
		// 移除不支持的字段
		temp.Name = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	marshalSystemMessageByArk = func(m SystemMessage) (b []byte, err error) {
		type Alias SystemMessage
		temp := struct {
			Role    string `json:"role"`
			Content any    `json:"content"`
			Alias
		}{
			Role:  "system",
			Alias: Alias(m),
		}

		// 根据内容类型设置 content 字段
		if len(m.MultimodalContent) > 0 {
			// 移除不支持的字段
			tempMultimodalContent := make([]ChatUserMsgPart, 0, len(m.MultimodalContent))
			for _, v := range m.MultimodalContent {
				if v.File != nil {
					continue
				}
				if v.ImageURL != nil {
					v.ImageURL.Detail = ""
				}
				tempMultimodalContent = append(tempMultimodalContent, v)
			}
			temp.Content = tempMultimodalContent
		} else {
			temp.Content = m.Content
		}
		// 移除不支持的字段
		temp.Name = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 策略映射
	systemMessageStrategies = map[consts.Provider]func(m SystemMessage) (b []byte, err error){
		consts.OpenAI:   marshalSystemMessageByOpenAI,
		consts.DeepSeek: marshalSystemMessageByOpenAI,
		consts.AliBL:    marshalSystemMessageByAliBL,
	}
)

// ChatSystemMsgPartType 多模态内容类型
//
//	提供商支持: Ark
type ChatSystemMsgPartType string

const (
	//	提供商支持: Ark
	ChatSystemMsgPartTypeText     ChatSystemMsgPartType = "text"
	ChatSystemMsgPartTypeImageURL ChatSystemMsgPartType = "image_url"
)

// ChatSystemMsgImageURLDetail 图像质量
//
//	提供商支持: Ark
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
//
//	提供商支持: OpenAI | AliBL
type ChatSystemMsgImageURL struct {
	// 图像URL，支持url和base64编码
	// 提供商支持: Ark
	URL string `json:"url,omitempty"`
	// 图像质量
	// 提供商支持: Ark
	Detail ChatSystemMsgImageURLDetail `json:"detail,omitempty"`
}

// ChatSystemMsgPart 多模态内容
//
//	提供商支持: Ark
type ChatSystemMsgPart struct {
	// 内容类型
	// 提供商支持: Ark
	Type ChatSystemMsgPartType `json:"type,omitempty"`
	// 文本内容
	// 提供商支持: Ark
	Text string `json:"text,omitempty"`
	// 图像URL
	// 提供商支持: OpenAI | AliBL
	ImageURL *ChatSystemMsgImageURL `json:"image_url,omitempty"`
}

// SystemMessage 系统消息
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type SystemMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Content string `json:"content,omitempty"`
	// 多模态内容
	// 提供商支持: Ark
	MultimodalContent []ChatUserMsgPart `json:"-"`
	// 参与者名称
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty"`
}

// GetRole 获取消息角色
func (m SystemMessage) GetRole() (role string) { return "system" }

// SetProvider 设置提供商
func (m *SystemMessage) SetProvider(provider consts.Provider) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m SystemMessage) MarshalJSON() (b []byte, err error) {
	strategy, ok := systemMessageStrategies[m.provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}
	return strategy(m)
}
