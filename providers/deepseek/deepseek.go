/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-13 19:15:31
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
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"github.com/liusuxian/go-aisdk/models"
	"net/http"
)

// deepseekProvider DeepSeek提供商
type deepseekProvider struct {
	supportedModels map[consts.ModelType]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig                 // 提供商配置
	hClient         *httpclient.HTTPClient               // HTTP 客户端
	lb              *loadbalancer.LoadBalancer           // 负载均衡器
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
		supportedModels: map[consts.ModelType]map[string]bool{
			consts.ChatModel: {
				// chat
				consts.DeepSeekChat:     true,
				consts.DeepSeekReasoner: true,
			},
		},
	}
	core.RegisterProvider(consts.DeepSeek, deepseekService)
}

// GetSupportedModels 获取支持的模型
func (s *deepseekProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *deepseekProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.hClient = httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *deepseekProvider) ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = s.executeRequest(ctx, http.MethodGet, apiModels, opts, &response)
	return
}

// CreateChatCompletion 创建聊天
func (s *deepseekProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = s.executeRequest(ctx, http.MethodPost, apiChatCompletions, opts, &response, httpclient.WithBody(request))
	return
}

// executeRequest 执行请求
func (s *deepseekProvider) executeRequest(ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, response httpclient.Response, reqSetters ...httpclient.RequestOption) (err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(reqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, method, s.hClient.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = s.hClient.SendRequest(req, response)
	return
}
