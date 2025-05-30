/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-01-XX XX:XX:XX
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 19:01:37
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
type RetryCondition func(err error, attempt int) (ok bool)

// RetryMiddlewareConfig 重试中间件配置
type RetryMiddlewareConfig struct {
	MaxAttempts int                          // 最大重试次数
	Strategy    RetryStrategy                // 重试策略
	BaseDelay   time.Duration                // 基础延迟时间
	MaxDelay    time.Duration                // 最大延迟时间
	Multiplier  float64                      // 重试间隔倍数（用于指数退避）
	Condition   RetryCondition               // 重试条件
	OnRetry     func(attempt int, err error) // 重试回调函数
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
		config.MaxDelay = 30 * time.Second
	}
	// 设置重试间隔倍数
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
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
	requestInfo := GetRequestInfo(ctx)

	var lastErr error
	for attempt := 1; attempt <= m.config.MaxAttempts; attempt++ {
		// 创建当前尝试的请求信息副本，避免并发问题
		currentInfo := *requestInfo // 浅拷贝
		currentInfo.RetryCount = attempt - 1
		currentCtx := SetRequestInfo(ctx, &currentInfo)
		// 执行请求
		response, err = next(currentCtx, request)
		// 如果成功或者不需要重试，直接返回
		if err == nil || !m.config.Condition(err, attempt) {
			// 更新原始请求信息的重试计数
			requestInfo.RetryCount = attempt - 1
			return
		}
		// 记录最后一次错误
		lastErr = err
		// 如果是最后一次尝试，不需要等待
		if attempt == m.config.MaxAttempts {
			break
		}
		// 调用重试回调
		if m.config.OnRetry != nil {
			m.config.OnRetry(attempt, err)
		}

		// 计算延迟时间并等待
		delay := m.calculateDelay(attempt)

		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// 继续下次重试
		}
	}
	// 更新最终的重试计数
	requestInfo.RetryCount = m.config.MaxAttempts - 1
	// 返回最后一次的错误
	return nil, fmt.Errorf("after %d attempts, last error: %w", m.config.MaxAttempts, lastErr)
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
			if int64(attempt) > int64(m.config.MaxDelay)/int64(m.config.BaseDelay) {
				delay = m.config.MaxDelay
			} else {
				delay = time.Duration(attempt) * m.config.BaseDelay
			}
		} else {
			delay = m.config.BaseDelay
		}
	case RetryStrategyExponential:
		if attempt <= 1 {
			delay = m.config.BaseDelay
		} else {
			if m.config.BaseDelay <= 0 || m.config.MaxDelay <= 0 || m.config.Multiplier <= 1 {
				delay = m.config.BaseDelay
			} else {
				maxExponent := math.Log(float64(m.config.MaxDelay)/float64(m.config.BaseDelay)) / math.Log(m.config.Multiplier)
				if math.IsInf(maxExponent, 0) || math.IsNaN(maxExponent) || float64(attempt-1) > maxExponent {
					delay = m.config.MaxDelay
				} else {
					delayFloat := float64(m.config.BaseDelay) * math.Pow(m.config.Multiplier, float64(attempt-1))
					if math.IsInf(delayFloat, 0) || math.IsNaN(delayFloat) {
						delay = m.config.MaxDelay
					} else {
						delay = time.Duration(delayFloat)
					}
				}
			}
		}
	case RetryStrategyJitter:
		var exponentialDelay time.Duration
		if attempt <= 1 {
			exponentialDelay = m.config.BaseDelay
		} else {
			if m.config.BaseDelay <= 0 || m.config.MaxDelay <= 0 || m.config.Multiplier <= 1 {
				exponentialDelay = m.config.BaseDelay
			} else {
				maxExponent := math.Log(float64(m.config.MaxDelay)/float64(m.config.BaseDelay)) / math.Log(m.config.Multiplier)
				if math.IsInf(maxExponent, 0) || math.IsNaN(maxExponent) || float64(attempt-1) > maxExponent {
					exponentialDelay = m.config.MaxDelay
				} else {
					delayFloat := float64(m.config.BaseDelay) * math.Pow(m.config.Multiplier, float64(attempt-1))
					if math.IsInf(delayFloat, 0) || math.IsNaN(delayFloat) {
						exponentialDelay = m.config.MaxDelay
					} else {
						exponentialDelay = time.Duration(delayFloat)
					}
				}
			}
		}
		// 添加随机抖动（10%）
		jitter := time.Duration(rng.Float64() * float64(exponentialDelay) * 0.1)
		delay = exponentialDelay + jitter
	default:
		delay = m.config.BaseDelay
	}
	// 确保不超过最大延迟
	if delay > m.config.MaxDelay {
		delay = m.config.MaxDelay
	}
	// 确保延迟不为负数
	if delay < 0 {
		delay = m.config.BaseDelay
	}
	return
}

// DefaultRetryCondition 默认重试条件
func DefaultRetryCondition(err error, attempt int) (ok bool) {
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
		MaxAttempts: 3,
		Strategy:    RetryStrategyExponential,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Condition:   DefaultRetryCondition,
		OnRetry: func(attempt int, err error) {
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
	Never: func(err error, attempt int) (ok bool) {
		return false
	},
	Always: func(err error, attempt int) (ok bool) {
		return err != nil
	},
	NetworkOnly: func(err error, attempt int) (ok bool) {
		return err != nil && isNetworkError(err)
	},
	HTTPOnly: func(err error, attempt int) (ok bool) {
		return err != nil && isRetryableHTTPError(err)
	},
}
