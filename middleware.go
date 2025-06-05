/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-02 04:49:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-06 01:07:32
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient/middleware"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/sdkerrors"
	"time"
)

// SDKClientOption SDK客户端选项
type SDKClientOption func(c *clientOption)

// clientOption 客户端选项
type clientOption struct {
	middlewares []middleware.Middleware
}

// WithMiddleware 添加中间件
func WithMiddleware(m middleware.Middleware) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, m)
	}
}

// WithLogging 添加日志中间件
func WithLogging(config middleware.LoggingMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, middleware.NewLoggingMiddleware(config))
	}
}

// WithMetrics 添加监控中间件
func WithMetrics(config middleware.MetricsMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, middleware.NewMetricsMiddleware(config))
	}
}

// WithRetry 添加重试中间件
func WithRetry(config middleware.RetryMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, middleware.NewRetryMiddleware(config))
	}
}

// WithDefaultMiddlewares 添加默认中间件（日志、监控、重试）
func WithDefaultMiddlewares() (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares,
			middleware.NewMetricsMiddleware(middleware.DefaultMetricsConfig()),
			middleware.NewRetryMiddleware(middleware.DefaultRetryConfig()),
			middleware.NewLoggingMiddleware(middleware.DefaultLoggingConfig()),
		)
	}
}

// GetMetrics 获取指标数据（如果启用了监控中间件）
func (c *SDKClient) GetMetrics() (metrics map[string]any) {
	for _, mw := range c.middlewareChain.GetMiddlewares() {
		if metricsMiddleware, ok := mw.(*middleware.MetricsMiddleware); ok {
			return metricsMiddleware.GetMetrics()
		}
	}
	return
}

// handlerRequest 处理请求
func (c *SDKClient) handlerRequest(
	ctx context.Context,
	userId string,
	modelInfo models.ModelInfo,
	method string,
	request any,
	handler func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error),
) (resp any, err error) {
	// 生成唯一请求ID
	var requestId string
	if requestId, err = c.flakeInstance.RequestID(); err != nil {
		err = &sdkerrors.SDKError{RequestID: requestId, Err: err}
		return
	}
	// 设置请求信息到上下文
	ctx = middleware.SetRequestInfo(ctx, &middleware.RequestInfo{
		Provider:  string(modelInfo.Provider),
		ModelType: string(modelInfo.ModelType),
		Model:     modelInfo.Model,
		Method:    method,
		StartTime: time.Now(),
		RequestID: requestId,
		UserID:    userId,
	})
	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req any) (resp any, err error) {
		// 获取提供商
		var (
			ps core.ProviderService
			e  error
		)
		if ps, e = core.GetProvider(modelInfo.Provider); e != nil {
			return nil, e
		}
		// 根据方法名称决定是否需要判断模型支持
		if !c.noCheckMethods[method] {
			// 判断模型是否支持
			if e = core.IsModelSupported(ps, modelInfo); e != nil {
				return nil, e
			}
		}
		// 执行具体的处理逻辑
		return handler(ctx, ps, req)
	}
	// 执行中间件链
	if resp, err = c.middlewareChain.Execute(ctx, request, finalHandler); err != nil {
		err = &sdkerrors.SDKError{RequestID: requestId, Err: err}
		return
	}
	return
}
