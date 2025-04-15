/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 19:22:11
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/aisdk/consts"
	"github.com/liusuxian/aisdk/models"
)

// OpenAIProvider OpenAI提供商
type OpenAIProvider struct {
}

// GetProvider 获取提供商
func (p *OpenAIProvider) GetProvider() (provider consts.Provider) {
	return consts.OpenAI
}

// CreateChatCompletion
func (p *OpenAIProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error) {
	return
}
