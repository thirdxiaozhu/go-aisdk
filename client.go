/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 19:21:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"github.com/liusuxian/aisdk/conf"
	"github.com/liusuxian/aisdk/core"
)

// SDKClient SDK客户端
type SDKClient struct {
	// factory     *core.AIClientFactory
	// config      *config.Manager
	// options     core.ClientOptions
	// initialized bool
	configManager *conf.SDKConfigManager // 配置管理器
}

// NewSDKClient 创建一个SDK客户端
func NewSDKClient(configPath string) (client *SDKClient, err error) {
	// 创建SDK配置管理器
	var configManager *conf.SDKConfigManager
	if configManager, err = conf.NewSDKConfigManager(configPath); err != nil {
		return
	}
	// 创建 ProviderFactory
	providerFactory := core.NewProviderFactory()
	// 创建SDK客户端
	client = &SDKClient{
		configManager: configManager,
	}
	// 获取整个配置
	config := client.configManager.GetConfig()
	// 注册所有提供商
	for providerName := range config.Providers {
		providerFactory.RegisterProvider(providerName)
	}
	return
}
