/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:45:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 19:15:27
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"context"
	"github.com/liusuxian/aisdk/consts"
	"github.com/liusuxian/aisdk/models"
)

// ProviderService AI服务提供商的服务接口
type ProviderService interface {
	GetProvider() (provider consts.Provider) // 获取提供商

	// 聊天相关
	CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error)
	// CreateChatCompletionStream(ctx context.Context, request models.ChatRequest) (response models.ChatResponseStream, err error)

	// TODO 图像相关

	// TODO 视频相关

	// TODO 音频相关
}
