/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:45:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-28 17:06:43
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/models"
	"slices"
)

// ProviderService AI服务提供商的服务接口
type ProviderService interface {
	// 获取支持的模型
	GetSupportedModels() (supportedModels map[consts.ModelType][]string)
	// 初始化提供商配置
	InitializeProviderConfig(config *conf.ProviderConfig)

	// 列出模型
	ListModels(ctx context.Context) (response models.ListModelsResponse, err error)
	// 聊天相关
	CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error)
	// CreateChatCompletionStream(ctx context.Context, request models.ChatRequest) (response models.ChatResponseStream, err error)

	// TODO 图像相关

	// TODO 视频相关

	// TODO 音频相关
}

// IsModelSupported 判断模型是否支持
func IsModelSupported(s ProviderService, modelInfo models.ModelInfo) (err error) {
	// 获取支持的模型
	supportedModels := s.GetSupportedModels()
	if supportedModels == nil {
		return fmt.Errorf("error: provider [%s] has no supported models list", modelInfo.Provider)
	}
	// 获取指定模型类型支持的模型列表
	var (
		modelList []string
		ok        bool
	)
	if modelList, ok = supportedModels[modelInfo.ModelType]; !ok {
		return fmt.Errorf("error: provider [%s] does not support model type [%s]", modelInfo.Provider, modelInfo.ModelType)
	}
	// 判断模型是否支持
	if slices.Contains(modelList, modelInfo.Model) {
		return
	}
	return fmt.Errorf("error: provider [%s] does not support model [%s] of type [%s]", modelInfo.Provider, modelInfo.Model, modelInfo.ModelType)
}
