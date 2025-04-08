/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-08 12:19:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 13:28:40
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai_test

import (
	"context"
	"github.com/liusuxian/go-openai"
	"os"
	"testing"
)

func TestBaiduAppBuilderIntegrated(t *testing.T) {
	var (
		ctx       = context.Background()
		authToken = os.Getenv("BAIDU_APPBUILDER_AUTH_TOKEN")
		response  *openai.BaiduAppBuilderIntegratedRes
		err       error
	)

	c := openai.NewClient("https://appbuilder.baidu.com", authToken)
	if response, err = c.BaiduAppBuilderIntegrated(ctx, openai.BaiduAppBuilderIntegratedReq{
		Query: "请帮我写一遍新中式装修的小红书营销文案",
	}); err != nil {
		t.Fatalf("BaiduAppBuilderIntegrated error: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("BaiduAppBuilderIntegrated error: %v", response.Message)
	}

	if response.Result == nil {
		t.Fatalf("BaiduAppBuilderIntegrated error: %v", "response.Result is nil")
	}

	if response.Result.Answer == "" {
		t.Fatalf("BaiduAppBuilderIntegrated error: %v", "response.Result.Answer is empty")
	}
}
