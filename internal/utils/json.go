/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 17:30:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 17:33:28
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import "encoding/json"

// String 将数据转换为json字符串
func String(v any) (str string, err error) {
	if v == nil {
		return
	}
	var b []byte
	b, err = json.Marshal(v)
	str = string(b)
	return
}

// MustString 将数据转换为json字符串，如果转换失败，返回空字符串
func MustString(v any) (str string) {
	str, _ = String(v)
	return
}
