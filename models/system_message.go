/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 18:26:07
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
	// 序列化系统消息函数（Aliyunbl）
	marshalSystemMessageByAliyunbl = func(m SystemMessage) (b []byte, err error) {
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
	// 策略映射
	systemMessageStrategies = map[consts.Provider]func(m SystemMessage) (b []byte, err error){
		consts.OpenAI:   marshalSystemMessageByOpenAI,
		consts.DeepSeek: marshalSystemMessageByOpenAI,
		consts.Aliyunbl: marshalSystemMessageByAliyunbl,
	}
)

// SystemMessage 系统消息
//
//	提供商支持: OpenAI | DeepSeek | Aliyunbl
type SystemMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | Aliyunbl
	Content string `json:"content,omitempty"`
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
