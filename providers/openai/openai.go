/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 20:59:08
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

// openAIProvider OpenAI提供商
type openAIProvider struct {
	core.DefaultProviderService
	supportedModels map[consts.ModelType]map[string]uint // 支持的模型
	providerConfig  *conf.ProviderConfig                 // 提供商配置
	lb              *loadbalancer.LoadBalancer           // 负载均衡器
}

var (
	openaiService *openAIProvider // OpenAI提供商实例
)

const (
	apiModels = "/models"
)

// init 包初始化时创建 openAIProvider 实例并注册到工厂
func init() {
	openaiService = &openAIProvider{
		supportedModels: map[consts.ModelType]map[string]uint{
			consts.ChatModel: {
				// chat
				consts.OpenAIO1Mini:                         0,
				consts.OpenAIO1Mini20240912:                 0,
				consts.OpenAIO1Preview:                      0,
				consts.OpenAIO1Preview20240912:              0,
				consts.OpenAIO1:                             0,
				consts.OpenAIO1_20241217:                    0,
				consts.OpenAIO1Pro:                          0,
				consts.OpenAIO1Pro20250319:                  0,
				consts.OpenAIO3:                             1,
				consts.OpenAIO3_20250416:                    1,
				consts.OpenAIO3Mini:                         1,
				consts.OpenAIO3Mini20250131:                 1,
				consts.OpenAIO4Mini:                         1,
				consts.OpenAIO4Mini20250416:                 1,
				consts.OpenAIGPT4_32K0613:                   0,
				consts.OpenAIGPT4_32K0314:                   0,
				consts.OpenAIGPT4_32K:                       0,
				consts.OpenAIGPT4_0613:                      0,
				consts.OpenAIGPT4_0314:                      0,
				consts.OpenAIGPT4o:                          1,
				consts.OpenAIGPT4o20240513:                  1,
				consts.OpenAIGPT4o20240806:                  1,
				consts.OpenAIGPT4o20241120:                  1,
				consts.OpenAIChatGPT4oLatest:                1,
				consts.OpenAIGPT4oMini:                      1,
				consts.OpenAIGPT4oMini20240718:              1,
				consts.OpenAIGPT4oSearchPreview:             1,
				consts.OpenAIGPT4oSearchPreview20250311:     1,
				consts.OpenAIGPT4oMiniSearchPreview:         1,
				consts.OpenAIGPT4oMiniSearchPreview20250311: 1,
				consts.OpenAIGPT4Turbo:                      1,
				consts.OpenAIGPT4TurboPreview:               1,
				consts.OpenAIGPT4Turbo20240409:              1,
				consts.OpenAIGPT4_0125Preview:               1,
				consts.OpenAIGPT4_1106Preview:               1,
				consts.OpenAIGPT4VisionPreview:              1,
				consts.OpenAIGPT4:                           0,
				consts.OpenAIGPT4Dot1:                       1,
				consts.OpenAIGPT4Dot1_20250414:              1,
				consts.OpenAIGPT4Dot1Mini:                   1,
				consts.OpenAIGPT4Dot1Mini20250414:           1,
				consts.OpenAIGPT4Dot1Nano:                   1,
				consts.OpenAIGPT4Dot1Nano20250414:           1,
				consts.OpenAIGPT4Dot5Preview:                1,
				consts.OpenAIGPT4Dot5Preview20250227:        1,
				consts.OpenAIGPT3Dot5Turbo0125:              0,
				consts.OpenAIGPT3Dot5Turbo1106:              0,
				consts.OpenAIGPT3Dot5Turbo0613:              0,
				consts.OpenAIGPT3Dot5Turbo0301:              0,
				consts.OpenAIGPT3Dot5Turbo16k:               0,
				consts.OpenAIGPT3Dot5Turbo16K0613:           0,
				consts.OpenAIGPT3Dot5Turbo:                  0,
				consts.OpenAIGPT3Dot5TurboInstruct:          0,
				consts.OpenAIGPT3Dot5TurboInstruct0914:      0,
				consts.OpenAIDavinci002:                     0,
				consts.OpenAIBabbage002:                     0,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                1,
				consts.OpenAIGPT4oAudioPreview20241001:        1,
				consts.OpenAIGPT4oAudioPreview20241217:        1,
				consts.OpenAIGPT4oAudioPreview20250603:        1,
				consts.OpenAIGPT4oRealtimePreview:             1,
				consts.OpenAIGPT4oRealtimePreview20241001:     1,
				consts.OpenAIGPT4oRealtimePreview20241217:     1,
				consts.OpenAIGPT4oRealtimePreview20250603:     1,
				consts.OpenAIGPT4oMiniAudioPreview:            1,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    1,
				consts.OpenAIGPT4oMiniRealtimePreview:         1,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: 1,
			},
			consts.ImageModel: {
				// image
				consts.OpenAIDallE2:    1,
				consts.OpenAIDallE3:    1,
				consts.OpenAIGPTImage1: 1,
			},
			consts.AudioModel: {
				// audio
				consts.OpenAITTS1:                1,
				consts.OpenAITTS1_1106:           1,
				consts.OpenAITTS1HD:              1,
				consts.OpenAITTS1HD1106:          1,
				consts.OpenAIWhisper1:            1,
				consts.OpenAIGPT4oTranscribe:     1,
				consts.OpenAIGPT4oMiniTranscribe: 1,
				consts.OpenAIGPT4oMiniTTS:        1,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                1,
				consts.OpenAIGPT4oAudioPreview20241001:        1,
				consts.OpenAIGPT4oAudioPreview20241217:        1,
				consts.OpenAIGPT4oAudioPreview20250603:        1,
				consts.OpenAIGPT4oRealtimePreview:             1,
				consts.OpenAIGPT4oRealtimePreview20241001:     1,
				consts.OpenAIGPT4oRealtimePreview20241217:     1,
				consts.OpenAIGPT4oRealtimePreview20250603:     1,
				consts.OpenAIGPT4oMiniAudioPreview:            1,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    1,
				consts.OpenAIGPT4oMiniRealtimePreview:         1,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: 1,
			},
			// moderation
			consts.ModerationModel: {
				consts.OpenAIOmniModerationLatest:   0,
				consts.OpenAIOmniModeration20240926: 0,
			},
			// embed
			consts.EmbedModel: {
				consts.OpenAITextEmbedding3Small: 0,
				consts.OpenAITextEmbedding3Large: 0,
				consts.OpenAITextEmbeddingAda002: 0,
			},
		},
	}
	core.RegisterProvider(consts.OpenAI, openaiService)
}

// GetSupportedModels 获取支持的模型
func (s *openAIProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]uint) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *openAIProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *openAIProvider) ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Provider: consts.OpenAI,
		Method:   http.MethodGet,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiModels,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
	})
	return
}
