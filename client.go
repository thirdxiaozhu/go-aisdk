/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-05 19:12:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/flake"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/httpclient/middleware"
	"github.com/liusuxian/go-aisdk/models"
	_ "github.com/liusuxian/go-aisdk/providers"
	"github.com/liusuxian/go-aisdk/sdkerrors"
	"sort"
)

// SDKClient SDK客户端
type SDKClient struct {
	configManager   *conf.SDKConfigManager // 配置管理器
	flakeInstance   *flake.Flake           // 分布式唯一ID生成器
	middlewareChain *middleware.Chain      // 中间件链
	noCheckMethods  map[string]bool        // 不需要检查模型支持的方法
}

// NewSDKClient 创建一个SDK客户端
func NewSDKClient(configPath string, opts ...SDKClientOption) (client *SDKClient, err error) {
	// 创建SDK配置管理器
	var configManager *conf.SDKConfigManager
	if configManager, err = conf.NewSDKConfigManager(configPath); err != nil {
		return
	}
	// 创建一个分布式唯一ID生成器
	var flakeInstance *flake.Flake
	if flakeInstance, err = flake.New(flake.Settings{}); err != nil {
		return
	}
	// 初始化所有提供商
	for _, provider := range core.ListProviders() {
		// 获取提供商
		var ps core.ProviderService
		if ps, err = core.GetProvider(provider); err != nil {
			return
		}
		// 获取提供商配置
		providerConfig := configManager.GetProviderConfig(provider)
		// 初始化提供商配置
		ps.InitializeProviderConfig(&providerConfig)
	}
	// 处理选项
	cliOpt := &clientOption{}
	for _, opt := range opts {
		opt(cliOpt)
	}
	// 按优先级排序中间件
	sort.Slice(cliOpt.middlewares, func(i, j int) bool {
		return cliOpt.middlewares[i].Priority() < cliOpt.middlewares[j].Priority()
	})
	// 创建中间件链
	middlewareChain := middleware.NewChain(cliOpt.middlewares...)
	// 创建SDK客户端
	client = &SDKClient{
		configManager:   configManager,
		flakeInstance:   flakeInstance,
		middlewareChain: middlewareChain,
		noCheckMethods: map[string]bool{
			"ListModels": true,
		},
	}
	return
}

// ListModels 列出模型
func (c *SDKClient) ListModels(ctx context.Context, userId string, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		// 列出模型
		return ps.ListModels(ctx, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, userId, models.ModelInfo{Provider: provider}, "ListModels", nil, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ListModelsResponse)
	return
}

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, userId string, request models.Request, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		// 判断是否流式传输
		if chatReq.Stream {
			return nil, sdkerrors.ErrCompletionStreamNotSupported
		}

		if err = ps.CheckRequestValidation(request); err != nil {
			return nil, err
		}

		// 创建聊天
		return ps.CreateChatCompletion(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, userId, request.GetModelInfo(), "CreateChatCompletion", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ChatResponse)
	return
}

// CreateChatCompletionStream  创建聊天
func (c *SDKClient) CreateChatCompletionStream(ctx context.Context, userId string, request models.Request, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (interface{}, error) {
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
	if _, err = c.handlerRequest(ctx, userId, request.GetModelInfo(), "CreateChatCompletionStream", request, handler); err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *SDKClient) CreateImageGeneration(ctx context.Context, userId string, request models.Request, opts ...httpclient.HTTPClientOption) (response httpclient.Response, err error) {
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
	if resp, err = c.handlerRequest(ctx, userId, request.GetModelInfo(), "CreateImageGeneration", request, handler); err != nil {
		return nil, err
	}
	response = resp.(httpclient.Response)
	return
}

type MultipleRoundDialogue struct {
}
