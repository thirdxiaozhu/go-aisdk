/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 17:30:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-30 17:33:28
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package middleware

import "encoding/json"

// jsonString 将数据转换为json字符串
func jsonString(v any) (str string, err error) {
	if v == nil {
		return
	}
	var b []byte
	b, err = json.Marshal(v)
	str = string(b)
	return
}

// mustJsonString 将数据转换为json字符串，如果转换失败，返回空字符串
func mustJsonString(v any) (str string) {
	str, _ = jsonString(v)
	return
}
