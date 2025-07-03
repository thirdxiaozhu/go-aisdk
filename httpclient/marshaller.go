/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:29:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 18:29:27
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient

import "encoding/json"

// Marshaller 序列化接口
type Marshaller interface {
	Marshal(val any) (b []byte, err error) // 序列化
}

// JSONMarshaller JSON 序列化
type JSONMarshaller struct{}

// Marshal 序列化
func (jm *JSONMarshaller) Marshal(val any) (b []byte, err error) {
	return json.Marshal(val)
}

// MustString 将数据转换为json字符串，如果转换失败，返回空字符串
func MustString(v any) (str string) {
	str, _ = String(v)
	return
}

// String 将数据转换为json字符串
func String(v any) (str string, err error) {
	if v == nil {
		return
	}

	var b []byte
	if b, err = json.Marshal(v); err != nil {
		return
	}

	str = string(b)
	return
}
