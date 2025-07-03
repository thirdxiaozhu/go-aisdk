/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-07-03 13:54:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 17:39:51
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

// Deserializer 反序列化器
type Deserializer struct {
	provider   string // 提供商
	streamable bool   // 是否流式
}

// NewDeserializer 创建反序列化器
func NewDeserializer(provider string, streamable bool) (s *Deserializer) {
	return &Deserializer{
		provider:   provider,
		streamable: streamable,
	}
}

// Unmarshal 反序列化
func (s *Deserializer) Unmarshal(data []byte, v any) (err error) {
	// 设置provider（如果目标对象实现了SetProvider接口）
	if setter, ok := v.(interface{ SetProvider(provider string) }); ok {
		setter.SetProvider(s.provider)
	}
	// 设置streamable（如果目标对象实现了SetStreamable接口）
	if setter, ok := v.(interface{ SetStreamable(streamable bool) }); ok {
		setter.SetStreamable(s.streamable)
	}
	// 反序列化
	return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}

// Decode 解码响应数据
func (s *Deserializer) Decode(body io.Reader, v any) (err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		// 读取全部数据
		var data []byte
		if data, err = io.ReadAll(body); err != nil {
			return
		}
		return s.Unmarshal(data, v)
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
