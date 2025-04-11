/*
* @Author: liusuxian 382185882@qq.com
* @Date: 2025-04-11 10:34:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-11 11:22:31
* @Description:
*
* Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
*/
package aisdk_test

import (
	"context"
	"github.com/liusuxian/aisdk"
	"reflect"
	"testing"
)

// MockProviderService 用于测试的模拟提供商实现
type MockProviderService struct {
	provider aisdk.Provider
}

// GetProvider 获取提供商
func (m *MockProviderService) GetProvider() (provider aisdk.Provider) {
	return m.provider
}

// CreateChatCompletion 聊天相关
func (m *MockProviderService) CreateChatCompletion(ctx context.Context, request aisdk.ChatRequest) (response aisdk.ChatResponse, err error) {
	return aisdk.ChatResponse{}, nil
}

func TestProviderFactory(t *testing.T) {
	// 创建工厂实例
	factory := aisdk.NewProviderFactory()
	// 创建模拟服务提供商
	openaiService := &MockProviderService{provider: aisdk.OpenAI}
	claudeService := &MockProviderService{provider: aisdk.Claude}
	geminiService := &MockProviderService{provider: aisdk.Gemini}
	// 测试 RegisterProvider 方法
	factory.RegisterProvider(openaiService)
	factory.RegisterProvider(claudeService)
	factory.RegisterProvider(geminiService)
	// 测试 GetProvider 方法
	service, err := factory.GetProvider(aisdk.OpenAI)
	if err != nil {
		t.Fatalf("GetProvider(OpenAI) returned error: %v", err)
	}
	if service != openaiService {
		t.Fatalf("GetProvider(OpenAI) returned wrong service")
	}
	service, err = factory.GetProvider(aisdk.Claude)
	if err != nil {
		t.Fatalf("GetProvider(Claude) returned error: %v", err)
	}
	if service != claudeService {
		t.Fatalf("GetProvider(Claude) returned wrong service")
	}
	service, err = factory.GetProvider(aisdk.Gemini)
	if err != nil {
		t.Fatalf("GetProvider(Gemini) returned error: %v", err)
	}
	if service != geminiService {
		t.Fatalf("GetProvider(Gemini) returned wrong service")
	}
	// 测试获取未注册的提供商
	_, err = factory.GetProvider(aisdk.Provider("UnknownProvider"))
	if err == nil {
		t.Fatalf("GetProvider(UnknownProvider) should return error")
	}
	// 测试 ListProviders 方法
	providers := factory.ListProviders()
	expectedProviders := []aisdk.Provider{aisdk.OpenAI, aisdk.Claude, aisdk.Gemini}
	if len(providers) != len(expectedProviders) {
		t.Fatalf("ListProviders() returned %d providers, expected %d", len(providers), len(expectedProviders))
	}
	// 检查所有预期的提供商是否都在返回的列表中
	providerMap := make(map[aisdk.Provider]bool)
	for _, p := range providers {
		providerMap[p] = true
	}
	for _, expected := range expectedProviders {
		if !providerMap[expected] {
			t.Fatalf("ListProviders() did not return expected provider: %s", expected)
		}
	}
	// 检查所有预期的提供商是否都在返回的列表中
	providerSet1 := make(map[aisdk.Provider]struct{})
	for _, p := range providers {
		providerSet1[p] = struct{}{}
	}
	providerSet2 := make(map[aisdk.Provider]struct{})
	for _, p := range expectedProviders {
		providerSet2[p] = struct{}{}
	}
	if !reflect.DeepEqual(providerSet1, providerSet2) {
		t.Fatalf("ListProviders() returned %v, expected %v", providers, expectedProviders)
	}
}
