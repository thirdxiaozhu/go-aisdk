/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:31:10
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 17:41:31
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package alibl

import (
	"fmt"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/loadbalancer"
)

// aliblProvider AliBL提供商
type aliblProvider struct {
	core.DefaultProviderService
	supportedModels map[fmt.Stringer]map[string]uint // 支持的模型
	providerConfig  *conf.ProviderConfig             // 提供商配置
	lb              *loadbalancer.LoadBalancer       // 负载均衡器
}

var (
	aliblService *aliblProvider // alibl提供商实例
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	aliblService = &aliblProvider{
		supportedModels: map[fmt.Stringer]map[string]uint{
			consts.ChatModel: {
				// chat
				consts.AliBLQwqPlus:                       0,
				consts.AliBLQwqPlusLatest:                 0,
				consts.AliBLQwqPlus20250305:               0,
				consts.AliBLQwenMax:                       0,
				consts.AliBLQwenMaxLatest:                 0,
				consts.AliBLQwenMax20250125:               0,
				consts.AliBLQwenMax20240919:               0,
				consts.AliBLQwenMax20240428:               0,
				consts.AliBLQwenMax20240403:               0,
				consts.AliBLQwenPlus:                      0,
				consts.AliBLQwenPlusLatest:                0,
				consts.AliBLQwenPlus20250428:              0,
				consts.AliBLQwenPlus20250125:              0,
				consts.AliBLQwenPlus20250112:              0,
				consts.AliBLQwenPlus20241220:              0,
				consts.AliBLQwenPlus20241127:              0,
				consts.AliBLQwenPlus20241125:              0,
				consts.AliBLQwenPlus20240919:              0,
				consts.AliBLQwenPlus20240806:              0,
				consts.AliBLQwenPlus20240723:              0,
				consts.AliBLQwenTurbo:                     0,
				consts.AliBLQwenTurboLatest:               0,
				consts.AliBLQwenTurbo20250428:             0,
				consts.AliBLQwenTurbo20250211:             0,
				consts.AliBLQwenTurbo20241101:             0,
				consts.AliBLQwenTurbo20240919:             0,
				consts.AliBLQwenTurbo20240624:             0,
				consts.AliBLQwenLong:                      0,
				consts.AliBLQwenLongLatest:                0,
				consts.AliBLQwenLong20250125:              0,
				consts.AliBLQwenOmniTurbo:                 0,
				consts.AliBLQwenOmniTurboLatest:           0,
				consts.AliBLQwenOmniTurbo20250326:         0,
				consts.AliBLQwenOmniTurbo20250119:         0,
				consts.AliBLQwenOmniTurboRealtime:         0,
				consts.AliBLQwenOmniTurboRealtimeLatest:   0,
				consts.AliBLQwenOmniTurboRealtime20250508: 0,
				consts.AliBLQvqMax:                        0,
				consts.AliBLQvqMaxLatest:                  0,
				consts.AliBLQvqMax20250515:                0,
				consts.AliBLQvqMax20250325:                0,
				consts.AliBLQvqPlus:                       0,
				consts.AliBLQvqPlusLatest:                 0,
				consts.AliBLQvqPlus20250515:               0,
				consts.AliBLQwenVlMax:                     1,
				consts.AliBLQwenVlMaxLatest:               1,
				consts.AliBLQwenVlMax20250408:             1,
				consts.AliBLQwenVlMax20250402:             1,
				consts.AliBLQwenVlMax20250125:             1,
				consts.AliBLQwenVlMax20241230:             1,
				consts.AliBLQwenVlMax20241119:             1,
				consts.AliBLQwenVlMax20241030:             1,
				consts.AliBLQwenVlMax20240809:             1,
				consts.AliBLQwenVlPlus:                    1,
				consts.AliBLQwenVlPlusLatest:              1,
				consts.AliBLQwenVlPlus20250507:            1,
				consts.AliBLQwenVlPlus20250125:            1,
				consts.AliBLQwenVlPlus20250102:            1,
				consts.AliBLQwenVlPlus20240809:            1,
				consts.AliBLQwenVlPlus20231201:            1,
				consts.AliBLQwenVlOcr:                     1,
				consts.AliBLQwenVlOcrLatest:               1,
				consts.AliBLQwenVlOcr20250413:             1,
				consts.AliBLQwenVlOcr20241028:             1,
				consts.AliBLQwenAudioTurbo:                1,
				consts.AliBLQwenAudioTurboLatest:          1,
				consts.AliBLQwenAudioTurbo20241204:        1,
				consts.AliBLQwenAudioTurbo20240807:        1,
				consts.AliBLQwenAudioAsr:                  1,
				consts.AliBLQwenAudioAsrLatest:            1,
				consts.AliBLQwenAudioAsr20241204:          1,
				consts.AliBLQwenMathPlus:                  0,
				consts.AliBLQwenMathPlusLatest:            0,
				consts.AliBLQwenMathPlus20240919:          0,
				consts.AliBLQwenMathPlus20240816:          0,
				consts.AliBLQwenMathTurbo:                 0,
				consts.AliBLQwenMathTurboLatest:           0,
				consts.AliBLQwenMathTurbo20240919:         0,
				consts.AliBLQwenCoderPlus:                 0,
				consts.AliBLQwenCoderPlusLatest:           0,
				consts.AliBLQwenCoderPlus20241106:         0,
				consts.AliBLQwenCoderTurbo:                0,
				consts.AliBLQwenCoderTurboLatest:          0,
				consts.AliBLQwenCoderTurbo20240919:        0,
				consts.AliBLQwenMtPlus:                    0,
				consts.AliBLQwenMtTurbo:                   0,
				consts.AliBLQwen3_235bA22b:                0,
				consts.AliBLQwen3_32b:                     0,
				consts.AliBLQwen3_30bA3b:                  0,
				consts.AliBLQwen3_14b:                     0,
				consts.AliBLQwen3_8b:                      0,
				consts.AliBLQwen3_4b:                      0,
				consts.AliBLQwen3_17b:                     0,
				consts.AliBLQwen3_06b:                     0,
				consts.AliBLQwq32b:                        0,
				consts.AliBLQwq32bPreview:                 0,
				consts.AliBLQwen2Dot5_14bInstruct1m:       0,
				consts.AliBLQwen2Dot5_7bInstruct1m:        0,
				consts.AliBLQwen2Dot5_72bInstruct:         0,
				consts.AliBLQwen2Dot5_32bInstruct:         0,
				consts.AliBLQwen2Dot5_14bInstruct:         0,
				consts.AliBLQwen2Dot5_7bInstruct:          0,
				consts.AliBLQwen2Dot5_3bInstruct:          0,
				consts.AliBLQwen2Dot5_15bInstruct:         0,
				consts.AliBLQwen2Dot5_05bInstruct:         0,
				consts.AliBLQwen2_72bInstruct:             0,
				consts.AliBLQwen2_57bA14bInstruct:         0,
				consts.AliBLQwen2_7bInstruct:              0,
				consts.AliBLQwen2_15bInstruct:             0,
				consts.AliBLQwen2_05bInstruct:             0,
				consts.AliBLQwen1Dot5_110bChat:            0,
				consts.AliBLQwen1Dot5_72bChat:             0,
				consts.AliBLQwen1Dot5_32bChat:             0,
				consts.AliBLQwen1Dot5_14bChat:             0,
				consts.AliBLQwen1Dot5_7bChat:              0,
				consts.AliBLQwen1Dot5_18bChat:             0,
				consts.AliBLQwen1Dot5_05bChat:             0,
				consts.AliBLQvq72bPreview:                 0,
				consts.AliBLQwen2Dot5Omni7b:               0,
				consts.AliBLQwen2Dot5Vl72bInstruct:        1,
				consts.AliBLQwen2Dot5Vl32bInstruct:        1,
				consts.AliBLQwen2Dot5Vl7bInstruct:         1,
				consts.AliBLQwen2Dot5Vl3bInstruct:         1,
				consts.AliBLQwen2Vl72bInstruct:            1,
				consts.AliBLQwen2Vl7bInstruct:             1,
				consts.AliBLQwen2Vl2bInstruct:             1,
				consts.AliBLQwenVlV1:                      1,
				consts.AliBLQwenVlChatV1:                  1,
				consts.AliBLQwen2AudioInstruct:            1,
				consts.AliBLQwenAudioChat:                 1,
				consts.AliBLQwen2Dot5Math72bInstruct:      0,
				consts.AliBLQwen2Dot5Math7bInstruct:       0,
				consts.AliBLQwen2Dot5Math15bInstruct:      0,
				consts.AliBLQwen2Dot5Coder32bInstruct:     0,
				consts.AliBLQwen2Dot5Coder14bInstruct:     0,
				consts.AliBLQwen2Dot5Coder7bInstruct:      0,
				consts.AliBLQwen2Dot5Coder3bInstruct:      0,
				consts.AliBLQwen2Dot5Coder15bInstruct:     0,
				consts.AliBLQwen2Dot5Coder05bInstruct:     0,
				consts.AliBLDeepSeekR1:                    0,
				consts.AliBLDeepSeekR1_0528:               0,
				consts.AliBLDeepSeekV3:                    0,
				consts.AliBLDeepSeekR1DistillQwen15b:      0,
				consts.AliBLDeepSeekR1DistillQwen7b:       0,
				consts.AliBLDeepSeekR1DistillQwen14b:      0,
				consts.AliBLDeepSeekR1DistillQwen32b:      0,
				consts.AliBLDeepSeekR1DistillLlama8b:      0,
				consts.AliBLDeepSeekR1DistillLlama70b:     0,
				consts.AliBLLlama3Dot3_70bInstruct:        0,
				consts.AliBLLlama3Dot2_3bInstruct:         0,
				consts.AliBLLlama3Dot2_1bInstruct:         0,
				consts.AliBLLlama3Dot1_405bInstruct:       0,
				consts.AliBLLlama3Dot1_70bInstruct:        0,
				consts.AliBLLlama3Dot1_8bInstruct:         0,
				consts.AliBLLlama3_70bInstruct:            0,
				consts.AliBLLlama3_8bInstruct:             0,
				consts.AliBLLlama2_13bChatV2:              0,
				consts.AliBLLlama2_7bChatV2:               0,
				consts.AliBLLlama4Scout17b16eInstruct:     0,
				consts.AliBLLlama4Maverick17b128eInstruct: 0,
				consts.AliBLLlama3Dot2_90bVisionInstruct:  0,
				consts.AliBLLlama3Dot2_11bVision:          0,
				consts.AliBLBaichuan2Turbo:                0,
				consts.AliBLBaichuan2_13bChatV1:           0,
				consts.AliBLBaichuan2_7bChatV1:            0,
				consts.AliBLBaichuan7bV1:                  0,
				consts.AliBLChatglm3_6b:                   0,
				consts.AliBLChatglm6bV2:                   0,
				consts.AliBLYiLarge:                       0,
				consts.AliBLYiMedium:                      0,
				consts.AliBLYiLargeRag:                    0,
				consts.AliBLYiLargeTurbo:                  0,
				consts.AliBLAbab6Dot5gChat:                0,
				consts.AliBLAbab6Dot5tChat:                0,
				consts.AliBLAbab6Dot5sChat:                0,
				consts.AliBLZiyaLlama13bV1:                0,
				consts.AliBLBelleLlama13b2mV1:             0,
				consts.AliBLChatyuanLargeV2:               0,
				consts.AliBLBilla7bSftV1:                  0,
			},
		},
	}
	core.RegisterProvider(consts.AliBL, aliblService)
}

// GetSupportedModels 获取支持的模型
func (s *aliblProvider) GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]uint) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *aliblProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}
