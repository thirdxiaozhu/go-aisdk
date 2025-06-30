/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:11:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 00:20:34
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"io"
)

var (
	// 序列化创建图像请求函数（OpenAI）
	marshalImageRequestByOpenAI = func(r ImageRequest) (b []byte, err error) {
		type Alias ImageRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			User   string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
			Alias
		}{
			User:  r.UserInfo.User,
			Alias: Alias(r),
		}
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}

	// 序列化创建图像请求函数（Ark）
	marshalImageRequestByArk = func(r ImageRequest) (b []byte, err error) {
		type Alias ImageRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			User   string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
			Alias
		}{
			User:  r.UserInfo.User,
			Alias: Alias(r),
		}
		temp.Provider = ""
		temp.Background = ""
		temp.Moderation = ""
		temp.N = 0

		// 序列化JSON
		return json.Marshal(temp)
	}
	// 策略映射（创建图像请求）
	imageRequestStrategies = map[consts.Provider]func(r ImageRequest) (b []byte, err error){
		consts.OpenAI: marshalImageRequestByOpenAI,
		consts.Ark:    marshalImageRequestByArk,
	}
	// 序列化编辑图像请求函数（OpenAI）
	marshalImageEditRequestByOpenAI = func(r ImageEditRequest) (b []byte, err error) {
		type Alias ImageEditRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			User   string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
			Alias
		}{
			User:  r.UserInfo.User,
			Alias: Alias(r),
		}
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}
	// 策略映射（编辑图像请求）
	imageEditRequestStrategies = map[consts.Provider]func(r ImageEditRequest) (b []byte, err error){
		consts.OpenAI: marshalImageEditRequestByOpenAI,
	}
	// 序列化变换图像请求函数（OpenAI）
	marshalImageVariationRequestByOpenAI = func(r ImageVariationRequest) (b []byte, err error) {
		type Alias ImageVariationRequest
		temp := struct {
			UserID string `json:"user_id,omitempty"`
			User   string `json:"user,omitempty"` // 用户标识符，用于监控和滥用检测
			Alias
		}{
			User:  r.UserInfo.User,
			Alias: Alias(r),
		}
		// 序列化JSON
		temp.Provider = ""
		return json.Marshal(temp)
	}
	// 策略映射（变换图像请求）
	imageVariationRequestStrategies = map[consts.Provider]func(r ImageVariationRequest) (b []byte, err error){
		consts.OpenAI: marshalImageVariationRequestByOpenAI,
	}
)

// ImageBackground 背景透明度类型
//
//	提供商支持: OpenAI
type ImageBackground string

const (
	// 透明
	// 提供商支持: OpenAI
	ImageBackgroundTransparent ImageBackground = "transparent"
	// 不透明
	// 提供商支持: OpenAI
	ImageBackgroundOpaque ImageBackground = "opaque"
	// 自动
	// 提供商支持: OpenAI
	ImageBackgroundAuto ImageBackground = "auto"
)

// ImageModeration 内容审核级别
//
//	提供商支持: OpenAI
type ImageModeration string

const (
	// 低
	// 提供商支持: OpenAI
	ImageModerationLow ImageModeration = "low"
	// 自动
	// 提供商支持: OpenAI
	ImageModerationAuto ImageModeration = "auto"
)

// ImageOutputFormat 返回图像格式
//
//	提供商支持: OpenAI
type ImageOutputFormat string

const (
	// PNG
	// 提供商支持: OpenAI
	ImageOutputFormatPNG ImageOutputFormat = "png"
	// JPEG
	// 提供商支持: OpenAI
	ImageOutputFormatJPEG ImageOutputFormat = "jpeg"
	// WebP
	// 提供商支持: OpenAI
	ImageOutputFormatWebP ImageOutputFormat = "webp"
)

// ImageQuality 图像质量
//
//	提供商支持: OpenAI
type ImageQuality string

const (
	// 标准
	// 提供商支持: OpenAI
	ImageQualityStandard ImageQuality = "standard"
	// 高清
	// 提供商支持: OpenAI
	ImageQualityHD ImageQuality = "hd"
	// 高
	// 提供商支持: OpenAI
	ImageQualityHigh ImageQuality = "high"
	// 中
	// 提供商支持: OpenAI
	ImageQualityMedium ImageQuality = "medium"
	// 低
	// 提供商支持: OpenAI
	ImageQualityLow ImageQuality = "low"
	// 自动
	// 提供商支持: OpenAI
	ImageQualityAuto ImageQuality = "auto"
)

// ImageResponseFormat 响应格式
//
//	提供商支持: OpenAI
type ImageResponseFormat string

const (
	// URL
	// 提供商支持: OpenAI
	ImageResponseFormatURL ImageResponseFormat = "url"
	// Base64
	// 提供商支持: OpenAI
	ImageResponseFormatB64JSON ImageResponseFormat = "b64_json"
)

// ImageSize 图像尺寸
//
//	提供商支持: OpenAI
type ImageSize string

const (
	// 自动
	// 提供商支持: OpenAI
	ImageSizeAuto ImageSize = "auto"
	// 256x256
	// 提供商支持: OpenAI
	ImageSize256x256 ImageSize = "256x256"
	// 512x512
	// 提供商支持: OpenAI
	ImageSize512x512 ImageSize = "512x512"
	// 1024x1024
	// 提供商支持: OpenAI | Ark
	ImageSize1024x1024 ImageSize = "1024x1024"
	// 1792x1024
	// 提供商支持: OpenAI
	ImageSize1792x1024 ImageSize = "1792x1024"
	// 1024x1792
	// 提供商支持: OpenAI
	ImageSize1024x1792 ImageSize = "1024x1792"
	// 1536x1024
	// 提供商支持: OpenAI
	ImageSize1536x1024 ImageSize = "1536x1024"
	// 1024x1536
	// 提供商支持: OpenAI
	ImageSize1024x1536 ImageSize = "1024x1536"
	// 864x1152
	// 提供商支持: Ark
	ImageSize864x1152 ImageSize = "864x1152"
	// 1152x864
	// 提供商支持: Ark
	ImageSize1152x864 ImageSize = "1152x864"
	// 1280x720
	// 提供商支持: Ark
	ImageSize1280x720 ImageSize = "1280x720"
	// 720x1280
	// 提供商支持: Ark
	ImageSize720x1280 ImageSize = "720x1280"
	// 832x1248
	// 提供商支持: Ark
	ImageSize832x1248 ImageSize = "832x1248"
	// 1248x832
	// 提供商支持: Ark
	ImageSize1248x832 ImageSize = "1248x832"
	// 1512x648
	// 提供商支持: Ark
	ImageSize1512x648 ImageSize = "1512x648"
)

// ImageStyle 图像风格
//
//	提供商支持: OpenAI
type ImageStyle string

const (
	// 超现实
	// 提供商支持: OpenAI
	ImageStyleVivid ImageStyle = "vivid"
	// 自然
	// 提供商支持: OpenAI
	ImageStyleNatural ImageStyle = "natural"
)

// ImageRequest 创建图像请求
//
//	提供商支持: OpenAI
type ImageRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 提示词
	// 提供商支持: OpenAI | Ark
	Prompt string `json:"prompt,omitempty"`
	// 设置生成图像的背景透明度
	// 提供商支持: OpenAI
	Background ImageBackground `json:"background,omitempty"`
	// 模型名称
	// 提供商支持: OpenAI | Ark
	Model string `json:"model,omitempty"`
	// 内容审核级别
	// 提供商支持: OpenAI
	Moderation ImageModeration `json:"moderation,omitempty"`
	// 生成的图像数量
	// 提供商支持: OpenAI
	N int `json:"n,omitempty"`
	// 图像压缩级别(0-100%)
	// 提供商支持: OpenAI
	OutputCompression int `json:"output_compression,omitempty"`
	// 返回图像格式
	// 提供商支持: OpenAI
	OutputFormat ImageOutputFormat `json:"output_format,omitempty"`
	// 图像质量
	// 提供商支持: OpenAI
	Quality ImageQuality `json:"quality,omitempty"`
	// 响应格式
	// 提供商支持: OpenAI | Ark
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty"`
	// 图像尺寸
	// 提供商支持: OpenAI | Ark
	Size ImageSize `json:"size,omitempty"`
	// 图像风格
	// 提供商支持: OpenAI
	Style ImageStyle `json:"style,omitempty"`
	// 随机数种子 用于控制模型生成内容的随机性
	// 提供商支持: Ark
	Seed int `json:"seed,omitempty"`
	// 模型输出结果与prompt的一致程度
	// 提供商支持: Ark
	GuidanceScale float64 `json:"guidance_scale,omitempty"`
	// 是否添加水印
	// 提供商支持: Ark
	Watermark bool `json:"watermark,omitempty"`
}

// MarshalJSON 序列化JSON
func (r ImageRequest) MarshalJSON() (b []byte, err error) {
	strategy, ok := imageRequestStrategies[r.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", r.Provider)
	}

	return strategy(r)
}

// ImageEditRequest 编辑图像请求
//
//	提供商支持: OpenAI
type ImageEditRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 要编辑的图像源数组
	// 提供商支持: OpenAI
	Image []io.Reader `json:"image,omitempty"`
	// 提示词
	// 提供商支持: OpenAI
	Prompt string `json:"prompt,omitempty"`
	// 设置生成图像的背景透明度
	// 提供商支持: OpenAI
	Background ImageBackground `json:"background,omitempty"`
	// mask图像源，其中完全透明的区域指示应该编辑的位置
	// 提供商支持: OpenAI
	Mask io.Reader `json:"mask,omitempty"`
	// 模型名称
	// 提供商支持: OpenAI
	Model string `json:"model,omitempty"`
	// 生成的图像数量。必须在1到10之间
	// 提供商支持: OpenAI
	N int `json:"n,omitempty"`
	// 图像压缩级别(0-100%)
	// 提供商支持: OpenAI
	OutputCompression int `json:"output_compression,omitempty"`
	// 返回图像格式
	// 提供商支持: OpenAI
	OutputFormat ImageOutputFormat `json:"output_format,omitempty"`
	// 图像质量
	// 提供商支持: OpenAI
	Quality ImageQuality `json:"quality,omitempty"`
	// 响应格式
	// 提供商支持: OpenAI
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty"`
	// 图像尺寸
	// 提供商支持: OpenAI
	Size ImageSize `json:"size,omitempty"`
}

// MarshalJSON 序列化JSON
func (r ImageEditRequest) MarshalJSON() (b []byte, err error) {
	strategy, ok := imageEditRequestStrategies[r.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", r.Provider)
	}
	return strategy(r)
}

// ImageVariationRequest 变换图像请求
//
//	提供商支持: OpenAI
type ImageVariationRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 用作变换的基础图像。必须是有效的PNG文件，小于4MB，且为正方形
	// 提供商支持: OpenAI
	Image io.Reader `json:"image,omitempty"`
	// 模型名称
	// 提供商支持: OpenAI
	Model string `json:"model,omitempty"`
	// 生成的图像数量
	// 提供商支持: OpenAI
	N int `json:"n,omitempty"`
	// 响应格式
	// 提供商支持: OpenAI
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty"`
	// 图像尺寸
	// 提供商支持: OpenAI
	Size ImageSize `json:"size,omitempty"`
}

// MarshalJSON 序列化JSON
func (r ImageVariationRequest) MarshalJSON() (b []byte, err error) {
	strategy, ok := imageVariationRequestStrategies[r.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", r.Provider)
	}
	return strategy(r)
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
	Created int64               `json:"created,omitempty"` // 创建图像完成时的 Unix 时间戳（以秒为单位）
	Data    []ImageResponseData `json:"data,omitempty"`    // 生成图像的列表
	Usage   *ImageUsage         `json:"usage,omitempty"`   // 图像生成的token使用信息
	httpclient.HttpHeader
}
