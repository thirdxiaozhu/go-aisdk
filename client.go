/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:22:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 12:18:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/liusuxian/go-openai/internal"
	"io"
	"net/http"
	"strings"
)

// Client 客户端
type Client struct {
	config            ClientConfig                           // 客户端配置
	requestBuilder    utils.RequestBuilder                   // 请求构建器
	createFormBuilder func(body io.Writer) utils.FormBuilder // 表单构建器
}

// Response 响应
type Response interface {
	SetHeader(http.Header)
}

type httpHeader http.Header

// SetHeader 设置请求头
func (h *httpHeader) SetHeader(header http.Header) {
	*h = httpHeader(header)
}

// Header 获取请求头
func (h *httpHeader) Header() (header http.Header) {
	return http.Header(*h)
}

// GetRateLimitHeaders 获取速率限制请求头
func (h *httpHeader) GetRateLimitHeaders() (rateLimit RateLimitHeaders) {
	return newRateLimitHeaders(h.Header())
}

// RawResponse 原始响应
type RawResponse struct {
	io.ReadCloser
	httpHeader
}

// NewClient 新建客户端
func NewClient(baseURL, authToken string) (c *Client) {
	return NewClientWithConfig(DefaultConfig(baseURL, authToken))
}

// NewClientWithConfig 通过客户端配置新建客户端
func NewClientWithConfig(config ClientConfig) (c *Client) {
	return &Client{
		config:         config,
		requestBuilder: utils.NewRequestBuilder(),
		createFormBuilder: func(body io.Writer) utils.FormBuilder {
			return utils.NewFormBuilder(body)
		},
	}
}

// requestOptions 请求选项
type requestOptions struct {
	body   any
	header http.Header
}

// requestOption 请求选项配置器
type requestOption func(reqOpts *requestOptions)

// setBody 设置 HTTP 请求的主体内容
func setBody(body any) (reqOpt requestOption) {
	return func(reqOpts *requestOptions) {
		reqOpts.body = body
	}
}

// setContentType 设置 HTTP 请求头的 Content-Type 字段
func setContentType(contentType string) (reqOpt requestOption) {
	return func(reqOpts *requestOptions) {
		reqOpts.header.Set("Content-Type", contentType)
	}
}

// setCookie 设置 HTTP 请求头的 Cookie 字段
func setCookie(cookies []*http.Cookie) (reqOpt requestOption) {
	return func(reqOpts *requestOptions) {
		cookieList := make([]string, 0, len(cookies))
		for _, v := range cookies {
			cookieList = append(cookieList, fmt.Sprintf("%s=%s", v.Name, v.Value))
		}
		reqOpts.header.Set("Cookie", strings.Join(cookieList, "; "))
	}
}

// setKeyValue 设置 HTTP 请求头的键值对
func setKeyValue(key, value string) (reqOpt requestOption) {
	return func(reqOpts *requestOptions) {
		reqOpts.header.Set(key, value)
	}
}

// addKeyValue 添加 HTTP 请求头的键值对
func addKeyValue(key, value string) (reqOpt requestOption) {
	return func(reqOpts *requestOptions) {
		reqOpts.header.Add(key, value)
	}
}

// newRequest 新建请求
func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (req *http.Request, err error) {
	reqOpts := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(reqOpts)
	}

	if req, err = c.requestBuilder.Build(ctx, method, url, reqOpts.body, reqOpts.header); err != nil {
		return
	}
	return
}

// sendRequest 发送请求
func (c *Client) sendRequest(req *http.Request, v Response) (err error) {
	// 设置默认请求头
	for _, v := range [][]string{
		{"Content-Type", "application/json"},
		{"Accept", "application/json"},
	} {
		if req.Header.Get(v[0]) == "" {
			req.Header.Set(v[0], v[1])
		}
	}

	var resp *http.Response
	if resp, err = c.config.HTTPClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if v != nil {
		v.SetHeader(resp.Header)
	}

	if isFailureStatusCode(resp) {
		return c.handleErrorResp(resp)
	}

	return decodeResponse(resp.Body, v)
}

// sendRequestRaw 发送请求
func (c *Client) sendRequestRaw(req *http.Request) (response RawResponse, err error) {
	var resp *http.Response
	if resp, err = c.config.HTTPClient.Do(req); err != nil {
		return
	}

	if isFailureStatusCode(resp) {
		err = c.handleErrorResp(resp)
		return
	}

	response.SetHeader(resp.Header)
	response.ReadCloser = resp.Body
	return
}

// sendRequestStream 发送请求
func sendRequestStream[T streamable](client *Client, req *http.Request) (stream *streamReader[T], err error) {
	// 设置默认请求头
	for _, v := range [][]string{
		{"Content-Type", "application/json"},
		{"Accept", "text/event-stream"},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
	} {
		if req.Header.Get(v[0]) == "" {
			req.Header.Set(v[0], v[1])
		}
	}

	var resp *http.Response
	if resp, err = client.config.HTTPClient.Do(req); err != nil {
		stream = &streamReader[T]{}
		return
	}

	if isFailureStatusCode(resp) {
		stream = &streamReader[T]{}
		err = client.handleErrorResp(resp)
		return
	}

	stream = &streamReader[T]{
		emptyMessagesLimit: client.config.EmptyMessagesLimit,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		errAccumulator:     utils.NewErrorAccumulator(),
		unmarshaler:        &utils.JSONUnmarshaler{},
	}
	return
}

// isFailureStatusCode 是否失败状态码
func isFailureStatusCode(resp *http.Response) (ok bool) {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

// decodeResponse 解码响应数据
func decodeResponse(body io.Reader, v any) (err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

// decodeString 解码字符串
func decodeString(body io.Reader, output *string) (err error) {
	var b []byte
	if b, err = io.ReadAll(body); err != nil {
		return
	}
	*output = string(b)
	return
}

// fullURL 获取完整链接
func (c *Client) fullURL(suffix string) (url string) {
	baseURL := strings.TrimRight(c.config.BaseURL, "/")
	return fmt.Sprintf("%s%s", baseURL, suffix)
}

// handleErrorResp 处理错误响应
func (c *Client) handleErrorResp(resp *http.Response) (err error) {
	// 读取响应体
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("error, reading response body: %w", err)
	}
	// 尝试将响应体解析为 JSON
	var errResp map[string]any
	if err = json.Unmarshal(body, &errResp); err != nil {
		// 如果解析失败，返回包含原始响应体的错误
		return &RequestError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
			Body:           body,
		}
	}
	// 处理 errResp 为空的情况
	if len(errResp) == 0 {
		return &RequestError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			Err:            fmt.Errorf("empty error response"),
			Body:           body,
		}
	}
	// 成功解析 JSON 后，返回包含错误信息的 RequestError
	return &RequestError{
		HTTPStatus:     resp.Status,
		HTTPStatusCode: resp.StatusCode,
		Err:            fmt.Errorf("%v", errResp),
		Body:           body,
	}
}
