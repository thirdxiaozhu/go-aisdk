/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 19:09:15
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-28 16:24:02
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
)

var (
	errConfigFileEmpty = errors.New("config file path is empty")
)

// ProviderConfig AI服务提供商的配置
type ProviderConfig struct {
	BaseURL          string            `json:"base_url"`          // 基础URL，用于自定义API服务器的地址
	APIKeys          []string          `json:"api_keys"`          // API密钥列表
	OrgID            string            `json:"org_id"`            // 组织ID，对于某些提供商可能需要
	APIVersion       string            `json:"api_version"`       // API版本，对于某些提供商可能需要
	AssistantVersion string            `json:"assistant_version"` // 助手版本，对于某些提供商可能需要
	Extra            map[string]string `json:"extra"`             // 额外参数，对于某些提供商可能需要
}

// SDKConfig SDK整体配置
type SDKConfig struct {
	Providers map[string]ProviderConfig `json:"providers"` // AI服务提供商的配置
}

// SDKConfigManager SDK配置管理器
type SDKConfigManager struct {
	configPath string    // 配置文件路径
	config     SDKConfig // SDK配置
}

// NewSDKConfigManager 创建SDK配置管理器
func NewSDKConfigManager(configPath string) (manager *SDKConfigManager, err error) {
	if configPath == "" {
		err = errConfigFileEmpty
		return
	}

	manager = &SDKConfigManager{
		configPath: configPath,
		config: SDKConfig{
			Providers: make(map[string]ProviderConfig),
		},
	}
	// 尝试加载配置
	if _, err = os.Stat(configPath); err == nil {
		if err = manager.Load(); err != nil {
			err = fmt.Errorf("failed to load config: %w", err)
			return
		}
		return
	}

	if !os.IsNotExist(err) {
		err = fmt.Errorf("failed to check config file: %w", err)
		return
	}
	return
}

// Load 从文件加载配置
func (m *SDKConfigManager) Load() (err error) {
	var data []byte
	if data, err = os.ReadFile(m.configPath); err != nil {
		err = fmt.Errorf("failed to read config file: %w", err)
		return
	}

	if err = json.Unmarshal(data, &m.config); err != nil {
		err = fmt.Errorf("failed to unmarshal config: %w", err)
		return
	}
	return
}

// GetConfig 获取整个配置
func (m *SDKConfigManager) GetConfig() (configCopy SDKConfig) {
	// 返回配置的副本，防止外部修改
	configCopy = SDKConfig{
		Providers: make(map[string]ProviderConfig),
	}

	for k, v := range m.config.Providers {
		configCopy.Providers[k] = cloneProviderConfig(v)
	}
	return
}

// GetProviderConfig 获取提供商配置
func (m *SDKConfigManager) GetProviderConfig(provider string) (config ProviderConfig) {
	if cfg, ok := m.config.Providers[provider]; ok {
		config = cloneProviderConfig(cfg)
		return
	}

	config = ProviderConfig{}
	return
}

// cloneProviderConfig 深拷贝 ProviderConfig
func cloneProviderConfig(source ProviderConfig) (dest ProviderConfig) {
	var extraCopy map[string]string
	if source.Extra != nil {
		extraCopy = make(map[string]string)
		maps.Copy(extraCopy, source.Extra)
	}

	dest = ProviderConfig{
		BaseURL:          source.BaseURL,
		APIKeys:          slices.Clone(source.APIKeys),
		OrgID:            source.OrgID,
		APIVersion:       source.APIVersion,
		AssistantVersion: source.AssistantVersion,
		Extra:            extraCopy,
	}
	return
}
