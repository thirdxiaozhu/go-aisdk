/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-07-09 15:34:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-09 17:10:41
 * @Description: 模型特性位掩码定义
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package consts

import "strings"

// ModelFeature 模型特性位掩码类型
type ModelFeature uint

// 无特性常量
const (
	ModelFeatureNone ModelFeature = 0 // 0b00000000 = 0
)

// 模型特性位掩码常量，每一位代表一个特性开关
const (
	// 第1位 (bit 0): 是否支持多模态 (文本+图像/音频/视频)
	ModelFeatureMultimodal ModelFeature = 1 << iota // 0b00000001 = 1
	// 第2位 (bit 1): 是否支持推理 (具备强推理能力)
	ModelFeatureReasoning ModelFeature = 1 << iota // 0b00000010 = 2
	// 第3位 (bit 2): 是否支持仅流式传输 (不支持阻塞式调用)
	ModelFeatureStreamingOnly ModelFeature = 1 << iota // 0b00000100 = 4
	// 更多特性位可以继续扩展...
)

// 预定义的常用特性组合
const (
	// 仅流式推理模型：推理 + 仅流式传输
	ModelFeatureReasoningStream = ModelFeatureReasoning | ModelFeatureStreamingOnly
	// 仅流式多模态模型：多模态 + 仅流式传输
	ModelFeatureMultimodalStream = ModelFeatureMultimodal | ModelFeatureStreamingOnly
	// 多模态推理模型：推理 + 多模态
	ModelFeatureAdvanced = ModelFeatureReasoning | ModelFeatureMultimodal
	// 仅流式多模态推理模型：推理 + 多模态 + 仅流式传输
	ModelFeatureAdvancedStream = ModelFeatureReasoning | ModelFeatureMultimodal | ModelFeatureStreamingOnly
)

// HasFeature 检查模型是否具有指定特性
func (f ModelFeature) HasFeature(feature ModelFeature) (ok bool) {
	return f&feature != 0
}

// IsMultimodal 检查是否支持多模态
func (f ModelFeature) IsMultimodal() (ok bool) {
	return f.HasFeature(ModelFeatureMultimodal)
}

// IsReasoningModel 检查是否支持推理
func (f ModelFeature) IsReasoningModel() (ok bool) {
	return f.HasFeature(ModelFeatureReasoning)
}

// IsStreamingOnly 检查是否支持仅流式传输
func (f ModelFeature) IsStreamingOnly() (ok bool) {
	return f.HasFeature(ModelFeatureStreamingOnly)
}

// String 返回特性的可读描述
func (f ModelFeature) String() (str string) {
	var features []string

	if f.IsMultimodal() {
		features = append(features, "multimodal")
	}
	if f.IsReasoningModel() {
		features = append(features, "reasoning")
	}
	if f.IsStreamingOnly() {
		features = append(features, "streaming-only")
	}

	if len(features) == 0 {
		return "none"
	}

	return strings.Join(features, "|")
}
