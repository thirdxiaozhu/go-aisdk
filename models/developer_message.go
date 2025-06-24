/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 18:25:47
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
	// 序列化开发者消息函数（OpenAI）
	marshalDeveloperMessageByOpenAI = func(m DeveloperMessage) (b []byte, err error) {
		type Alias DeveloperMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "developer",
			Alias: Alias(m),
		}
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 策略映射
	developerMessageStrategies = map[consts.Provider]func(m DeveloperMessage) (b []byte, err error){
		consts.OpenAI: marshalDeveloperMessageByOpenAI,
	}
)

// DeveloperMessage 开发者消息
//
//	提供商支持: OpenAI
type DeveloperMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 文本内容
	// 提供商支持: OpenAI
	Content string `json:"content,omitempty"`
	// 参与者名称
	// 提供商支持: OpenAI
	Name string `json:"name,omitempty"`
}

// GetRole 获取消息角色
func (m DeveloperMessage) GetRole() (role string) { return "developer" }

// SetProvider 设置提供商
func (m *DeveloperMessage) SetProvider(provider consts.Provider) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m DeveloperMessage) MarshalJSON() (b []byte, err error) {
	strategy, ok := developerMessageStrategies[m.provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}
	return strategy(m)
}
