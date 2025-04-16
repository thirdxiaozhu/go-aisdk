/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-16 17:52:16
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/aisdk/conf"
	"github.com/liusuxian/aisdk/consts"
	"github.com/liusuxian/aisdk/core"
	"github.com/liusuxian/aisdk/models"
)

// openAIProvider OpenAI提供商
type openAIProvider struct {
	supportedModels   map[consts.ModelType][]string // 支持的模型
	providerConfig    *conf.ProviderConfig          // 提供商配置
	connectionOptions *conf.ConnectionOptions       // 连接选项
}

var (
	openaiService *openAIProvider // OpenAI提供商实例
)

// init 包初始化时创建 openAIProvider 实例并注册到工厂
func init() {
	openaiService = &openAIProvider{
		supportedModels: map[consts.ModelType][]string{
			consts.ChatModel:  {},
			consts.ImageModel: {},
			consts.VideoModel: {},
			consts.AudioModel: {},
			consts.EmbedModel: {},
		},
	}
	core.RegisterProvider(consts.OpenAI, openaiService)
}

// GetSupportedModels 获取支持的模型
func (s *openAIProvider) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *openAIProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
}

// InitializeConnectionOptions 初始化连接选项
func (s *openAIProvider) InitializeConnectionOptions(options *conf.ConnectionOptions) {
	s.connectionOptions = options
}

// TODO: CreateChatCompletion 创建聊天
func (s *openAIProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error) {
	return
}
