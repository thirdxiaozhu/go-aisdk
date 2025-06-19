/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:37:53
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 19:12:14
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"net/http"
)

const (
	apiImagesGenerations = "/images/generations"
)

// CreateImage 创建图像
func (s *openAIProvider) CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = s.executeRequest(ctx, http.MethodPost, apiImagesGenerations, opts, &response, httpclient.WithBody(request))
	return
}
