/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:31:10
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 12:58:05
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
	supportedModels map[fmt.Stringer]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig             // 提供商配置
	lb              *loadbalancer.LoadBalancer       // 负载均衡器
}

var (
	aliblService *aliblProvider // alibl提供商实例
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	aliblService = &aliblProvider{
		supportedModels: map[fmt.Stringer]map[string]bool{
			consts.ChatModel: {
				// chat
				consts.AliBLQwqPlus:                       true,
				consts.AliBLQwqPlusLatest:                 true,
				consts.AliBLQwqPlus20250305:               true,
				consts.AliBLQwenMax:                       true,
				consts.AliBLQwenMaxLatest:                 true,
				consts.AliBLQwenMax20250125:               true,
				consts.AliBLQwenMax20240919:               true,
				consts.AliBLQwenMax20240428:               true,
				consts.AliBLQwenMax20240403:               true,
				consts.AliBLQwenPlus:                      true,
				consts.AliBLQwenPlusLatest:                true,
				consts.AliBLQwenPlus20250428:              true,
				consts.AliBLQwenPlus20250125:              true,
				consts.AliBLQwenPlus20250112:              true,
				consts.AliBLQwenPlus20241220:              true,
				consts.AliBLQwenPlus20241127:              true,
				consts.AliBLQwenPlus20241125:              true,
				consts.AliBLQwenPlus20240919:              true,
				consts.AliBLQwenPlus20240806:              true,
				consts.AliBLQwenPlus20240723:              true,
				consts.AliBLQwenTurbo:                     true,
				consts.AliBLQwenTurboLatest:               true,
				consts.AliBLQwenTurbo20250428:             true,
				consts.AliBLQwenTurbo20250211:             true,
				consts.AliBLQwenTurbo20241101:             true,
				consts.AliBLQwenTurbo20240919:             true,
				consts.AliBLQwenTurbo20240624:             true,
				consts.AliBLQwenLong:                      true,
				consts.AliBLQwenLongLatest:                true,
				consts.AliBLQwenLong20250125:              true,
				consts.AliBLQwenOmniTurbo:                 true,
				consts.AliBLQwenOmniTurboLatest:           true,
				consts.AliBLQwenOmniTurbo20250326:         true,
				consts.AliBLQwenOmniTurbo20250119:         true,
				consts.AliBLQwenOmniTurboRealtime:         true,
				consts.AliBLQwenOmniTurboRealtimeLatest:   true,
				consts.AliBLQwenOmniTurboRealtime20250508: true,
				consts.AliBLQvqMax:                        true,
				consts.AliBLQvqMaxLatest:                  true,
				consts.AliBLQvqMax20250515:                true,
				consts.AliBLQvqMax20250325:                true,
				consts.AliBLQvqPlus:                       true,
				consts.AliBLQvqPlusLatest:                 true,
				consts.AliBLQvqPlus20250515:               true,
				consts.AliBLQwenVlMax:                     true,
				consts.AliBLQwenVlMaxLatest:               true,
				consts.AliBLQwenVlMax20250408:             true,
				consts.AliBLQwenVlMax20250402:             true,
				consts.AliBLQwenVlMax20250125:             true,
				consts.AliBLQwenVlMax20241230:             true,
				consts.AliBLQwenVlMax20241119:             true,
				consts.AliBLQwenVlMax20241030:             true,
				consts.AliBLQwenVlMax20240809:             true,
				consts.AliBLQwenVlPlus:                    true,
				consts.AliBLQwenVlPlusLatest:              true,
				consts.AliBLQwenVlPlus20250507:            true,
				consts.AliBLQwenVlPlus20250125:            true,
				consts.AliBLQwenVlPlus20250102:            true,
				consts.AliBLQwenVlPlus20240809:            true,
				consts.AliBLQwenVlPlus20231201:            true,
				consts.AliBLQwenVlOcr:                     true,
				consts.AliBLQwenVlOcrLatest:               true,
				consts.AliBLQwenVlOcr20250413:             true,
				consts.AliBLQwenVlOcr20241028:             true,
				consts.AliBLQwenAudioTurbo:                true,
				consts.AliBLQwenAudioTurboLatest:          true,
				consts.AliBLQwenAudioTurbo20241204:        true,
				consts.AliBLQwenAudioTurbo20240807:        true,
				consts.AliBLQwenAudioAsr:                  true,
				consts.AliBLQwenAudioAsrLatest:            true,
				consts.AliBLQwenAudioAsr20241204:          true,
				consts.AliBLQwenMathPlus:                  true,
				consts.AliBLQwenMathPlusLatest:            true,
				consts.AliBLQwenMathPlus20240919:          true,
				consts.AliBLQwenMathPlus20240816:          true,
				consts.AliBLQwenMathTurbo:                 true,
				consts.AliBLQwenMathTurboLatest:           true,
				consts.AliBLQwenMathTurbo20240919:         true,
				consts.AliBLQwenCoderPlus:                 true,
				consts.AliBLQwenCoderPlusLatest:           true,
				consts.AliBLQwenCoderPlus20241106:         true,
				consts.AliBLQwenCoderTurbo:                true,
				consts.AliBLQwenCoderTurboLatest:          true,
				consts.AliBLQwenCoderTurbo20240919:        true,
				consts.AliBLQwenMtPlus:                    true,
				consts.AliBLQwenMtTurbo:                   true,
				consts.AliBLQwen3_235bA22b:                true,
				consts.AliBLQwen3_32b:                     true,
				consts.AliBLQwen3_30bA3b:                  true,
				consts.AliBLQwen3_14b:                     true,
				consts.AliBLQwen3_8b:                      true,
				consts.AliBLQwen3_4b:                      true,
				consts.AliBLQwen3_17b:                     true,
				consts.AliBLQwen3_06b:                     true,
				consts.AliBLQwq32b:                        true,
				consts.AliBLQwq32bPreview:                 true,
				consts.AliBLQwen2Dot5_14bInstruct1m:       true,
				consts.AliBLQwen2Dot5_7bInstruct1m:        true,
				consts.AliBLQwen2Dot5_72bInstruct:         true,
				consts.AliBLQwen2Dot5_32bInstruct:         true,
				consts.AliBLQwen2Dot5_14bInstruct:         true,
				consts.AliBLQwen2Dot5_7bInstruct:          true,
				consts.AliBLQwen2Dot5_3bInstruct:          true,
				consts.AliBLQwen2Dot5_15bInstruct:         true,
				consts.AliBLQwen2Dot5_05bInstruct:         true,
				consts.AliBLQwen2_72bInstruct:             true,
				consts.AliBLQwen2_57bA14bInstruct:         true,
				consts.AliBLQwen2_7bInstruct:              true,
				consts.AliBLQwen2_15bInstruct:             true,
				consts.AliBLQwen2_05bInstruct:             true,
				consts.AliBLQwen1Dot5_110bChat:            true,
				consts.AliBLQwen1Dot5_72bChat:             true,
				consts.AliBLQwen1Dot5_32bChat:             true,
				consts.AliBLQwen1Dot5_14bChat:             true,
				consts.AliBLQwen1Dot5_7bChat:              true,
				consts.AliBLQwen1Dot5_18bChat:             true,
				consts.AliBLQwen1Dot5_05bChat:             true,
				consts.AliBLQvq72bPreview:                 true,
				consts.AliBLQwen2Dot5Omni7b:               true,
				consts.AliBLQwen2Dot5Vl72bInstruct:        true,
				consts.AliBLQwen2Dot5Vl32bInstruct:        true,
				consts.AliBLQwen2Dot5Vl7bInstruct:         true,
				consts.AliBLQwen2Dot5Vl3bInstruct:         true,
				consts.AliBLQwen2Vl72bInstruct:            true,
				consts.AliBLQwen2Vl7bInstruct:             true,
				consts.AliBLQwen2Vl2bInstruct:             true,
				consts.AliBLQwenVlV1:                      true,
				consts.AliBLQwenVlChatV1:                  true,
				consts.AliBLQwen2AudioInstruct:            true,
				consts.AliBLQwenAudioChat:                 true,
				consts.AliBLQwen2Dot5Math72bInstruct:      true,
				consts.AliBLQwen2Dot5Math7bInstruct:       true,
				consts.AliBLQwen2Dot5Math15bInstruct:      true,
				consts.AliBLQwen2Dot5Coder32bInstruct:     true,
				consts.AliBLQwen2Dot5Coder14bInstruct:     true,
				consts.AliBLQwen2Dot5Coder7bInstruct:      true,
				consts.AliBLQwen2Dot5Coder3bInstruct:      true,
				consts.AliBLQwen2Dot5Coder15bInstruct:     true,
				consts.AliBLQwen2Dot5Coder05bInstruct:     true,
			},
		},
	}
	core.RegisterProvider(consts.AliBL, aliblService)
}

// GetSupportedModels 获取支持的模型
func (s *aliblProvider) GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *aliblProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}
