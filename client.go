/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 00:28:24
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/flake"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	_ "github.com/liusuxian/go-aisdk/providers"
	"sort"
	"time"
)

// SDKClient SDK客户端
type SDKClient struct {
	configManager   *conf.SDKConfigManager // 配置管理器
	flakeInstance   *flake.Flake           // 分布式唯一ID生成器
	middlewareChain *httpclient.Chain      // 中间件链
	noCheckMethods  map[string]bool        // 不需要检查模型支持的方法
}

// SDKClientOption SDK客户端选项
type SDKClientOption func(c *clientOption)

// clientOption 客户端选项
type clientOption struct {
	middlewares []httpclient.Middleware
}

// NewSDKClient 创建一个SDK客户端
func NewSDKClient(configPath string, opts ...SDKClientOption) (client *SDKClient, err error) {
	// 创建SDK配置管理器
	var configManager *conf.SDKConfigManager
	if configManager, err = conf.NewSDKConfigManager(configPath); err != nil {
		err = wrapFailedToCreateConfigManager(err.Error())
		return
	}
	// 创建一个分布式唯一ID生成器
	var flakeInstance *flake.Flake
	if flakeInstance, err = flake.New(flake.Settings{}); err != nil {
		err = wrapFailedToCreateFlakeInstance(err.Error())
		return
	}
	// 初始化所有提供商
	for _, provider := range core.ListProviders() {
		// 获取提供商
		var ps core.ProviderService
		if ps = core.GetProvider(provider); ps == nil {
			err = wrapProviderNotSupported(provider)
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
	middlewareChain := httpclient.NewChain(cliOpt.middlewares...)
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
func (c *SDKClient) ListModels(ctx context.Context, request models.ListModelsRequest, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		// 列出模型
		return ps.ListModels(ctx, request.Provider, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider: request.Provider,
	}, request.UserInfo, "ListModels", nil, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ListModelsResponse)
	return
}

// handlerRequest 处理请求
func (c *SDKClient) handlerRequest(
	ctx context.Context,
	modelInfo models.ModelInfo,
	userInfo models.UserInfo,
	method string,
	request any,
	handler func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error),
) (resp any, err error) {
	// 生成唯一请求ID
	var requestId string
	if requestId, err = c.flakeInstance.RequestID(); err != nil {
		err = &SDKError{RequestID: requestId, Err: err}
		return
	}
	// 设置请求信息到上下文
	ctx = httpclient.SetRequestInfo(ctx, &httpclient.RequestInfo{
		Provider:  string(modelInfo.Provider),
		ModelType: string(modelInfo.ModelType),
		Model:     modelInfo.Model,
		Method:    method,
		StartTime: time.Now(),
		RequestID: requestId,
		User:      userInfo.User,
	})
	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req any) (resp any, err error) {
		// 获取提供商
		var ps core.ProviderService
		if ps = core.GetProvider(modelInfo.Provider); ps == nil {
			return nil, wrapProviderNotSupported(modelInfo.Provider)
		}
		// 根据方法名称决定是否需要判断模型支持
		var e error
		if !c.noCheckMethods[method] {
			// 判断模型是否支持
			if e = c.isModelSupported(ps, modelInfo); e != nil {
				return nil, e
			}
		}
		// 执行具体的处理逻辑
		return handler(ctx, ps, req)
	}
	// 执行中间件链
	if resp, err = c.middlewareChain.Execute(ctx, request, finalHandler); err != nil {
		err = &SDKError{RequestID: requestId, Err: err}
		return
	}
	return
}

// isModelSupported 判断模型是否支持
func (c *SDKClient) isModelSupported(s core.ProviderService, modelInfo models.ModelInfo) (err error) {
	// 获取支持的模型
	supportedModels := s.GetSupportedModels()
	if len(supportedModels) == 0 {
		return wrapProviderNotSupported(modelInfo.Provider)
	}
	// 获取指定模型类型支持的模型列表
	var (
		modelMap map[string]bool
		ok       bool
	)
	if modelMap, ok = supportedModels[modelInfo.ModelType]; !ok {
		return wrapModelTypeNotSupported(modelInfo.Provider, modelInfo.ModelType)
	}
	// 判断模型是否支持
	var modelSupported bool
	if modelSupported, ok = modelMap[modelInfo.Model]; !ok {
		return wrapModelNotSupported(modelInfo.Provider, modelInfo.Model, modelInfo.ModelType)
	}
	if !modelSupported {
		return wrapModelNotSupported(modelInfo.Provider, modelInfo.Model, modelInfo.ModelType)
	}
	return
}
