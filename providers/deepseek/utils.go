/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-18 15:00:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 14:14:31
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"net/http"
)

// executeRequest 执行请求
func (s *deepseekProvider) executeRequest(ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, response httpclient.Response, reqSetters ...httpclient.RequestOption) (err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	// 设置客户端选项
	for _, opt := range opts {
		opt(hc)
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
	if req, err = hc.NewRequest(ctx, method, hc.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = hc.SendRequest(req, response)
	return
}

// executeStreamRequest 执行流式传输请求
func executeStreamRequest[T httpclient.Streamable](s *deepseekProvider, ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, reqSetters ...httpclient.RequestOption) (stream *httpclient.StreamReader[T], err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	// 设置客户端选项
	for _, opt := range opts {
		opt(hc)
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
	if req, err = hc.NewRequest(ctx, method, hc.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送流式请求
	return httpclient.SendRequestStream[T](hc, req)
}
