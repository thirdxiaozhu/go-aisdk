/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 16:42:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 16:28:24
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/liusuxian/go-aisdk/httpclient"

// ListModelsRequest 列出模型请求
type ListModelsRequest struct {
	ModelInfo
	UserInfo
}

// ListModelsResponse 列出模型响应
type ListModelsResponse struct {
	Object string   `json:"object"` // 对象类型
	Data   []Models `json:"data"`   // 模型列表
	httpclient.HttpHeader
}

// Models 模型信息
type Models struct {
	Created int64  `json:"created,omitempty"` // 创建时间
	ID      string `json:"id"`                // 模型的标识符
	Object  string `json:"object"`            // 对象的类型，其值为 model
	OwnedBy string `json:"owned_by"`          // 拥有该模型的组织
}

type Request interface {
	GetModelInfo() ModelInfo
	MarshalJSON() (b []byte, err error)
}
