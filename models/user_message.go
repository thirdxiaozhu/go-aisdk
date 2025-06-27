/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 10:57:42
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
)

var (
	// 序列化用户消息函数（OpenAI）
	marshalUserMessageByOpenAI = func(m UserMessage) (b []byte, err error) {
		type Alias UserMessage
		temp := struct {
			Role    string `json:"role"`
			Content any    `json:"content"`
			Alias
		}{
			Role:  "user",
			Alias: Alias(m),
		}
		// 根据内容类型设置 content 字段
		if len(m.MultimodalContent) > 0 {
			// 移除不支持的字段
			tempMultimodalContent := make([]ChatUserMsgPart, 0, len(m.MultimodalContent))
			for _, v := range m.MultimodalContent {
				if len(v.Video) > 0 || v.VideoURL != nil {
					continue
				}
				if v.ImageURL != nil {
					v.MinPixels = 0
					v.MaxPixels = 0
				}
				tempMultimodalContent = append(tempMultimodalContent, v)
			}
			temp.Content = tempMultimodalContent
		} else {
			temp.Content = m.Content
		}
		// 移除 multimodal_content 字段
		temp.MultimodalContent = nil
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化用户消息函数（DeepSeek）
	marshalUserMessageByDeepSeek = func(m UserMessage) (b []byte, err error) {
		type Alias UserMessage
		temp := struct {
			Role string `json:"role"`
			Alias
		}{
			Role:  "user",
			Alias: Alias(m),
		}
		// 移除不支持的字段
		temp.MultimodalContent = nil
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化用户消息函数（AliBL）
	marshalUserMessageByAliBL = func(m UserMessage) (b []byte, err error) {
		type Alias UserMessage
		temp := struct {
			Role    string `json:"role"`
			Content any    `json:"content"`
			Alias
		}{
			Role:  "user",
			Alias: Alias(m),
		}
		// 根据内容类型设置 content 字段
		if len(m.MultimodalContent) > 0 {
			// 移除不支持的字段
			tempMultimodalContent := make([]ChatUserMsgPart, 0, len(m.MultimodalContent))
			for _, v := range m.MultimodalContent {
				if v.File != nil {
					continue
				}
				if v.ImageURL != nil {
					v.ImageURL.Detail = ""
				}
				tempMultimodalContent = append(tempMultimodalContent, v)
			}
			temp.Content = tempMultimodalContent
		} else {
			temp.Content = m.Content
		}
		// 移除 multimodal_content 字段
		temp.MultimodalContent = nil
		// 移除不支持的字段
		temp.Name = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 序列化用户消息函数（Ark）
	marshalUserMessageByArk = func(m UserMessage) (b []byte, err error) {
		type Alias UserMessage
		temp := struct {
			Role    string `json:"role"`
			Content any    `json:"content"`
			Alias
		}{
			Role:  "user",
			Alias: Alias(m),
		}
		// 根据内容类型设置 content 字段
		if len(m.MultimodalContent) > 0 {
			// 移除不支持的字段
			tempMultimodalContent := make([]ChatUserMsgPart, 0, len(m.MultimodalContent))
			for _, v := range m.MultimodalContent {
				if v.InputAudio != nil {
					v.InputAudio = nil
				}
				if v.File != nil {
					continue
				}
				if len(v.Video) != 0 {
					v.Video = nil
				}
				v.MinPixels = 0
				v.MaxPixels = 0
				tempMultimodalContent = append(tempMultimodalContent, v)
			}
			temp.Content = tempMultimodalContent
		} else {
			temp.Content = m.Content
		}
		// 移除 multimodal_content 字段
		temp.MultimodalContent = nil
		// 移除不支持的字段
		temp.Name = ""
		// 序列化JSON
		temp.provider = ""
		return json.Marshal(temp)
	}
	// 策略映射
	userMessageStrategies = map[consts.Provider]func(m UserMessage) (b []byte, err error){
		consts.OpenAI:   marshalUserMessageByOpenAI,
		consts.DeepSeek: marshalUserMessageByDeepSeek,
		consts.AliBL:    marshalUserMessageByAliBL,
		consts.Ark:      marshalUserMessageByArk,
	}
)

// ChatUserMsgPartType 多模态内容类型
//
//	提供商支持: OpenAI | AliBL
type ChatUserMsgPartType string

const (
	//	提供商支持: OpenAI | AliBL
	ChatUserMsgPartTypeText ChatUserMsgPartType = "text"
	//	提供商支持: OpenAI | AliBL
	ChatUserMsgPartTypeImageURL ChatUserMsgPartType = "image_url"
	//	提供商支持: OpenAI | AliBL
	ChatUserMsgPartTypeInputAudio ChatUserMsgPartType = "input_audio"
	//	提供商支持: OpenAI
	ChatUserMsgPartTypeFile ChatUserMsgPartType = "file"
	//	提供商支持: AliBL
	ChatUserMsgPartTypeVideo ChatUserMsgPartType = "video"
	//	提供商支持: AliBL
	ChatUserMsgPartTypeVideoURL ChatUserMsgPartType = "video_url"
)

// ChatUserMsgImageURLDetail 图像质量
//
//	提供商支持: OpenAI
type ChatUserMsgImageURLDetail string

const (
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailHigh ChatUserMsgImageURLDetail = "high"
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailLow ChatUserMsgImageURLDetail = "low"
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailAuto ChatUserMsgImageURLDetail = "auto"
)

// ChatUserMsgImageURLPixelLimit 允许设置图片的像素大小设置
//
// 提供商支持: Ark
type ChatUserMsgImageURLPixelLimit struct {
	// 传入图片最大像素值限制
	// 提供商支持: Ark
	MaxPixels int `json:"max_pixels,omitempty"`
	// 传入图片最小像素值限制
	// 提供商支持: Ark
	MinPixels int `json:"min_pixels,omitempty"`
}

// ChatUserMsgImageURL 图像URL
//
//	提供商支持: OpenAI | AliBL
type ChatUserMsgImageURL struct {
	// 图像URL，支持url和base64编码
	// 提供商支持: OpenAI | AliBL | Ark
	URL string `json:"url,omitempty"`
	// 图像质量
	// 提供商支持: OpenAI | Ark
	Detail ChatUserMsgImageURLDetail `json:"detail,omitempty"`
	// 像素值限制
	// 提供商支持: Ark
	ImagePixelLimit *ChatUserMsgImageURL `json:"image_pixel_limit,omitempty"`
}

// ChatUserMsgInputAudioFormat 音频格式
//
//	提供商支持: OpenAI | AliBL
type ChatUserMsgInputAudioFormat string

const (
	//	提供商支持: OpenAI | AliBL
	ChatUserMsgInputAudioFormatMP3 ChatUserMsgInputAudioFormat = "mp3"
	//	提供商支持: OpenAI | AliBL
	ChatUserMsgInputAudioFormatWAV ChatUserMsgInputAudioFormat = "wav"
)

// ChatUserMsgInputAudio 输入音频
//
//	提供商支持: OpenAI | AliBL
type ChatUserMsgInputAudio struct {
	// 音频数据，支持url和base64编码
	// 提供商支持: OpenAI | AliBL
	Data string `json:"data,omitempty"`
	// 音频格式
	// 提供商支持: OpenAI | AliBL
	Format ChatUserMsgInputAudioFormat `json:"format,omitempty"`
}

// ChatUserMsgFile 文件
//
//	提供商支持: OpenAI
type ChatUserMsgFile struct {
	// 文件数据，支持base64编码
	// 提供商支持: OpenAI
	FileData string `json:"file_data,omitempty"`
	// 文件ID
	// 提供商支持: OpenAI
	FileID string `json:"file_id,omitempty"`
	// 文件名
	// 提供商支持: OpenAI
	FileName string `json:"filename,omitempty"`
}

// ChatUserMsgVideoURL 视频URL
//
//	提供商支持: AliBL | Ark
type ChatUserMsgVideoURL struct {
	// 视频URL，支持url和base64编码
	// 提供商支持: AliBL | Ark
	URL string `json:"url,omitempty"`
	// 每秒钟从视频中抽取指定数量的图像
	// 提供商支持: Ark
	Fps float64 `json:"fps,omitempty"`
}

// ChatUserMsgPart 多模态内容
//
//	提供商支持: OpenAI | AliBL
type ChatUserMsgPart struct {
	// 内容类型
	// 提供商支持: OpenAI | AliBL | Ark
	Type ChatUserMsgPartType `json:"type,omitempty"`
	// 文本内容
	// 提供商支持: OpenAI | AliBL | Ark
	Text string `json:"text,omitempty"`
	// 图像URL
	// 提供商支持: OpenAI | AliBL | Ark
	ImageURL *ChatUserMsgImageURL `json:"image_url,omitempty"`
	// 输入音频
	// 提供商支持: OpenAI | AliBL
	InputAudio *ChatUserMsgInputAudio `json:"input_audio,omitempty"`
	// 文件
	// 提供商支持: OpenAI
	File *ChatUserMsgFile `json:"file,omitempty"`
	// 视频列表
	// 提供商支持: AliBL
	Video []string `json:"video,omitempty"`
	// 视频URL
	// 提供商支持: AliBL | Ark
	VideoURL *ChatUserMsgVideoURL `json:"video_url,omitempty"`
	// 模型限制输入图像的最小像素
	// 提供商支持: AliBL
	MinPixels int `json:"min_pixels,omitempty"`
	// 模型限制输入图像的最大像素
	// 提供商支持: AliBL
	MaxPixels int `json:"max_pixels,omitempty"`
}

// UserMessage 用户消息，支持多模态内容
//
//	提供商支持: OpenAI | DeepSeek | AliBL
type UserMessage struct {
	provider consts.Provider `json:"-"` // 用于序列化参数时，处理差异化参数
	// 文本内容
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Content string `json:"content,omitempty"`
	// 多模态内容
	// 提供商支持: OpenAI | AliBL | Ark
	MultimodalContent []ChatUserMsgPart `json:"-"`
	// 参与者名称
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty"`
}

// GetRole 获取消息角色
func (m UserMessage) GetRole() (role string) { return "user" }

// SetProvider 设置提供商
func (m *UserMessage) SetProvider(provider consts.Provider) {
	fmt.Println(m)
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m UserMessage) MarshalJSON() (b []byte, err error) {
	strategy, ok := userMessageStrategies[m.provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}
	return strategy(m)
}
