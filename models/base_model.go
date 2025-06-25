/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 22:10:58
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/liusuxian/go-aisdk/consts"

// ModelInfo 模型信息结构
type ModelInfo struct {
	Provider  consts.Provider  `json:"provider,omitempty"`   // 提供商
	ModelType consts.ModelType `json:"model_type,omitempty"` // 模型类型
	Model     string           `json:"model,omitempty"`      // 模型名称
}

// UserInfo 用户信息结构
type UserInfo struct {
	UserID string `json:"user_id,omitempty"` // 用户ID
}
