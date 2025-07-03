/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-18 15:01:49
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:38:03
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiChatCompletions = "/chat/completions"
)

// CreateChatCompletion 创建聊天
func (s *deepseekProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Provider: consts.DeepSeek,
		Method:   http.MethodPost,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiChatCompletions,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
		ReqSetters: []httpclient.RequestOption{
			httpclient.WithBody(request),
		},
	})
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *deepseekProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	var stream *httpclient.StreamReader[models.ChatBaseResponse]
	if stream, err = common.ExecuteStreamRequest[models.ChatBaseResponse](ctx, &common.ExecuteRequestContext{
		Provider: consts.DeepSeek,
		Method:   http.MethodPost,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiChatCompletions,
		Opts:     opts,
		LB:       s.lb,
		ReqSetters: []httpclient.RequestOption{
			httpclient.WithBody(request),
		},
	}); err != nil {
		return
	}
	response = models.ChatResponseStream{
		StreamReader: stream,
	}
	return
}
