/*
* @Author: liusuxian 382185882@qq.com
* @Date: 2025-04-11 10:34:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 10:31:50
* @Description:
*
* Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
*/
package aisdk

import (
	"fmt"
	"sync"
)

// ProviderFactory 管理所有AI服务提供商的工厂
type ProviderFactory struct {
	providers map[Provider]ProviderService // 所有提供商
	mu        sync.RWMutex                 // 读写锁
}

// NewProviderFactory 创建 ProviderFactory
func NewProviderFactory() (factory *ProviderFactory) {
	factory = &ProviderFactory{
		providers: make(map[Provider]ProviderService),
	}
	return
}

// RegisterProvider 注册提供商
func (f *ProviderFactory) RegisterProvider(service ProviderService) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.providers[service.GetProvider()] = service
}

// GetProvider 获取提供商
func (f *ProviderFactory) GetProvider(provider Provider) (service ProviderService, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if p, ok := f.providers[provider]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider %s not registered", provider)
}

// ListProviders 列出所有注册的提供商
func (f *ProviderFactory) ListProviders() (providers []Provider) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	providers = make([]Provider, 0, len(f.providers))
	for p := range f.providers {
		providers = append(providers, p)
	}
	return
}
