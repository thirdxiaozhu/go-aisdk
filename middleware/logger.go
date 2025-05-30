/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 15:14:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 19:02:54
 * @Description: 日志中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

const (
	maxSanitizeDepth = 10 // 最大脱敏递归深度
)

// Logger 日志接口
type Logger interface {
	Debug(ctx context.Context, format string, args ...any) // 调试日志
	Info(ctx context.Context, format string, args ...any)  // 信息日志
	Warn(ctx context.Context, format string, args ...any)  // 警告日志
	Error(ctx context.Context, format string, args ...any) // 错误日志
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	logger *log.Logger
	level  LogLevel
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger(level LogLevel) (l *DefaultLogger) {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

// Debug 调试日志
func (l *DefaultLogger) Debug(ctx context.Context, format string, args ...any) {
	if LogLevelDebug < l.level {
		return
	}
	l.logger.Printf("[DEBUG] %s", fmt.Sprintf(format, args...))
}

// Info 信息日志
func (l *DefaultLogger) Info(ctx context.Context, format string, args ...any) {
	if LogLevelInfo < l.level {
		return
	}
	l.logger.Printf("[INFO] %s", fmt.Sprintf(format, args...))
}

// Warn 警告日志
func (l *DefaultLogger) Warn(ctx context.Context, format string, args ...any) {
	if LogLevelWarn < l.level {
		return
	}
	l.logger.Printf("[WARN] %s", fmt.Sprintf(format, args...))
}

// Error 错误日志
func (l *DefaultLogger) Error(ctx context.Context, format string, args ...any) {
	if LogLevelError < l.level {
		return
	}
	l.logger.Printf("[ERROR] %s", fmt.Sprintf(format, args...))
}

// LoggingMiddlewareConfig 日志中间件配置
type LoggingMiddlewareConfig struct {
	Logger          Logger   // 日志器
	LogRequest      bool     // 是否记录请求
	LogResponse     bool     // 是否记录响应
	LogError        bool     // 是否记录错误
	SkipSuccessLog  bool     // 是否跳过成功请求的日志
	SensitiveFields []string // 敏感字段，会被脱敏
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	config LoggingMiddlewareConfig
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(config LoggingMiddlewareConfig) (lm *LoggingMiddleware) {
	if config.Logger == nil {
		config.Logger = NewDefaultLogger(LogLevelInfo)
	}

	return &LoggingMiddleware{
		config: config,
	}
}

// Process 处理请求
func (m *LoggingMiddleware) Process(ctx context.Context, request any, next Handler) (response any, err error) {
	requestInfo := GetRequestInfo(ctx)
	// 是否记录请求
	if m.config.LogRequest {
		startFields := map[string]any{
			"request_id": requestInfo.RequestID,
			"provider":   requestInfo.Provider,
			"model_type": requestInfo.ModelType,
			"model":      requestInfo.Model,
			"method":     requestInfo.Method,
			"user_id":    requestInfo.UserID,
		}
		// 脱敏处理请求数据
		if reqData := m.sanitizeData(request); reqData != nil {
			startFields["request"] = reqData
		}
		m.config.Logger.Info(ctx, "request started: %v", mustJsonString(startFields))
	}
	// 执行下一个处理器
	startTime := time.Now()
	response, err = next(ctx, request)
	duration := time.Since(startTime)
	// 记录请求结果
	endFields := map[string]any{
		"request_id":  requestInfo.RequestID,
		"provider":    requestInfo.Provider,
		"model_type":  requestInfo.ModelType,
		"model":       requestInfo.Model,
		"method":      requestInfo.Method,
		"duration":    formatDuration(duration),
		"success":     err == nil,
		"user_id":     requestInfo.UserID,
		"retry_count": requestInfo.RetryCount,
	}
	if err != nil {
		// 是否记录错误
		if m.config.LogError {
			endFields["error"] = err.Error()
			m.config.Logger.Error(ctx, "request failed: %v", mustJsonString(endFields))
		}
	} else {
		// 是否跳过成功请求的日志
		if !m.config.SkipSuccessLog {
			// 是否记录响应
			if m.config.LogResponse {
				// 脱敏处理响应数据
				if respData := m.sanitizeData(response); respData != nil {
					endFields["response"] = respData
				}
			}
			m.config.Logger.Info(ctx, "request completed: %v", mustJsonString(endFields))
		}
	}
	return
}

// Name 返回中间件名称
func (m *LoggingMiddleware) Name() (name string) {
	return "logging"
}

// Priority 返回中间件优先级
func (m *LoggingMiddleware) Priority() (priority int) {
	return 100 // 日志中间件优先级较低，在其他中间件执行后记录
}

// sanitizeData 脱敏数据
func (m *LoggingMiddleware) sanitizeData(data any) (newData any) {
	if data == nil {
		return nil
	}
	// 将数据转换为map进行处理
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("failed to marshal data: %v", err)
	}
	var result any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return string(jsonData) // 如果无法解析，直接返回字符串
	}
	// 递归脱敏，从深度0开始
	return m.sanitizeValue(result, 0)
}

// sanitizeValue 递归脱敏值，添加深度限制
func (m *LoggingMiddleware) sanitizeValue(value any, depth int) (newValue any) {
	// 防止无限递归
	if depth > maxSanitizeDepth {
		return "<max_depth_reached>"
	}

	switch v := value.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, val := range v {
			// 检查是否为敏感字段
			isSensitive := false
			for _, field := range m.config.SensitiveFields {
				if strings.Contains(strings.ToLower(key), strings.ToLower(field)) {
					isSensitive = true
					break
				}
			}
			if isSensitive {
				result[key] = "***"
			} else {
				result[key] = m.sanitizeValue(val, depth+1) // 递归处理，深度+1
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, val := range v {
			result[i] = m.sanitizeValue(val, depth+1) // 递归处理数组元素，深度+1
		}
		return result
	default:
		return v // 基本类型直接返回
	}
}

// DefaultLoggingConfig 默认日志配置
func DefaultLoggingConfig() (config LoggingMiddlewareConfig) {
	return LoggingMiddlewareConfig{
		Logger:          NewDefaultLogger(LogLevelInfo),
		LogRequest:      true,
		LogResponse:     false, // 默认不记录响应以减少日志量
		LogError:        true,
		SkipSuccessLog:  false,
		SensitiveFields: []string{},
	}
}
