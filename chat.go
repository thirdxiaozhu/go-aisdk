/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-16 14:31:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 16:17:38
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
)

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		// 判断是否流式传输
		if chatReq.Stream {
			return nil, ErrCompletionStreamNotSupported
		}
		// 创建聊天
		return ps.CreateChatCompletion(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.ModelInfo, "CreateChatCompletion", request.UserInfo.UserID, request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ChatResponse)
	return
}
