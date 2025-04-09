/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-09 17:07:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-09 20:45:11
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import "encoding/json"

// UserMessage 用户消息，支持多模态内容
type UserMessage struct {
	Content           string            `json:"content,omitempty"`            // 文本内容
	MultimodalContent []ChatUserMsgPart `json:"multimodal_content,omitempty"` // 多模态内容
	Name              string            `json:"name,omitempty"`               // 参与者名称
}

// GetRole 获取消息角色
func (m UserMessage) GetRole() (role string) { return "user" }

// MarshalJSON 序列化JSON
func (m UserMessage) MarshalJSON() (b []byte, err error) {
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
		temp.Content = m.MultimodalContent
	} else {
		temp.Content = m.Content
	}
	// 移除 multimodal_content 字段
	temp.MultimodalContent = nil
	return json.Marshal(temp)
}

// ChatUserMsgPartType 多模态内容类型
type ChatUserMsgPartType string

const (
	ChatUserMsgPartTypeText       ChatUserMsgPartType = "text"
	ChatUserMsgPartTypeImageURL   ChatUserMsgPartType = "image_url"
	ChatUserMsgPartTypeInputAudio ChatUserMsgPartType = "input_audio"
	ChatUserMsgPartTypeFile       ChatUserMsgPartType = "file"
)

// ChatUserMsgImageURLDetail 图像质量
type ChatUserMsgImageURLDetail string

const (
	ChatUserMsgImageURLDetailHigh ChatUserMsgImageURLDetail = "high"
	ChatUserMsgImageURLDetailLow  ChatUserMsgImageURLDetail = "low"
	ChatUserMsgImageURLDetailAuto ChatUserMsgImageURLDetail = "auto"
)

// ChatUserMsgImageURL 图像URL
type ChatUserMsgImageURL struct {
	URL    string                    `json:"url"`              // 图像URL，支持url和base64编码
	Detail ChatUserMsgImageURLDetail `json:"detail,omitempty"` // 图像质量
}

// ChatUserMsgInputAudioFormat 音频格式
type ChatUserMsgInputAudioFormat string

const (
	ChatUserMsgInputAudioFormatMP3 ChatUserMsgInputAudioFormat = "mp3"
	ChatUserMsgInputAudioFormatWAV ChatUserMsgInputAudioFormat = "wav"
)

// ChatUserMsgInputAudio 输入音频
type ChatUserMsgInputAudio struct {
	Data   string                      `json:"data"`   // 音频数据，支持base64编码
	Format ChatUserMsgInputAudioFormat `json:"format"` // 音频格式
}

// ChatUserMsgFile 文件
type ChatUserMsgFile struct {
	FileData string `json:"file_data,omitempty"` // 文件数据，支持base64编码
	FileID   string `json:"file_id,omitempty"`   // 文件ID
	FileName string `json:"filename,omitempty"`  // 文件名
}

// ChatUserMsgPart 多模态内容
type ChatUserMsgPart struct {
	Type       ChatUserMsgPartType    `json:"type"`                  // 内容类型
	Text       string                 `json:"text,omitempty"`        // 文本内容
	ImageURL   *ChatUserMsgImageURL   `json:"image_url,omitempty"`   // 图像URL
	InputAudio *ChatUserMsgInputAudio `json:"input_audio,omitempty"` // 输入音频
	File       *ChatUserMsgFile       `json:"file,omitempty"`        // 文件
}
