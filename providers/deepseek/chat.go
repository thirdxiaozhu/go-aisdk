/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-18 15:01:49
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-18 15:06:51
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"net/http"
)

const (
	apiChatCompletions = "/chat/completions"
)

// CreateChatCompletion 创建聊天
func (s *deepseekProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = s.executeRequest(ctx, http.MethodPost, apiChatCompletions, opts, &response, httpclient.WithBody(request))
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *deepseekProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	var stream *httpclient.StreamReader[models.ChatBaseResponse]
	if stream, err = executeStreamRequest[models.ChatBaseResponse](s, ctx, http.MethodPost, apiChatCompletions, opts, httpclient.WithBody(request)); err != nil {
		return
	}
	response = models.ChatResponseStream{
		StreamReader: stream,
	}
	return
}
