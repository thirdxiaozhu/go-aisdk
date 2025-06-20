/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:37:53
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-20 22:52:42
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
	"strconv"
)

const (
	apiImagesGenerations = "/images/generations"
	apiImagesEdits       = "/images/edits"
)

// CreateImage 创建图像
func (s *openAIProvider) CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, apiImagesGenerations, opts, s.lb, nil, &response, httpclient.WithBody(request))
	return
}

// CreateImageEdit 编辑图像
func (s *openAIProvider) CreateImageEdit(ctx context.Context, request models.ImageEditRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	formHandler := func(builder httpclient.FormBuilder) (e error) {
		// 要编辑的图像源数组
		for _, imageReader := range request.Image {
			if e = builder.CreateFormFileReader("image", imageReader, ""); e != nil {
				return
			}
		}
		// 提示词
		if e = builder.WriteField("prompt", request.Prompt); e != nil {
			return
		}
		// 设置生成图像的背景透明度
		if request.Background != "" {
			if e = builder.WriteField("background", string(request.Background)); e != nil {
				return
			}
		}
		// mask图像源，其中完全透明的区域指示应该编辑的位置
		if request.Mask != nil {
			if e = builder.CreateFormFileReader("mask", request.Mask, ""); e != nil {
				return
			}
		}
		// 模型名称
		if request.Model != "" {
			if e = builder.WriteField("model", request.Model); e != nil {
				return
			}
		}
		// 生成图像数量
		if request.N > 0 {
			if e = builder.WriteField("n", strconv.Itoa(request.N)); e != nil {
				return
			}
		}
		// 图像压缩级别(0-100%)
		if request.OutputCompression > 0 {
			if e = builder.WriteField("output_compression", strconv.Itoa(request.OutputCompression)); e != nil {
				return
			}
		}
		// 返回图像格式
		if request.OutputFormat != "" {
			if e = builder.WriteField("output_format", string(request.OutputFormat)); e != nil {
				return
			}
		}
		// 图像质量
		if request.Quality != "" {
			if e = builder.WriteField("quality", string(request.Quality)); e != nil {
				return
			}
		}
		// 响应格式
		if request.ResponseFormat != "" {
			if e = builder.WriteField("response_format", string(request.ResponseFormat)); e != nil {
				return
			}
		}
		// 图像尺寸
		if request.Size != "" {
			if e = builder.WriteField("size", string(request.Size)); e != nil {
				return
			}
		}
		// 用户
		if request.User != "" {
			if e = builder.WriteField("user", request.User); e != nil {
				return
			}
		}
		// 关闭构建器
		return builder.Close()
	}
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, apiImagesEdits, opts, s.lb, formHandler, &response)
	return
}
