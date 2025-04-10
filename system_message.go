/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 17:07:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-09 17:07:04
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import "encoding/json"

// SystemMessage 系统消息
type SystemMessage struct {
	Content string `json:"content"`        // 文本内容
	Name    string `json:"name,omitempty"` // 参与者名称
}

// GetRole 获取消息角色
func (m SystemMessage) GetRole() (role string) { return "system" }

// MarshalJSON 序列化JSON
func (m SystemMessage) MarshalJSON() (b []byte, err error) {
	type Alias SystemMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		Alias
	}{
		Role:  "system",
		Alias: Alias(m),
	})
}
