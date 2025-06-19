/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 18:18:15
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
	core.DefaultProviderService
	supportedModels map[fmt.Stringer]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig             // 提供商配置
	lb              *loadbalancer.LoadBalancer       // 负载均衡器
}

var (
	deepseekService *deepseekProvider // deepseek提供商实例
)

const (
	apiModels = "/models"
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	deepseekService = &deepseekProvider{
		supportedModels: map[fmt.Stringer]map[string]bool{
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
func (s *deepseekProvider) GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *deepseekProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *deepseekProvider) ListModels(ctx context.Context, provider fmt.Stringer, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = s.executeRequest(ctx, http.MethodGet, apiModels, opts, &response)
	return
}
