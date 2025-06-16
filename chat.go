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
	"github.com/liusuxian/go-aisdk/sdkerror"
)

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		// 判断是否流式传输
		if chatReq.Stream {
			return nil, sdkerror.ErrCompletionStreamNotSupported
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

// CreateChatCompletionStream  创建聊天
func (c *SDKClient) CreateChatCompletionStream(ctx context.Context, request models.Request, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (httpclient.Response, error) {
	var err error

	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		// 创建聊天
		if err = ps.CheckRequestValidation(request); err != nil {
			return nil, err
		}
		return ps.CreateChatCompletionStream(ctx, chatReq, cb, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.GetModelInfo(), "CreateChatCompletionStream", request.GetUserInfo().UserID, request, handler); err != nil {
		return nil, err
	}
	return resp.(httpclient.Response), nil
}

func (c *SDKClient) CreateImageGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (response httpclient.Response, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.Request)
		// 创建聊天
		if err = ps.CheckRequestValidation(request); err != nil {
			return nil, err
		}
		return ps.CreateImageGeneration(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.GetModelInfo(), "CreateImageGeneration", request.GetUserInfo().UserID, request, handler); err != nil {
		return nil, err
	}
	response = resp.(httpclient.Response)
	return
}

func (c *SDKClient) CreateVideoGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (response httpclient.Response, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.Request)
		// 创建聊天
		if err = ps.CheckRequestValidation(request); err != nil {
			return nil, err
		}
		return ps.CreateVideoGeneration(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.GetModelInfo(), "CreateVideoGeneration", request.GetUserInfo().UserID, request, handler); err != nil {
		return nil, err
	}
	response = resp.(httpclient.Response)
	return
}
