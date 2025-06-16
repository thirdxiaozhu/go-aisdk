/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:45:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-13 19:25:56
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
)

type StreamCallback func(ctx context.Context, response models.ChatResponse) error

// ProviderService AI服务提供商的服务接口
type ProviderService interface {
	CheckRequestValidation(request models.Request) error
	// 获取支持的模型
	GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]bool)
	// 初始化提供商配置
	InitializeProviderConfig(config *conf.ProviderConfig)

	// 列出模型
	ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error)
	// 聊天相关
	CreateChatCompletion(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error)
	CreateChatCompletionStream(ctx context.Context, request models.Request, cb StreamCallback, opts ...httpclient.HTTPClientOption) (httpclient.Response, error)

	// TODO 图像相关

	CreateImageGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (response httpclient.Response, err error)

	// TODO 视频相关

	CreateVideoGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (httpclient.Response, error)

	// TODO 音频相关
}
