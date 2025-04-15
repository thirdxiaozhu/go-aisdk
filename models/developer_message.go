/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 18:53:14
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "encoding/json"

// DeveloperMessage 开发者消息
type DeveloperMessage struct {
	Content string `json:"content"`        // 文本内容
	Name    string `json:"name,omitempty"` // 参与者名称
}

// GetRole 获取消息角色
func (m DeveloperMessage) GetRole() (role string) { return "developer" }

// MarshalJSON 序列化JSON
func (m DeveloperMessage) MarshalJSON() (b []byte, err error) {
	type Alias DeveloperMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		Alias
	}{
		Role:  "developer",
		Alias: Alias(m),
	})
}
