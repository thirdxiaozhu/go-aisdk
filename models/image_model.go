/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:11:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 10:44:57
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"io"
)

// ImageBackground 背景透明度类型
type ImageBackground string

const (
	ImageBackgroundTransparent ImageBackground = "transparent" // 透明
	ImageBackgroundOpaque      ImageBackground = "opaque"      // 不透明
	ImageBackgroundAuto        ImageBackground = "auto"        // 自动
)

// ImageModeration 内容审核级别
type ImageModeration string

const (
	ImageModerationLow  ImageModeration = "low"  // 低
	ImageModerationAuto ImageModeration = "auto" // 自动
)

// ImageOutputFormat 返回图像格式
type ImageOutputFormat string

const (
	ImageOutputFormatPNG  ImageOutputFormat = "png"
	ImageOutputFormatJPEG ImageOutputFormat = "jpeg"
	ImageOutputFormatWebP ImageOutputFormat = "webp"
)

// ImageQuality 图像质量
type ImageQuality string

const (
	ImageQualityStandard ImageQuality = "standard" // 标准
	ImageQualityHD       ImageQuality = "hd"       // 高清
	ImageQualityHigh     ImageQuality = "high"     // 高
	ImageQualityMedium   ImageQuality = "medium"   // 中
	ImageQualityLow      ImageQuality = "low"      // 低
	ImageQualityAuto     ImageQuality = "auto"     // 自动
)

// ImageResponseFormat 响应格式
type ImageResponseFormat string

const (
	ImageResponseFormatURL     ImageResponseFormat = "url"      // URL
	ImageResponseFormatB64JSON ImageResponseFormat = "b64_json" // Base64
)

// ImageSize 图像尺寸
type ImageSize string

const (
	ImageSizeAuto      ImageSize = "auto"      // 自动
	ImageSize256x256   ImageSize = "256x256"   // 256x256
	ImageSize512x512   ImageSize = "512x512"   // 512x512
	ImageSize1024x1024 ImageSize = "1024x1024" // 1024x1024
	ImageSize1792x1024 ImageSize = "1792x1024" // 1792x1024
	ImageSize1024x1792 ImageSize = "1024x1792" // 1024x1792
	ImageSize1536x1024 ImageSize = "1536x1024" // 1536x1024
	ImageSize1024x1536 ImageSize = "1024x1536" // 1024x1536
)

// ImageStyle 图像风格
type ImageStyle string

const (
	ImageStyleVivid   ImageStyle = "vivid"   // 超现实
	ImageStyleNatural ImageStyle = "natural" // 自然
)

// ImageRequest 创建图像请求
type ImageRequest struct {
	UserInfo
	Provider          consts.Provider     `json:"provider"`                     // 提供商
	Prompt            string              `json:"prompt"`                       // 提示词
	Background        ImageBackground     `json:"background,omitempty"`         // 设置生成图像的背景透明度
	Model             string              `json:"model,omitempty"`              // 模型名称
	Moderation        ImageModeration     `json:"moderation,omitempty"`         // 内容审核级别
	N                 int                 `json:"n,omitempty"`                  // 生成的图像数量。必须在1到10之间
	OutputCompression int                 `json:"output_compression,omitempty"` // 图像压缩级别(0-100%)
	OutputFormat      ImageOutputFormat   `json:"output_format,omitempty"`      // 返回图像格式
	Quality           ImageQuality        `json:"quality,omitempty"`            // 图像质量
	ResponseFormat    ImageResponseFormat `json:"response_format,omitempty"`    // 响应格式
	Size              ImageSize           `json:"size,omitempty"`               // 图像尺寸
	Style             ImageStyle          `json:"style,omitempty"`              // 图像风格
}

// MarshalJSON 序列化JSON
func (r ImageRequest) MarshalJSON() (b []byte, err error) {
	// 创建一个别名结构体
	type Alias ImageRequest
	temp := struct {
		Provider string `json:"provider,omitempty"`
		UserID   string `json:"user_id,omitempty"`
		User     string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
		Alias
	}{
		Provider: "",
		User:     r.UserInfo.UserID,
		Alias:    Alias(r),
	}
	// 序列化JSON
	return json.Marshal(temp)
}

// ImageEditRequest 编辑图像请求
type ImageEditRequest struct {
	UserInfo
	Provider          consts.Provider     `json:"provider"`                     // 提供商
	Image             []io.Reader         `json:"image,omitempty"`              // 要编辑的图像源数组
	Prompt            string              `json:"prompt"`                       // 提示词
	Background        ImageBackground     `json:"background,omitempty"`         // 设置生成图像的背景透明度
	Mask              io.Reader           `json:"mask,omitempty"`               // mask图像源，其中完全透明的区域指示应该编辑的位置
	Model             string              `json:"model,omitempty"`              // 模型名称
	N                 int                 `json:"n,omitempty"`                  // 生成的图像数量。必须在1到10之间
	OutputCompression int                 `json:"output_compression,omitempty"` // 图像压缩级别(0-100%)
	OutputFormat      ImageOutputFormat   `json:"output_format,omitempty"`      // 返回图像格式
	Quality           ImageQuality        `json:"quality,omitempty"`            // 图像质量
	ResponseFormat    ImageResponseFormat `json:"response_format,omitempty"`    // 响应格式
	Size              ImageSize           `json:"size,omitempty"`               // 图像尺寸
}

// MarshalJSON 序列化JSON
func (r ImageEditRequest) MarshalJSON() (b []byte, err error) {
	// 创建一个别名结构体
	type Alias ImageEditRequest
	temp := struct {
		Provider string `json:"provider,omitempty"`
		UserID   string `json:"user_id,omitempty"`
		User     string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
		Alias
	}{
		Provider: "",
		User:     r.UserInfo.UserID,
		Alias:    Alias(r),
	}
	// 序列化JSON
	return json.Marshal(temp)
}

// ImageVariationRequest 变换图像请求
type ImageVariationRequest struct {
	UserInfo
	Provider       consts.Provider     `json:"provider"`                  // 提供商
	Image          io.Reader           `json:"image"`                     // 用作变换的基础图像。必须是有效的PNG文件，小于4MB，且为正方形
	Model          string              `json:"model,omitempty"`           // 模型名称
	N              int                 `json:"n,omitempty"`               // 生成的图像数量。必须在1到10之间
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty"` // 响应格式
	Size           ImageSize           `json:"size,omitempty"`            // 图像尺寸
}

// MarshalJSON 序列化JSON
func (r ImageVariationRequest) MarshalJSON() (b []byte, err error) {
	// 创建一个别名结构体
	type Alias ImageVariationRequest
	temp := struct {
		Provider string `json:"provider,omitempty"`
		UserID   string `json:"user_id,omitempty"`
		User     string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
		Alias
	}{
		Provider: "",
		User:     r.UserInfo.UserID,
		Alias:    Alias(r),
	}
	// 序列化JSON
	return json.Marshal(temp)
}

// ImageResponseData 生成图像的列表
type ImageResponseData struct {
	URL           string `json:"url,omitempty"`            // 生成图像的URL地址
	B64JSON       string `json:"b64_json,omitempty"`       // 生成图像的base64编码JSON
	RevisedPrompt string `json:"revised_prompt,omitempty"` // 用于生成图像的修订提示词
}

// ImageUsageInputTokensDetails 图像生成输入token的详细信息
type ImageUsageInputTokensDetails struct {
	ImageTokens int `json:"image_tokens,omitempty"` // 输入提示中的图像token数量
	TextTokens  int `json:"text_tokens,omitempty"`  // 输入提示中的文本token数量
}

// ImageUsage 图像生成的token使用信息
type ImageUsage struct {
	InputTokens        int                           `json:"input_tokens,omitempty"`         // 输入提示中的token总数（包括图像和文本）
	InputTokensDetails *ImageUsageInputTokensDetails `json:"input_tokens_details,omitempty"` // 图像生成输入token的详细信息
	OutputTokens       int                           `json:"output_tokens,omitempty"`        // 输出图像中的图像token数量
	TotalTokens        int                           `json:"total_tokens,omitempty"`         // 图像生成使用的token总数（图像和文本）
}

// ImageResponse 创建图像响应
type ImageResponse struct {
	Created int64               `json:"created"`         // 创建图像完成时的 Unix 时间戳（以秒为单位）
	Data    []ImageResponseData `json:"data"`            // 生成图像的列表
	Usage   *ImageUsage         `json:"usage,omitempty"` // 图像生成的token使用信息
	httpclient.HttpHeader
}
