/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:47:24
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 19:19:43
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package alibl

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
	"strings"
)

const (
	apiChatCompletionsText       = "/services/aigc/text-generation/generation"
	apiChatCompletionsMultimodal = "/services/aigc/multimodal-generation/generation"
)

// AliBLChatBaseResponse 阿里百炼聊天响应基础信息
type AliBLChatBaseResponse struct {
	RequestId string `json:"request_id,omitempty"` // 本次调用的唯一标识符
	Output    *struct {
		Choices []struct {
			FinishReason models.ChatFinishReason `json:"finish_reason,omitempty"` // 模型停止生成 token 的原因
			Message      *struct {
				Role             string             `json:"role,omitempty"`              // 输出消息的角色
				Content          any                `json:"content,omitempty"`           // 输出消息的内容
				ReasoningContent string             `json:"reasoning_content,omitempty"` // 模型的深度思考内容
				ToolCalls        []models.ToolCalls `json:"tool_calls,omitempty"`        // 工具调用
			} `json:"message,omitempty"` // 模型输出的消息对象
			LogProbs *models.ChatLogProbs `json:"logprobs,omitempty"` // 当前 choices 对象的概率信息
		} `json:"choices,omitempty"` // 模型的输出信息
		SearchInfo *struct {
			SearchResults []struct {
				SiteName string `json:"site_name,omitempty"` // 搜索结果来源的网站名称
				Icon     string `json:"icon,omitempty"`      // 来源网站的图标URL，如果没有图标则为空字符串
				Index    int    `json:"index,omitempty"`     // 搜索结果的序号，表示该搜索结果在search_results中的索引
				Title    string `json:"title,omitempty"`     // 搜索结果的标题
				URL      string `json:"url,omitempty"`       // 搜索结果的链接地址
			} `json:"search_results,omitempty"` // 联网搜索到的结果
		} `json:"search_info,omitempty"` // 联网搜索到的信息，在设置search_options参数后会返回该参数
	} `json:"output,omitempty"` // 调用结果信息
	Usage *struct {
		InputTokens        int `json:"input_tokens,omitempty"`  // 用户输入内容转换成Token后的长度
		OutputTokens       int `json:"output_tokens,omitempty"` // 模型输出内容转换成Token后的长度
		InputTokensDetails *struct {
			TextTokens  int `json:"text_tokens,omitempty"`  // 输入的文本转换为Token后的长度
			ImageTokens int `json:"image_tokens,omitempty"` // 输入的图像转换为Token后的长度
			VideoTokens int `json:"video_tokens,omitempty"` // 输入的视频文件或图像列表转换为Token后的长度
		} `json:"input_tokens_details,omitempty"` // 输入内容转换成Token后的长度详情
		TotalTokens         int `json:"total_tokens,omitempty"` // 当输入为纯文本时返回该字段，为input_tokens与output_tokens之和
		ImageTokens         int `json:"image_tokens,omitempty"` // 输入内容包含image时返回该字段。为用户输入图片内容转换成Token后的长度
		VideoTokens         int `json:"video_tokens,omitempty"` // 输入内容包含video时返回该字段。为用户输入视频内容转换成Token后的长度
		AudioTokens         int `json:"audio_tokens,omitempty"` // 输入内容包含audio时返回该字段。为用户输入音频内容转换成Token后的长度
		OutputTokensDetails *struct {
			TextTokens      int `json:"text_tokens,omitempty"`      // 输出的文本转换为Token后的长度
			ReasoningTokens int `json:"reasoning_tokens,omitempty"` // 模型思考过程转换为Token后的长度
		} `json:"output_tokens_details,omitempty"` // 输出内容转换成 Token后的长度详情
		PromptTokensDetails *models.PromptTokensDetails `json:"prompt_tokens_details,omitempty"` // prompt tokens 的详细信息
	} `json:"usage,omitempty"` // 本次调用的 token 使用情况
}

// AliBLChatResponse 阿里百炼聊天响应
type AliBLChatResponse struct {
	AliBLChatBaseResponse
	httpclient.HttpHeader
}

// CreateChatCompletion 创建聊天
func (s *aliblProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	aliblResponse := AliBLChatResponse{}
	if err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, s.apiChatCompletions(request.Model), opts, s.lb, nil, &aliblResponse, withRequestOptions(request)...); err != nil {
		return
	}
	fmt.Printf("aliblResponse = %s\n", httpclient.MustString(aliblResponse))
	// 转换阿里百炼响应为 SDK 通用响应
	response = convertResponse(&aliblResponse)
	return
}

// TODO CreateChatCompletionStream 创建流式聊天
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
	if request.XDashScopeDataInspection != nil {
		reqSetters = append(reqSetters, httpclient.WithKeyValue("X-DashScope-DataInspection", httpclient.MustString(request.XDashScopeDataInspection)))
	}
	if models.BoolValue(request.Stream) {
		reqSetters = append(reqSetters, httpclient.WithKeyValue("X-DashScope-SSE", "enable"))
	}
	return
}

// convertResponse 转换阿里百炼响应为 SDK 通用响应
func convertResponse(resp *AliBLChatResponse) (response models.ChatResponse) {
	response = models.ChatResponse{
		ChatBaseResponse: models.ChatBaseResponse{},
		HttpHeader:       resp.HttpHeader,
	}
	if resp.Output != nil {
		if len(resp.Output.Choices) > 0 {
			choices := make([]models.ChatChoice, len(resp.Output.Choices))
			for i, choice := range resp.Output.Choices {
				choices[i] = models.ChatChoice{
					FinishReason: choice.FinishReason,
					Message: &models.ChatCompletionMessage{
						Role: choice.Message.Role,
						Content: func() string {
							if content, ok := choice.Message.Content.([]string); ok {
								return strings.Join(content, "\n") // TODO
							}
							if content, ok := choice.Message.Content.(string); ok {
								return content
							}
							return ""
						}(),
						ReasoningContent: choice.Message.ReasoningContent,
						ToolCalls:        choice.Message.ToolCalls,
					},
					LogProbs: choice.LogProbs,
				}
			}
			response.Choices = choices
		}
	}
	return
}
