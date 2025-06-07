/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-05 19:13:26
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/internal/loadbalancer"
	"github.com/liusuxian/go-aisdk/models"
)

// openAIProvider OpenAI提供商
type openAIProvider struct {
	supportedModels map[consts.ModelType][]string // 支持的模型
	providerConfig  *conf.ProviderConfig          // 提供商配置
	hClient         *httpclient.HTTPClient        // HTTP 客户端
	lb              *loadbalancer.LoadBalancer    // 负载均衡器
}

var (
	openaiService *openAIProvider // OpenAI提供商实例
)

const (
	apiModels          = "/models"
	apiChatCompletions = "/chat/completions"
)

// init 包初始化时创建 openAIProvider 实例并注册到工厂
func init() {
	openaiService = &openAIProvider{
		supportedModels: map[consts.ModelType][]string{
			consts.ChatModel:  {},
			consts.ImageModel: {},
			consts.VideoModel: {},
			consts.AudioModel: {},
			consts.EmbedModel: {},
		},
	}
	core.RegisterProvider(consts.OpenAI, openaiService)
}

// GetSupportedModels 获取支持的模型
func (s *openAIProvider) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *openAIProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.hClient = httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// TODO ListModels 列出模型
func (s *openAIProvider) ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	var (
		setters = []httpclient.RequestOption{
			httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
		}
		req *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, http.MethodGet, s.hClient.FullURL(apiModels), setters...); err != nil {
		return
	}
	err = s.hClient.SendRequest(req, &response)
	return
}

// TODO CreateChatCompletion 创建聊天
func (s *openAIProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	return
}

// TODO CreateChatCompletionStream 创建聊天
func (s *openAIProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (interface{}, error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	return nil, nil
}
