/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-04 11:56:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-04 22:29:30
 * @Description: 重试中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package middleware

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	"strings"
	"time"
)

// rng 随机数生成器
var rng *rand.Rand

// 包初始化时设置随机数种子
func init() {
	now := time.Now().UnixNano()
	rng = rand.New(rand.NewPCG(uint64(now), uint64(now>>32)))
}

// RetryStrategy 重试策略
type RetryStrategy string

const (
	RetryStrategyFixed       RetryStrategy = "fixed"       // 固定间隔
	RetryStrategyLinear      RetryStrategy = "linear"      // 线性递增
	RetryStrategyExponential RetryStrategy = "exponential" // 指数退避
	RetryStrategyJitter      RetryStrategy = "jitter"      // 带抖动的指数退避
)

// RetryCondition 重试条件函数
type RetryCondition func(attempt int, err error) (ok bool)

// RetryCallback 重试失败回调函数
type RetryCallback func(ctx context.Context, requestInfo RequestInfo)

// RetryMiddlewareConfig 重试中间件配置
type RetryMiddlewareConfig struct {
	MaxAttempts   int            // 最大重试次数
	Strategy      RetryStrategy  // 重试策略
	BaseDelay     time.Duration  // 基础延迟时间
	MaxDelay      time.Duration  // 最大延迟时间
	Multiplier    float64        // 重试间隔倍数（用于指数退避）
	JitterPercent float64        // 抖动百分比（用于抖动策略，范围0-1，如0.1表示±10%）
	Condition     RetryCondition // 重试条件
	OnRetry       RetryCallback  // 重试失败回调函数（同步执行会阻塞重试流程，建议仅用于轻量级操作如日志记录、监控上报等，耗时操作请在回调内使用goroutine异步处理）
}

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	config RetryMiddlewareConfig
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(config RetryMiddlewareConfig) (retry *RetryMiddleware) {
	// 设置最大重试次数
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	// 设置重试策略
	if config.Strategy == "" {
		config.Strategy = RetryStrategyExponential
	}
	// 设置基础延迟时间
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	// 设置最大延迟时间
	if config.MaxDelay <= 0 {
		config.MaxDelay = 10 * time.Second
	}
	// 确保最大延迟不小于基础延迟
	if config.MaxDelay < config.BaseDelay {
		config.MaxDelay = config.BaseDelay * 10 // 设置为基础延迟的10倍
	}
	// 设置重试间隔倍数
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}
	// 设置抖动百分比
	if config.JitterPercent <= 0 || config.JitterPercent > 1 {
		config.JitterPercent = 0.1 // 默认±10%抖动
	}
	// 设置重试条件
	if config.Condition == nil {
		config.Condition = DefaultRetryCondition
	}
	return &RetryMiddleware{
		config: config,
	}
}

// Process 处理请求
func (m *RetryMiddleware) Process(ctx context.Context, request any, next Handler) (response any, err error) {
	// 从上下文中获取请求信息
	requestInfo := GetRequestInfo(ctx)
	requestInfo.MaxAttempts = m.config.MaxAttempts

	for attempt := 0; attempt <= m.config.MaxAttempts; attempt++ {
		requestInfo.Attempt = attempt
		// 执行请求
		response, err = next(ctx, request)
		// 如果成功或者不需要重试，直接返回
		if err == nil || !m.config.Condition(attempt, err) {
			return
		}
		// 如果重试回调不为空，则执行重试回调
		if m.config.OnRetry != nil && attempt > 0 {
			// TODO
			m.config.OnRetry(ctx, *requestInfo)
		}
		// 如果是最后一次尝试，不需要等待
		if attempt == m.config.MaxAttempts {
			break
		}
		// 计算延迟时间
		delay := m.calculateDelay(attempt + 1)
		// 等待延迟时间
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// 继续下次重试
		}
	}
	// 返回最后一次的错误
	return nil, fmt.Errorf("after %d attempts, last error: %w", requestInfo.Attempt, err)
}

// Name 返回中间件名称
func (m *RetryMiddleware) Name() (name string) {
	return "retry"
}

// Priority 返回中间件优先级
func (m *RetryMiddleware) Priority() (priority int) {
	return 20 // 重试中间件优先级较高，在监控之后执行
}

// calculateDelay 计算延迟时间
func (m *RetryMiddleware) calculateDelay(attempt int) (delay time.Duration) {
	switch m.config.Strategy {
	case RetryStrategyFixed:
		delay = m.config.BaseDelay
	case RetryStrategyLinear:
		if attempt > 0 && m.config.BaseDelay > 0 {
			// 防止整数溢出，先检查是否会超过最大延迟
			maxAttempts := int64(m.config.MaxDelay / m.config.BaseDelay)
			if int64(attempt) > maxAttempts {
				delay = m.config.MaxDelay
			} else {
				delay = min(time.Duration(attempt)*m.config.BaseDelay, m.config.MaxDelay)
			}
		} else {
			delay = m.config.BaseDelay
		}
	case RetryStrategyExponential:
		delay = m.calculateExponentialDelay(attempt)
	case RetryStrategyJitter:
		exponentialDelay := m.calculateExponentialDelay(attempt)
		// 添加双向随机抖动（基于配置的百分比）
		jitterRange := float64(exponentialDelay) * m.config.JitterPercent
		jitter := time.Duration((rng.Float64() - 0.5) * 2 * jitterRange)
		// 确保抖动后的延迟在合理范围内
		delay = max(exponentialDelay+jitter, exponentialDelay/2)
	default:
		delay = m.config.BaseDelay
	}
	// 确保不超过最大延迟
	delay = min(delay, m.config.MaxDelay)
	// 确保延迟不为负数或过小
	delay = max(delay, m.config.BaseDelay/2)
	return
}

// calculateExponentialDelay 计算指数退避延迟
func (m *RetryMiddleware) calculateExponentialDelay(attempt int) (delay time.Duration) {
	// 如果基础延迟、最大延迟小于等于0或重试间隔倍数小于等于1，返回基础延迟
	if m.config.BaseDelay <= 0 || m.config.MaxDelay <= 0 || m.config.Multiplier <= 1 {
		return m.config.BaseDelay
	}
	// 计算最大尝试次数
	maxExponent := math.Log(float64(m.config.MaxDelay/m.config.BaseDelay)) / math.Log(m.config.Multiplier)
	// 如果最大尝试次数为无穷大或NaN，或者尝试次数大于最大尝试次数，返回最大延迟
	if math.IsInf(maxExponent, 0) || math.IsNaN(maxExponent) || float64(attempt-1) > maxExponent {
		return m.config.MaxDelay
	}
	// 计算延迟时间
	delayFloat := float64(m.config.BaseDelay) * math.Pow(m.config.Multiplier, float64(attempt-1))
	// 如果延迟时间为无穷大或NaN，返回最大延迟
	if math.IsInf(delayFloat, 0) || math.IsNaN(delayFloat) {
		return m.config.MaxDelay
	}
	// 返回延迟时间
	return time.Duration(delayFloat)
}

// DefaultRetryCondition 默认重试条件
func DefaultRetryCondition(attempt int, err error) (ok bool) {
	if err == nil {
		return false
	}
	// 网络错误通常可以重试
	if isNetworkError(err) {
		return true
	}
	// HTTP状态码错误
	if isRetryableHTTPError(err) {
		return true
	}
	// 超时错误
	if isTimeoutError(err) {
		return true
	}
	// 临时错误
	if isTemporaryError(err) {
		return true
	}
	return false
}

// isNetworkError 判断是否为网络错误
func isNetworkError(err error) (ok bool) {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	// DNS错误
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	// 连接错误
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "broken pipe") {
		return true
	}
	return false
}

// isRetryableHTTPError 判断是否为可重试的HTTP错误
func isRetryableHTTPError(err error) (ok bool) {
	// 检查是否包含HTTP状态码
	errStr := err.Error()
	// 5xx 服务器错误通常可以重试
	if strings.Contains(errStr, "500") || // Internal Server Error
		strings.Contains(errStr, "502") || // Bad Gateway
		strings.Contains(errStr, "503") || // Service Unavailable
		strings.Contains(errStr, "504") || // Gateway Timeout
		strings.Contains(errStr, "507") || // Insufficient Storage
		strings.Contains(errStr, "508") || // Loop Detected
		strings.Contains(errStr, "509") || // Bandwidth Limit Exceeded
		strings.Contains(errStr, "510") { // Not Extended
		return true
	}
	// 429 Too Many Requests 也可以重试
	if strings.Contains(errStr, "429") {
		return true
	}
	// 408 Request Timeout
	if strings.Contains(errStr, "408") {
		return true
	}
	return false
}

// isTimeoutError 判断是否为超时错误
func isTimeoutError(err error) (ok bool) {
	// 检查是否实现了Timeout接口
	type timeoutError interface {
		Timeout() (ok bool)
	}
	if te, ok := err.(timeoutError); ok {
		return te.Timeout()
	}
	// 检查错误信息中是否包含超时关键字
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context deadline exceeded")
}

// isTemporaryError 判断是否为临时错误
func isTemporaryError(err error) (ok bool) {
	// 检查是否实现了Temporary接口
	type temporaryError interface {
		Temporary() (ok bool)
	}
	if te, ok := err.(temporaryError); ok {
		return te.Temporary()
	}
	return false
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() (config RetryMiddlewareConfig) {
	return RetryMiddlewareConfig{
		MaxAttempts:   3,
		Strategy:      RetryStrategyExponential,
		BaseDelay:     1 * time.Second,
		MaxDelay:      10 * time.Second,
		Multiplier:    2.0,
		JitterPercent: 0.1, // 默认±10%抖动
		Condition:     DefaultRetryCondition,
		OnRetry: func(ctx context.Context, requestInfo RequestInfo) {
			fmt.Printf("111111111111: %+v\n", requestInfo)
		},
	}
}

// RetryConditions 预定义的重试条件
var RetryConditions = struct {
	// Never 永不重试
	Never RetryCondition
	// Always 总是重试（除了成功的情况）
	Always RetryCondition
	// NetworkOnly 仅网络错误重试
	NetworkOnly RetryCondition
	// HTTPOnly 仅HTTP错误重试
	HTTPOnly RetryCondition
}{
	Never: func(attempt int, err error) (ok bool) {
		return false
	},
	Always: func(attempt int, err error) (ok bool) {
		return err != nil
	},
	NetworkOnly: func(attempt int, err error) (ok bool) {
		return err != nil && isNetworkError(err)
	},
	HTTPOnly: func(attempt int, err error) (ok bool) {
		return err != nil && isRetryableHTTPError(err)
	},
}
