/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-16 16:16:29
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 16:16:37
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
)

// ListModels 列出模型
func (c *SDKClient) ListModels(ctx context.Context, request models.ListModelsRequest, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		// 列出模型
		return ps.ListModels(ctx, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.ModelInfo, "ListModels", request.UserInfo.UserID, nil, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ListModelsResponse)
	return
}
