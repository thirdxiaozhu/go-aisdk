/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 19:08:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"crypto/rand"
	"fmt"
	"sort"
	"time"

	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/middleware"
	"github.com/liusuxian/go-aisdk/models"
	_ "github.com/liusuxian/go-aisdk/providers"
)

// SDKClient SDK客户端
type SDKClient struct {
	configManager   *conf.SDKConfigManager // 配置管理器
	middlewareChain *middleware.Chain      // 中间件链
}

// ClientOption 客户端选项
type ClientOption func(*clientOptions)

// clientOptions 客户端配置选项
type clientOptions struct {
	middlewares []middleware.Middleware
}

// WithMiddleware 添加中间件
func WithMiddleware(m middleware.Middleware) ClientOption {
	return func(opts *clientOptions) {
		opts.middlewares = append(opts.middlewares, m)
	}
}

// WithLogging 添加日志中间件
func WithLogging(config middleware.LoggingMiddlewareConfig) ClientOption {
	return func(opts *clientOptions) {
		opts.middlewares = append(opts.middlewares, middleware.NewLoggingMiddleware(config))
	}
}

// WithMetrics 添加监控中间件
func WithMetrics(config middleware.MetricsMiddlewareConfig) ClientOption {
	return func(opts *clientOptions) {
		opts.middlewares = append(opts.middlewares, middleware.NewMetricsMiddleware(config))
	}
}

// WithRetry 添加重试中间件
func WithRetry(config middleware.RetryMiddlewareConfig) ClientOption {
	return func(opts *clientOptions) {
		opts.middlewares = append(opts.middlewares, middleware.NewRetryMiddleware(config))
	}
}

// WithDefaultMiddlewares 添加默认中间件（日志、监控、重试）
func WithDefaultMiddlewares() ClientOption {
	return func(opts *clientOptions) {
		opts.middlewares = append(opts.middlewares,
			middleware.NewMetricsMiddleware(middleware.DefaultMetricsConfig()),
			middleware.NewRetryMiddleware(middleware.DefaultRetryConfig()),
			middleware.NewLoggingMiddleware(middleware.DefaultLoggingConfig()),
		)
	}
}

// NewSDKClient 创建一个SDK客户端
func NewSDKClient(configPath string, options ...ClientOption) (client *SDKClient, err error) {
	// 创建SDK配置管理器
	var configManager *conf.SDKConfigManager
	if configManager, err = conf.NewSDKConfigManager(configPath); err != nil {
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
	opts := &clientOptions{}
	for _, option := range options {
		option(opts)
	}

	// 按优先级排序中间件
	sort.Slice(opts.middlewares, func(i, j int) bool {
		return opts.middlewares[i].Priority() < opts.middlewares[j].Priority()
	})

	// 创建中间件链
	middlewareChain := middleware.NewChain(opts.middlewares...)

	// 创建SDK客户端
	client = &SDKClient{
		configManager:   configManager,
		middlewareChain: middlewareChain,
	}
	return
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// 如果加密随机数生成失败，使用时间戳作为fallback
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", bytes)
}

// createRequestInfo 创建请求信息
func (c *SDKClient) createRequestInfo(provider consts.Provider, modelInfo models.ModelInfo, method string) *middleware.RequestInfo {
	return &middleware.RequestInfo{
		Provider:   string(provider),
		ModelType:  string(modelInfo.ModelType),
		Model:      modelInfo.Model,
		Method:     method,
		RequestID:  generateRequestID(),
		StartTime:  time.Now(),
		RetryCount: 0,
	}
}

// ListModels 列出模型
func (c *SDKClient) ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 创建请求信息
	requestInfo := c.createRequestInfo(provider, models.ModelInfo{Provider: provider}, "ListModels")
	ctx = middleware.SetRequestInfo(ctx, requestInfo)
	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// 获取提供商
		var ps core.ProviderService
		if ps, err = core.GetProvider(provider); err != nil {
			return nil, err
		}
		// 列出模型
		return ps.ListModels(ctx, opts...)
	}
	// 执行中间件链
	resp, err := c.middlewareChain.Execute(ctx, nil, finalHandler)
	if err != nil {
		return models.ListModelsResponse{}, err
	}
	return resp.(models.ListModelsResponse), nil
}

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 创建请求信息
	requestInfo := c.createRequestInfo(request.ModelInfo.Provider, request.ModelInfo, "CreateChatCompletion")
	ctx = middleware.SetRequestInfo(ctx, requestInfo)

	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		chatReq := req.(models.ChatRequest)
		// 获取提供商
		var ps core.ProviderService
		if ps, err = core.GetProvider(chatReq.ModelInfo.Provider); err != nil {
			return nil, err
		}
		// 判断模型是否支持
		if err = core.IsModelSupported(ps, chatReq.ModelInfo); err != nil {
			return nil, err
		}
		// 判断是否流式传输
		if chatReq.Stream {
			return nil, consts.ErrCompletionStreamNotSupported
		}
		// 创建聊天
		return ps.CreateChatCompletion(ctx, chatReq, opts...)
	}

	// 执行中间件链
	resp, err := c.middlewareChain.Execute(ctx, request, finalHandler)
	if err != nil {
		return models.ChatResponse{}, err
	}

	return resp.(models.ChatResponse), nil
}

// CreateChatCompletionStream  创建聊天
func (c *SDKClient) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (interface{}, error) {
	// 创建请求信息
	requestInfo := c.createRequestInfo(request.ModelInfo.Provider, request.ModelInfo, "CreateChatCompletion")
	ctx = middleware.SetRequestInfo(ctx, requestInfo)
	var err error

	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		chatReq := req.(models.ChatRequest)
		// 获取提供商
		var ps core.ProviderService
		if ps, err = core.GetProvider(chatReq.ModelInfo.Provider); err != nil {
			return nil, err
		}
		// 判断模型是否支持
		if err = core.IsModelSupported(ps, chatReq.ModelInfo); err != nil {
			return nil, err
		}
		// 判断是否流式传输
		if !chatReq.Stream {
			return nil, consts.ErrCompletionNotStream
		}
		// 创建聊天
		return ps.CreateChatCompletionStream(ctx, chatReq, cb, opts...)
	}

	// 执行中间件链
	_, err = c.middlewareChain.Execute(ctx, request, finalHandler)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetMetrics 获取指标数据（如果启用了监控中间件）
func (c *SDKClient) GetMetrics() map[string]interface{} {
	for _, mw := range c.middlewareChain.GetMiddlewares() {
		if metricsMiddleware, ok := mw.(*middleware.MetricsMiddleware); ok {
			return metricsMiddleware.GetMetrics()
		}
	}
	return nil
}
