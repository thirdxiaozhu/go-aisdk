/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-02 04:49:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 15:54:41
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import "github.com/liusuxian/go-aisdk/httpclient/middleware"

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
