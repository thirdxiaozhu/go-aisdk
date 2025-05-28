/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-28 17:07:17
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/models"
	_ "github.com/liusuxian/go-aisdk/providers"
)

// SDKClient SDK客户端
type SDKClient struct {
	configManager *conf.SDKConfigManager // 配置管理器
}

// NewSDKClient 创建一个SDK客户端
func NewSDKClient(configPath string) (client *SDKClient, err error) {
	// 创建SDK配置管理器
	var configManager *conf.SDKConfigManager
	if configManager, err = conf.NewSDKConfigManager(configPath); err != nil {
		return
	}
	// 初始化所有提供商
	for _, provider := range core.ListProviders() {
		// 获取提供商
		var ps core.ProviderService
		if ps, err = core.GetProvider(provider); err != nil {
			return
		}
		// 获取提供商配置
		providerConfig := configManager.GetProviderConfig(provider)
		// 初始化提供商配置
		ps.InitializeProviderConfig(&providerConfig)
	}
	// 创建SDK客户端
	client = &SDKClient{
		configManager: configManager,
	}
	return
}

// ListModels 列出模型
func (c *SDKClient) ListModels(ctx context.Context, provider consts.Provider) (response models.ListModelsResponse, err error) {
	// 获取提供商
	var ps core.ProviderService
	if ps, err = core.GetProvider(provider); err != nil {
		return
	}
	// 列出模型
	return ps.ListModels(ctx)
}

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error) {
	// 获取提供商
	var ps core.ProviderService
	if ps, err = core.GetProvider(request.ModelInfo.Provider); err != nil {
		return
	}
	// 判断模型是否支持
	if err = core.IsModelSupported(ps, request.ModelInfo); err != nil {
		return
	}
	if request.Stream {
		// TODO
		return
	}
	// 创建聊天
	return ps.CreateChatCompletion(ctx, request)
}
