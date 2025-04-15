/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 20:54:21
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 11:01:37
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import "net/http"

const (
	defaultEmptyMessagesLimit uint = 300 // 默认空消息限制
)

// HTTPDoer HTTP 请求执行器接口
type HTTPDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// ClientConfig 客户端配置
type ClientConfig struct {
	BaseURL            string   // API 的基础 URL 地址
	AuthToken          string   // 认证令牌，用于 API 访问授权
	HTTPClient         HTTPDoer // HTTP 客户端实现，用于发送请求
	EmptyMessagesLimit uint     // 空消息限制
}

// DefaultConfig 默认客户端配置
func DefaultConfig(baseURL, authToken string) (config ClientConfig) {
	return ClientConfig{
		BaseURL:            baseURL,
		AuthToken:          authToken,
		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
