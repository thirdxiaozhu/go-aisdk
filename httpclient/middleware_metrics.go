/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-17 18:24:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-17 20:47:17
 * @Description: 监控中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"
)

const (
	maxResponseTimeRecords = 1000 // 每个key最多保存1000条响应时间记录
)

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	// 记录请求开始
	RecordRequestStart(provider, modelType, model, method string)
	// 记录请求完成
	RecordRequestComplete(provider, modelType, model, method string, duration time.Duration, success bool)
	// 记录错误
	RecordError(provider, modelType, model, method string, errorType string)
	// 记录重试
	RecordRetry(provider, modelType, model, method string, retryCount int)
	// 获取指标数据
	GetMetrics() map[string]any
	// 重置指标
	Reset()
}

// DefaultMetricsCollector 默认指标收集器
type DefaultMetricsCollector struct {
	mu sync.RWMutex
	// 请求计数器
	totalRequests   map[string]int64 // 总请求数
	successRequests map[string]int64 // 成功请求数
	failedRequests  map[string]int64 // 失败请求数
	// 响应时间统计
	responseTimes map[string][]int64 // 响应时间列表（毫秒）
	// 错误统计
	errorCounts map[string]int64 // 错误计数
	// 重试统计
	retryCounts map[string]int64 // 重试计数
	// 活跃请求数
	activeRequests map[string]int64 // 当前活跃请求数
	// 时间范围内的统计
	startTime time.Time // 统计开始时间
}

// NewDefaultMetricsCollector 创建默认指标收集器
func NewDefaultMetricsCollector() (collector *DefaultMetricsCollector) {
	return &DefaultMetricsCollector{
		totalRequests:   make(map[string]int64),
		successRequests: make(map[string]int64),
		failedRequests:  make(map[string]int64),
		responseTimes:   make(map[string][]int64),
		errorCounts:     make(map[string]int64),
		retryCounts:     make(map[string]int64),
		activeRequests:  make(map[string]int64),
		startTime:       time.Now(),
	}
}

// getKey 获取指标键值
func (c *DefaultMetricsCollector) getKey(provider, modelType, model, method string) (key string) {
	for i, v := range []string{provider, modelType, model, method} {
		if i == 0 {
			key = v // provider必需存在
		} else if v != "" {
			key = fmt.Sprintf("%s:%s", key, v)
		}
	}
	return
}

// RecordRequestStart 记录请求开始
func (c *DefaultMetricsCollector) RecordRequestStart(provider, modelType, model, method string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(provider, modelType, model, method)
	c.totalRequests[key]++
	c.activeRequests[key]++
}

// RecordRequestComplete 记录请求完成
func (c *DefaultMetricsCollector) RecordRequestComplete(provider, modelType, model, method string, duration time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(provider, modelType, model, method)
	// 记录响应时间，限制记录数量防止内存泄漏
	durationMs := duration.Milliseconds()
	if _, exists := c.responseTimes[key]; !exists {
		c.responseTimes[key] = make([]int64, 0, maxResponseTimeRecords)
	}
	// 如果已达到最大记录数，移除最旧的记录（使用更高效的方式）
	if len(c.responseTimes[key]) >= maxResponseTimeRecords {
		// 保留后面80%的记录，丢弃前面20%
		keepFromIndex := maxResponseTimeRecords / 5 // 20%的位置
		c.responseTimes[key] = append(c.responseTimes[key][:0], c.responseTimes[key][keepFromIndex:]...)
	}
	c.responseTimes[key] = append(c.responseTimes[key], durationMs)
	// 记录成功/失败
	if success {
		c.successRequests[key]++
	} else {
		c.failedRequests[key]++
	}
	// 减少活跃请求数
	if c.activeRequests[key] > 0 {
		c.activeRequests[key]--
	}
}

// RecordError 记录错误
func (c *DefaultMetricsCollector) RecordError(provider, modelType, model, method string, errorType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(provider, modelType, model, method) + ":" + errorType
	c.errorCounts[key]++
}

// RecordRetry 记录重试
func (c *DefaultMetricsCollector) RecordRetry(provider, modelType, model, method string, retryCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(provider, modelType, model, method)
	c.retryCounts[key] += int64(retryCount)
}

// GetMetrics 获取指标数据
func (c *DefaultMetricsCollector) GetMetrics() (metrics map[string]any) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics = make(map[string]any)
	// 深拷贝基础统计数据
	totalRequests := make(map[string]int64)
	maps.Copy(totalRequests, c.totalRequests)
	// 深拷贝成功请求数据
	successRequests := make(map[string]int64)
	maps.Copy(successRequests, c.successRequests)
	// 深拷贝失败请求数据
	failedRequests := make(map[string]int64)
	maps.Copy(failedRequests, c.failedRequests)
	// 深拷贝活跃请求数据
	activeRequests := make(map[string]int64)
	maps.Copy(activeRequests, c.activeRequests)
	// 深拷贝重试计数数据
	retryCounts := make(map[string]int64)
	maps.Copy(retryCounts, c.retryCounts)
	// 深拷贝错误计数数据
	errorCounts := make(map[string]int64)
	maps.Copy(errorCounts, c.errorCounts)
	// 基础统计
	metrics["total_requests"] = totalRequests
	metrics["success_requests"] = successRequests
	metrics["failed_requests"] = failedRequests
	metrics["active_requests"] = activeRequests
	metrics["retry_counts"] = retryCounts
	metrics["error_counts"] = errorCounts
	// 计算成功率和平均响应时间
	successRates := make(map[string]float64)
	avgResponseTimes := make(map[string]float64)
	// 计算成功率和平均响应时间
	for key, total := range totalRequests {
		if total > 0 {
			success := successRequests[key]
			successRates[key] = float64(success) / float64(total) * 100
		}
		// 计算平均响应时间
		if times, exists := c.responseTimes[key]; exists && len(times) > 0 {
			var sum int64
			for _, t := range times {
				sum += t
			}
			avgResponseTimes[key] = float64(sum) / float64(len(times))
		}
	}
	// 添加成功率、平均响应时间和运行时间
	metrics["success_rates"] = successRates
	metrics["avg_response_times"] = avgResponseTimes
	metrics["uptime_seconds"] = time.Since(c.startTime).Seconds()
	return
}

// Reset 重置指标数据
func (c *DefaultMetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalRequests = make(map[string]int64)
	c.successRequests = make(map[string]int64)
	c.failedRequests = make(map[string]int64)
	c.responseTimes = make(map[string][]int64)
	c.errorCounts = make(map[string]int64)
	c.retryCounts = make(map[string]int64)
	c.activeRequests = make(map[string]int64)
	c.startTime = time.Now()
}

// MetricsMiddlewareConfig 监控中间件配置
type MetricsMiddlewareConfig struct {
	Collector      MetricsCollector // 指标收集器
	EnableDetailed bool             // 是否启用详细指标
}

// MetricsMiddleware 监控中间件
type MetricsMiddleware struct {
	config MetricsMiddlewareConfig
}

// NewMetricsMiddleware 创建监控中间件
func NewMetricsMiddleware(config MetricsMiddlewareConfig) (mm *MetricsMiddleware) {
	if config.Collector == nil {
		config.Collector = NewDefaultMetricsCollector()
	}

	return &MetricsMiddleware{
		config: config,
	}
}

// Name 返回中间件名称
func (m *MetricsMiddleware) Name() (name string) {
	return "metrics"
}

// Priority 返回中间件优先级
func (m *MetricsMiddleware) Priority() (priority int) {
	return 10 // 监控中间件优先级较高，尽早执行
}

// Process 处理请求
func (m *MetricsMiddleware) Process(ctx context.Context, req any, next Handler) (resp any, err error) {
	requestInfo := GetRequestInfo(ctx)
	// 记录请求开始
	m.config.Collector.RecordRequestStart(
		requestInfo.Provider,
		requestInfo.ModelType,
		requestInfo.Model,
		requestInfo.Method,
	)
	// 执行下一个处理器
	startTime := time.Now()
	resp, err = next(ctx, req)
	duration := time.Since(startTime)
	// 记录请求完成
	success := err == nil
	m.config.Collector.RecordRequestComplete(
		requestInfo.Provider,
		requestInfo.ModelType,
		requestInfo.Model,
		requestInfo.Method,
		duration,
		success,
	)
	// 记录错误
	if err != nil {
		errorType := m.classifyError(err)
		m.config.Collector.RecordError(
			requestInfo.Provider,
			requestInfo.ModelType,
			requestInfo.Model,
			requestInfo.Method,
			errorType,
		)
	}
	// 记录重试次数
	if requestInfo.Attempt > 0 {
		m.config.Collector.RecordRetry(
			requestInfo.Provider,
			requestInfo.ModelType,
			requestInfo.Model,
			requestInfo.Method,
			requestInfo.Attempt,
		)
	}
	return
}

// GetMetrics 获取指标数据
func (m *MetricsMiddleware) GetMetrics() (metrics map[string]any) {
	return m.config.Collector.GetMetrics()
}

// classifyError 分类错误类型
func (m *MetricsMiddleware) classifyError(err error) (errorType string) {
	if err == nil {
		return "none"
	}

	errStr := strings.ToLower(err.Error())
	// 网络相关错误
	if strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dns") {
		return "network"
	}
	// 超时错误
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline") {
		return "timeout"
	}
	// HTTP状态码错误
	if strings.Contains(errStr, "400") {
		return "bad_request"
	}
	if strings.Contains(errStr, "401") {
		return "unauthorized"
	}
	if strings.Contains(errStr, "403") {
		return "forbidden"
	}
	if strings.Contains(errStr, "404") {
		return "not_found"
	}
	if strings.Contains(errStr, "429") {
		return "rate_limit"
	}
	if strings.Contains(errStr, "5") && (strings.Contains(errStr, "500") ||
		strings.Contains(errStr, "502") || strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504")) {
		return "server_error"
	}
	// 认证相关错误
	if strings.Contains(errStr, "api") && strings.Contains(errStr, "key") {
		return "api_key_error"
	}
	// 模型相关错误
	if strings.Contains(errStr, "model") {
		return "model_error"
	}
	// 其他未知错误
	return "unknown"
}

// DefaultMetricsConfig 默认监控配置
func DefaultMetricsConfig() (config MetricsMiddlewareConfig) {
	return MetricsMiddlewareConfig{
		Collector:      NewDefaultMetricsCollector(),
		EnableDetailed: true,
	}
}
