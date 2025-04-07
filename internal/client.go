/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:22:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-07 18:25:11
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"encoding/json"
	"io"
	"net/http"
)

// IsFailureStatusCode 是否失败状态码
func IsFailureStatusCode(resp *http.Response) (ok bool) {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

// DecodeString 解码字符串
func DecodeString(body io.Reader, output *string) (err error) {
	var b []byte
	if b, err = io.ReadAll(body); err != nil {
		return
	}
	*output = string(b)
	return
}

// DecodeResponse 解码响应数据
func DecodeResponse(body io.Reader, v any) (err error) {
	if v == nil {
		return
	}
	if result, ok := v.(*string); ok {
		return DecodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}
