/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:31:10
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 17:07:39
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package ark

import (
	"fmt"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/loadbalancer"
)

// arkProvider 火山大模型提供商
type arkProvider struct {
	core.DefaultProviderService
	supportedModels map[fmt.Stringer]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig             // 提供商配置
	lb              *loadbalancer.LoadBalancer       // 负载均衡器
}

var (
	arkService *arkProvider // ark提供商实例
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	arkService = &arkProvider{
		supportedModels: map[fmt.Stringer]map[string]bool{
			consts.ChatModel: {
				// chat
				consts.Doubaoseed1_6: true,
			},
			consts.ImageModel: {
				consts.Doubaoseedream3: true,
			},
		},
	}
	core.RegisterProvider(consts.Ark, arkService)
}

// GetSupportedModels 获取支持的模型
func (s *arkProvider) GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *arkProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}
