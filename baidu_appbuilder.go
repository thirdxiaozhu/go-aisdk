/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-08 12:19:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 13:38:42
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

// BaiduAppBuilderIntegratedStream 集成API，流式返回
func (c *Client) BaiduAppBuilderIntegratedStream(ctx context.Context, request BaiduAppBuilderIntegratedReq) (stream *BaiduAppBuilderIntegratedResStream, err error) {
	var (
		setters = []requestOption{
			setBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "streaming",
				"conversation_id": request.ConversationId,
			}),
			setKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.config.AuthToken)),
		}
		req *http.Request
	)

	if req, err = c.newRequest(ctx, http.MethodPost, c.fullURL(baiduAppBuilderIntegratedURL), setters...); err != nil {
		return
	}

	var resp *streamReader[BaiduAppBuilderIntegratedResult]
	if resp, err = sendRequestStream[BaiduAppBuilderIntegratedResult](c, req); err != nil {
		return
	}

	stream = &BaiduAppBuilderIntegratedResStream{
		streamReader: resp,
	}
	return
}
