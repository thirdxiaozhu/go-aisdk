/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-18 15:06:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-20 01:32:38
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiChatCompletions = "/chat/completions"
)

// CreateChatCompletion 创建聊天
func (s *openAIProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, apiChatCompletions, opts, s.lb, nil, &response, httpclient.WithBody(request))
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *openAIProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	var stream *httpclient.StreamReader[models.ChatBaseResponse]
	if stream, err = common.ExecuteStreamRequest[models.ChatBaseResponse](ctx, http.MethodPost, s.providerConfig.BaseURL, apiChatCompletions, opts, s.lb, httpclient.WithBody(request)); err != nil {
		return
	}
	response = models.ChatResponseStream{
		StreamReader: stream,
	}
	return
}
