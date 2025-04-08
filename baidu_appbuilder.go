/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-08 12:19:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 12:41:04
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"fmt"
	"net/http"
)

// BaiduAppBuilderIntegratedReq 请求数据
type BaiduAppBuilderIntegratedReq struct {
	Query          string `json:"query"`                     // 用户的请求query
	ConversationId string `json:"conversation_id,omitempty"` // 对话ID，仅对话型应用生效。在对话型应用中，空:表示表新建会话，非空:表示在对应的会话中继续进行对话，服务内部维护对话历史
}

// BaiduAppBuilderIntegratedRes 响应数据
type BaiduAppBuilderIntegratedRes struct {
	Code    int                              `json:"code"`    // 错误码。非0为错误
	Message string                           `json:"message"` // 报错信息
	Result  *BaiduAppBuilderIntegratedResult `json:"result"`  // 返回结果
	httpHeader
}

// BaiduAppBuilderIntegratedResult 返回结果
type BaiduAppBuilderIntegratedResult struct {
	Answer         string `json:"answer"`         // 应用响应结果
	ConversationId string `json:"conversationId"` // 对话ID，仅对话式应用生效。如果是对话请求中没有conversation_id，则会自动生成一个
}

// BaiduAppBuilderIntegratedResStream 流式响应数据
type BaiduAppBuilderIntegratedResStream struct {
	*streamReader[BaiduAppBuilderIntegratedResult]
}

const (
	baiduAppBuilderIntegratedURL = "/rpc/2.0/cloud_hub/v1/ai_engine/agi_platform/v1/instance/integrated"
)

// BaiduAppBuilderIntegrated 集成API
func (c *Client) BaiduAppBuilderIntegrated(ctx context.Context, request BaiduAppBuilderIntegratedReq) (response *BaiduAppBuilderIntegratedRes, err error) {
	var (
		setters = []requestOption{
			setBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "blocking",
				"conversation_id": request.ConversationId,
			}),
			setKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.config.AuthToken)),
		}
		req *http.Request
	)

	if req, err = c.newRequest(ctx, http.MethodPost, c.fullURL(baiduAppBuilderIntegratedURL), setters...); err != nil {
		return
	}

	response = &BaiduAppBuilderIntegratedRes{}
	err = c.sendRequest(req, response)
	return
}
