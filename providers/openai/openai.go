/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-13 19:21:46
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

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

// openAIProvider OpenAI提供商
type openAIProvider struct {
	supportedModels map[consts.ModelType]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig                 // 提供商配置
	hClient         *httpclient.HTTPClient               // HTTP 客户端
	lb              *loadbalancer.LoadBalancer           // 负载均衡器
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
		supportedModels: map[consts.ModelType]map[string]bool{
			consts.ChatModel: {
				// chat
				consts.OpenAIO1Mini:                         true,
				consts.OpenAIO1Mini20240912:                 true,
				consts.OpenAIO1Preview:                      true,
				consts.OpenAIO1Preview20240912:              true,
				consts.OpenAIO1:                             true,
				consts.OpenAIO120241217:                     true,
				consts.OpenAIO1Pro:                          true,
				consts.OpenAIO1Pro20250319:                  true,
				consts.OpenAIO3:                             true,
				consts.OpenAIO320250416:                     true,
				consts.OpenAIO3Mini:                         true,
				consts.OpenAIO3Mini20250131:                 true,
				consts.OpenAIO4Mini:                         true,
				consts.OpenAIO4Mini20250416:                 true,
				consts.OpenAIGPT432K0613:                    true,
				consts.OpenAIGPT432K0314:                    true,
				consts.OpenAIGPT432K:                        true,
				consts.OpenAIGPT40613:                       true,
				consts.OpenAIGPT40314:                       true,
				consts.OpenAIGPT4o:                          true,
				consts.OpenAIGPT4o20240513:                  true,
				consts.OpenAIGPT4o20240806:                  true,
				consts.OpenAIGPT4o20241120:                  true,
				consts.OpenAIChatGPT4oLatest:                true,
				consts.OpenAIGPT4oMini:                      true,
				consts.OpenAIGPT4oMini20240718:              true,
				consts.OpenAIGPT4oSearchPreview:             true,
				consts.OpenAIGPT4oSearchPreview20250311:     true,
				consts.OpenAIGPT4oMiniSearchPreview:         true,
				consts.OpenAIGPT4oMiniSearchPreview20250311: true,
				consts.OpenAIGPT4Turbo:                      true,
				consts.OpenAIGPT4TurboPreview:               true,
				consts.OpenAIGPT4Turbo20240409:              true,
				consts.OpenAIGPT40125Preview:                true,
				consts.OpenAIGPT41106Preview:                true,
				consts.OpenAIGPT4VisionPreview:              true,
				consts.OpenAIGPT4:                           true,
				consts.OpenAIGPT41:                          true,
				consts.OpenAIGPT4120250414:                  true,
				consts.OpenAIGPT41Mini:                      true,
				consts.OpenAIGPT41Mini20250414:              true,
				consts.OpenAIGPT41Nano:                      true,
				consts.OpenAIGPT41Nano20250414:              true,
				consts.OpenAIGPT45Preview:                   true,
				consts.OpenAIGPT45Preview20250227:           true,
				consts.OpenAIGPT35Turbo0125:                 true,
				consts.OpenAIGPT35Turbo1106:                 true,
				consts.OpenAIGPT3Dot5Turbo0613:              true,
				consts.OpenAIGPT3Dot5Turbo0301:              true,
				consts.OpenAIGPT35Turbo16k:                  true,
				consts.OpenAIGPT3Dot5Turbo16K0613:           true,
				consts.OpenAIGPT35Turbo:                     true,
				consts.OpenAIGPT35TurboInstruct:             true,
				consts.OpenAIGPT35TurboInstruct0914:         true,
				consts.OpenAIDavinci002:                     true,
				consts.OpenAIBabbage002:                     true,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                true,
				consts.OpenAIGPT4oAudioPreview20241001:        true,
				consts.OpenAIGPT4oAudioPreview20241217:        true,
				consts.OpenAIGPT4oAudioPreview20250603:        true,
				consts.OpenAIGPT4oRealtimePreview:             true,
				consts.OpenAIGPT4oRealtimePreview20241001:     true,
				consts.OpenAIGPT4oRealtimePreview20241217:     true,
				consts.OpenAIGPT4oRealtimePreview20250603:     true,
				consts.OpenAIGPT4oMiniAudioPreview:            true,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    true,
				consts.OpenAIGPT4oMiniRealtimePreview:         true,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: true,
			},
			consts.ImageModel: {
				// image
				consts.OpenAIDallE2:    true,
				consts.OpenAIDallE3:    true,
				consts.OpenAIGPTImage1: true,
			},
			consts.AudioModel: {
				// audio
				consts.OpenAITTS1:                true,
				consts.OpenAITTS11106:            true,
				consts.OpenAITTS1HD:              true,
				consts.OpenAITTS1HD1106:          true,
				consts.OpenAIWhisper1:            true,
				consts.OpenAIGPT4oTranscribe:     true,
				consts.OpenAIGPT4oMiniTranscribe: true,
				consts.OpenAIGPT4oMiniTTS:        true,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                true,
				consts.OpenAIGPT4oAudioPreview20241001:        true,
				consts.OpenAIGPT4oAudioPreview20241217:        true,
				consts.OpenAIGPT4oAudioPreview20250603:        true,
				consts.OpenAIGPT4oRealtimePreview:             true,
				consts.OpenAIGPT4oRealtimePreview20241001:     true,
				consts.OpenAIGPT4oRealtimePreview20241217:     true,
				consts.OpenAIGPT4oRealtimePreview20250603:     true,
				consts.OpenAIGPT4oMiniAudioPreview:            true,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    true,
				consts.OpenAIGPT4oMiniRealtimePreview:         true,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: true,
			},
			// moderation
			consts.ModerationModel: {
				consts.OpenAIOmniModerationLatest:   true,
				consts.OpenAIOmniModeration20240926: true,
			},
			// embed
			consts.EmbedModel: {
				consts.OpenAITextEmbedding3Small: true,
				consts.OpenAITextEmbedding3Large: true,
				consts.OpenAITextEmbeddingAda002: true,
			},
		},
	}
	core.RegisterProvider(consts.OpenAI, openaiService)
}

// GetSupportedModels 获取支持的模型
func (s *openAIProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *openAIProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.hClient = httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *openAIProvider) ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = s.executeRequest(ctx, http.MethodGet, apiModels, opts, &response)
	return
}

// CreateChatCompletion 创建聊天
func (s *openAIProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = s.executeRequest(ctx, http.MethodPost, apiChatCompletions, opts, &response, httpclient.WithBody(request))
	return
}

// executeRequest 执行请求
func (s *openAIProvider) executeRequest(ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, response httpclient.Response, reqSetters ...httpclient.RequestOption) (err error) {
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
