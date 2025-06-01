/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 14:59:37
 * @Description: DeepSeek服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpClient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"github.com/liusuxian/go-aisdk/models"
	"net/http"
)

// deepseekProvider DeepSeek提供商
type deepseekProvider struct {
	supportedModels map[consts.ModelType][]string // 支持的模型
	providerConfig  *conf.ProviderConfig          // 提供商配置
	httpClient      *httpClient.HTTPClient        // HTTP 客户端
	lb              *loadbalancer.LoadBalancer    // 负载均衡器
}

var (
	deepseekService *deepseekProvider // deepseek提供商实例
)

const (
	apiModels          = "/models"
	apiChatCompletions = "/chat/completions"
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	deepseekService = &deepseekProvider{
		supportedModels: map[consts.ModelType][]string{
			consts.ChatModel: {
				consts.DeepSeekChat,
				consts.DeepSeekReasoner,
			},
		},
	}
	core.RegisterProvider(consts.DeepSeek, deepseekService)
}

// GetSupportedModels 获取支持的模型
func (s *deepseekProvider) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *deepseekProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.httpClient = httpClient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// TODO ListModels 列出模型
func (s *deepseekProvider) ListModels(ctx context.Context, opts ...httpClient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.httpClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	var (
		setters = []httpClient.RequestOption{
			httpClient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
		}
		req *http.Request
	)
	if req, err = s.httpClient.NewRequest(ctx, http.MethodGet, s.httpClient.FullURL(apiModels), setters...); err != nil {
		return
	}
	err = s.httpClient.SendRequest(req, &response)
	return
}

// TODO CreateChatCompletion 创建聊天
func (s *deepseekProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpClient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.httpClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	var (
		setters = []httpClient.RequestOption{
			httpClient.WithBody(request),
			httpClient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
		}
		req *http.Request
	)
	if req, err = s.httpClient.NewRequest(ctx, http.MethodPost, s.httpClient.FullURL(apiChatCompletions), setters...); err != nil {
		return
	}
	err = s.httpClient.SendRequest(req, &response)
	return
}
