/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-20 01:15:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-20 01:25:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"net/http"
)

// ExecuteRequest 执行请求
func ExecuteRequest(
	ctx context.Context,
	method, baseURL, apiPath string,
	opts []httpclient.HTTPClientOption,
	lb *loadbalancer.LoadBalancer,
	formHandler httpclient.FormBuilderHandler,
	response httpclient.Response,
	reqSetters ...httpclient.RequestOption,
) (err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClient(baseURL)
	// 设置客户端选项
	for _, opt := range opts {
		opt(hc)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = lb.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(reqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	// 构建表单请求体
	if formHandler != nil {
		var (
			formBody = &bytes.Buffer{}
			builder  = hc.GetFormBuilder(formBody)
		)
		if err = formHandler(builder); err != nil {
			return
		}
		setters = append(setters, httpclient.WithBody(formBody), httpclient.WithContentType(builder.FormDataContentType()))
	}
	if req, err = hc.NewRequest(ctx, method, hc.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = hc.SendRequest(req, response)
	return
}

// ExecuteStreamRequest 执行流式传输请求
func ExecuteStreamRequest[T httpclient.Streamable](
	ctx context.Context,
	method, baseURL, apiPath string,
	opts []httpclient.HTTPClientOption,
	lb *loadbalancer.LoadBalancer,
	reqSetters ...httpclient.RequestOption,
) (stream *httpclient.StreamReader[T], err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClient(baseURL)
	// 设置客户端选项
	for _, opt := range opts {
		opt(hc)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = lb.GetAPIKey(); err != nil {
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
