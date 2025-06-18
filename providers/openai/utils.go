/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-18 15:05:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-18 15:05:39
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"net/http"
)

// executeRequest 执行请求
func (s *openAIProvider) executeRequest(ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, response httpclient.Response, reqSetters ...httpclient.RequestOption) (err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(reqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, method, s.hClient.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = s.hClient.SendRequest(req, response)
	return
}

// executeStreamRequest 执行流式传输请求
func executeStreamRequest[T httpclient.Streamable](s *openAIProvider, ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, reqSetters ...httpclient.RequestOption) (stream *httpclient.StreamReader[T], err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(reqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, method, s.hClient.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送流式请求
	return httpclient.SendRequestStream[T](s.hClient, req)
}
