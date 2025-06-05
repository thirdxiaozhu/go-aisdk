/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 15:14:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-06 04:07:41
 * @Description: 中间件接口定义
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package middleware

import (
	"context"
	"time"
)

// Handler 处理器函数类型
type Handler func(ctx context.Context, request any) (response any, err error)

// Middleware 中间件接口
type Middleware interface {
	Process(ctx context.Context, request any, next Handler) (response any, err error) // 处理请求，next是下一个处理器
	Name() (name string)                                                              // 返回中间件名称，用于标识和排序
	Priority() (priority int)                                                         // 返回中间件优先级，数值越小优先级越高
}

// Chain 中间件链
type Chain struct {
	middlewares []Middleware
}

// NewChain 创建新的中间件链
func NewChain(middlewares ...Middleware) (chain *Chain) {
	return &Chain{
		middlewares: middlewares,
	}
}

// Execute 执行中间件链
func (c *Chain) Execute(ctx context.Context, request any, finalHandler Handler) (response any, err error) {
	if len(c.middlewares) == 0 {
		return finalHandler(ctx, request)
	}
	// 构建中间件调用链
	handler := finalHandler
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		var (
			mw          = c.middlewares[i]
			nextHandler = handler
		)
		handler = func(ctx context.Context, request any) (response any, err error) {
			return mw.Process(ctx, request, nextHandler)
		}
	}
	return handler(ctx, request)
}

// GetMiddlewares 获取中间件列表
func (c *Chain) GetMiddlewares() (middlewares []Middleware) {
	return c.middlewares
}

// RequestInfo 请求信息
type RequestInfo struct {
	Provider    string        `json:"provider"`     // 提供商
	ModelType   string        `json:"model_type"`   // 模型类型
	Model       string        `json:"model"`        // 模型名称
	Method      string        `json:"method"`       // 方法名称
	StartTime   time.Time     `json:"start_time"`   // 请求开始时间
	EndTime     time.Time     `json:"end_time"`     // 最新一次请求结束时间（重试过程中会更新）
	Duration    time.Duration `json:"duration"`     // 累计请求耗时（包含所有重试）
	IsSuccess   bool          `json:"is_success"`   // 当前请求状态（重试过程中会更新，最终表示是否成功）
	LastError   error         `json:"last_error"`   // 最后一次的错误信息（重试过程中会更新）
	RequestID   string        `json:"request_id"`   // 请求ID
	UserID      string        `json:"user_id"`      // 用户ID
	Attempt     int           `json:"attempt"`      // 第几次重试
	MaxAttempts int           `json:"max_attempts"` // 最大重试次数
}

// ContextKey 上下文键类型
type ContextKey string

const (
	RequestInfoKey ContextKey = "middleware_request_info" // 请求信息在上下文中的键
)

// GetRequestInfo 从上下文中获取请求信息
func GetRequestInfo(ctx context.Context) (reqInfo *RequestInfo) {
	if info, ok := ctx.Value(RequestInfoKey).(*RequestInfo); ok && info != nil {
		return info
	}
	// 返回安全的默认值
	return &RequestInfo{
		Provider:  "unknown",
		ModelType: "unknown",
		Model:     "unknown",
		Method:    "unknown",
		StartTime: time.Time{},
		RequestID: "unknown",
	}
}

// SetRequestInfo 设置请求信息到上下文
func SetRequestInfo(ctx context.Context, reqInfo *RequestInfo) (newCtx context.Context) {
	return context.WithValue(ctx, RequestInfoKey, reqInfo)
}
