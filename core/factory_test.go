/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 19:17:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-26 17:16:30
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core_test

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/models"
	"reflect"
	"testing"
)

// MockProviderService 用于测试的模拟提供商实现
type MockProviderService struct {
}

// GetSupportedModels 获取支持的模型
func (m *MockProviderService) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return
}

// InitializeProviderConfig 初始化提供商配置
func (m *MockProviderService) InitializeProviderConfig(config *conf.ProviderConfig) {
}

// InitializeConnectionOptions 初始化连接选项
func (m *MockProviderService) InitializeConnectionOptions(options *conf.ConnectionOptions) {
}

// CreateChatCompletion
func (m *MockProviderService) CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error) {
	return models.ChatResponse{}, nil
}

func TestProviderFactory(t *testing.T) {
	// 创建模拟服务提供商
	openaiService := &MockProviderService{}
	claudeService := &MockProviderService{}
	geminiService := &MockProviderService{}
	// 测试 RegisterProvider 方法
	core.RegisterProvider(consts.OpenAI, openaiService)
	core.RegisterProvider(consts.Claude, claudeService)
	core.RegisterProvider(consts.Gemini, geminiService)
	// 测试 GetProvider 方法
	service, err := core.GetProvider(consts.OpenAI)
	if err != nil {
		t.Fatalf("GetProvider(OpenAI) returned error: %v", err)
	}
	if service != openaiService {
		t.Fatalf("GetProvider(OpenAI) returned wrong service")
	}
	service, err = core.GetProvider(consts.Claude)
	if err != nil {
		t.Fatalf("GetProvider(Claude) returned error: %v", err)
	}
	if service != claudeService {
		t.Fatalf("GetProvider(Claude) returned wrong service")
	}
	service, err = core.GetProvider(consts.Gemini)
	if err != nil {
		t.Fatalf("GetProvider(Gemini) returned error: %v", err)
	}
	if service != geminiService {
		t.Fatalf("GetProvider(Gemini) returned wrong service")
	}
	// 测试获取未注册的提供商
	_, err = core.GetProvider(consts.Provider("UnknownProvider"))
	if err == nil {
		t.Fatalf("GetProvider(UnknownProvider) should return error")
	}
	// 测试 ListProviders 方法
	providers := core.ListProviders()
	expectedProviders := []consts.Provider{consts.OpenAI, consts.Claude, consts.Gemini}
	if len(providers) != len(expectedProviders) {
		t.Fatalf("ListProviders() returned %d providers, expected %d", len(providers), len(expectedProviders))
	}
	// 检查所有预期的提供商是否都在返回的列表中
	providerMap := make(map[consts.Provider]bool)
	for _, p := range providers {
		providerMap[p] = true
	}
	for _, expected := range expectedProviders {
		if !providerMap[expected] {
			t.Fatalf("ListProviders() did not return expected provider: %s", expected)
		}
	}
	// 检查所有预期的提供商是否都在返回的列表中
	providerSet1 := make(map[consts.Provider]struct{})
	for _, p := range providers {
		providerSet1[p] = struct{}{}
	}
	providerSet2 := make(map[consts.Provider]struct{})
	for _, p := range expectedProviders {
		providerSet2[p] = struct{}{}
	}
	if !reflect.DeepEqual(providerSet1, providerSet2) {
		t.Fatalf("ListProviders() returned %v, expected %v", providers, expectedProviders)
	}
}
