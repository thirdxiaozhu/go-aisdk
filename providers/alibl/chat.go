/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:47:24
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-02 22:15:16
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package alibl

import (
	"context"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiChatCompletionsText       = "/services/aigc/text-generation/generation"
	apiChatCompletionsMultimodal = "/services/aigc/multimodal-generation/generation"
)

// CreateChatCompletion 创建聊天
func (s *aliblProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, s.apiChatCompletions(request.Model), opts, s.lb, nil, &response, withRequestOptions(request)...)
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *aliblProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	var stream *httpclient.StreamReader[models.ChatBaseResponse]
	if stream, err = common.ExecuteStreamRequest[models.ChatBaseResponse](ctx, http.MethodPost, s.providerConfig.BaseURL, s.apiChatCompletions(request.Model), opts, s.lb, withRequestOptions(request)...); err != nil {
		return
	}
	response = models.ChatResponseStream{
		StreamReader: stream,
	}
	return
}

// apiChatCompletions 获取聊天接口
func (s *aliblProvider) apiChatCompletions(model string) (api string) {
	// 判断模型是否支持多模态
	if s.supportedModels[consts.ChatModel][model] == 1 {
		return apiChatCompletionsMultimodal
	}
	return apiChatCompletionsText
}

// withRequestOptions 添加请求选项
func withRequestOptions(request models.ChatRequest) (reqSetters []httpclient.RequestOption) {
	reqSetters = []httpclient.RequestOption{
		httpclient.WithBody(request),
	}
	if data, ok := request.Metadata["X-DashScope-DataInspection"]; ok {
		reqSetters = append(reqSetters, httpclient.WithKeyValue("X-DashScope-DataInspection", data))
	}
	if models.BoolValue(request.Stream) {
		reqSetters = append(reqSetters, httpclient.WithKeyValue("X-DashScope-SSE", "enable"))
	}
	return
}
