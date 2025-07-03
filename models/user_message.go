/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 14:59:04
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/liusuxian/go-aisdk/internal/utils"

// ChatUserMsgPartType 多模态内容类型
type ChatUserMsgPartType string

const (
	// 提供商支持: OpenAI | AliBL | Ark
	ChatUserMsgPartTypeText ChatUserMsgPartType = "text"
	// 提供商支持: OpenAI | AliBL | Ark
	ChatUserMsgPartTypeImageURL ChatUserMsgPartType = "image_url"
	// 提供商支持: OpenAI | AliBL | Ark
	ChatUserMsgPartTypeInputAudio ChatUserMsgPartType = "input_audio"
	// 提供商支持: OpenAI
	ChatUserMsgPartTypeFile ChatUserMsgPartType = "file"
	// 提供商支持: AliBL
	ChatUserMsgPartTypeVideo ChatUserMsgPartType = "video"
	// 提供商支持: AliBL | Ark
	ChatUserMsgPartTypeVideoURL ChatUserMsgPartType = "video_url"
)

// ChatUserMsgImageURLDetail 图像质量
type ChatUserMsgImageURLDetail string

const (
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailHigh ChatUserMsgImageURLDetail = "high"
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailLow ChatUserMsgImageURLDetail = "low"
	// 提供商支持: OpenAI
	ChatUserMsgImageURLDetailAuto ChatUserMsgImageURLDetail = "auto"
)

// ChatUserMsgImagePixelLimit 像素限制
type ChatUserMsgImagePixelLimit struct {
	// 允许设置图片的像素大小限制
	// 提供商支持: Ark
	MaxPixels int `json:"max_pixels,omitempty" providers:"ark"`
	// 允许设置图片的像素大小限制
	// 提供商支持: Ark
	MinPixels int `json:"min_pixels,omitempty" providers:"ark"`
}

// ChatUserMsgImageURL 图像URL
type ChatUserMsgImageURL struct {
	// 图像URL，支持url和base64编码
	//
	// 提供商支持: OpenAI | AliBL | Ark
	URL string `json:"url,omitempty" providers:"openai,alibl,ark" mapping:"alibl:image"`
	// 图像质量
	//
	// 提供商支持: OpenAI | Ark
	Detail ChatUserMsgImageURLDetail `json:"detail,omitempty" providers:"openai,ark"`
	// 使用OCR模型进行文字提取前对图像进行自动转正
	//
	// 提供商支持: AliBL
	EnableRotate *bool `json:"enable_rotate,omitempty" providers:"alibl"`
	// 使用OCR模型限制输入图像的最小像素，默认值：3136，最小值：100
	//
	// 提供商支持: AliBL
	MinPixels *int `json:"min_pixels,omitempty" providers:"alibl"`
	// 使用OCR模型限制输入图像的最大像素，默认值：6422528，最大值：23520000
	//
	// 提供商支持: AliBL
	MaxPixels *int `json:"max_pixels,omitempty" providers:"alibl"`
	// 图片像素限制
	// 提供商支持: Ark
	ImagePixelLimit *ChatUserMsgImagePixelLimit `json:"image_pixel_limit,omitempty" providers:"ark"`
}

// ChatUserMsgInputAudioFormat 音频格式
type ChatUserMsgInputAudioFormat string

const (
	// 提供商支持: OpenAI
	ChatUserMsgInputAudioFormatMP3 ChatUserMsgInputAudioFormat = "mp3"
	// 提供商支持: OpenAI
	ChatUserMsgInputAudioFormatWAV ChatUserMsgInputAudioFormat = "wav"
)

// ChatUserMsgInputAudio 输入音频
type ChatUserMsgInputAudio struct {
	// 音频数据，支持base64编码，AliBL只支持url
	//
	// 提供商支持: OpenAI | AliBL
	Data string `json:"data,omitempty" providers:"openai,alibl" mapping:"alibl:audio"`
	// 音频格式
	//
	// 提供商支持: OpenAI
	Format ChatUserMsgInputAudioFormat `json:"format,omitempty" providers:"openai"`
}

// ChatUserMsgFile 文件
type ChatUserMsgFile struct {
	// 文件数据，支持base64编码
	//
	// 提供商支持: OpenAI
	FileData string `json:"file_data,omitempty" providers:"openai"`
	// 文件ID
	//
	// 提供商支持: OpenAI
	FileID string `json:"file_id,omitempty" providers:"openai"`
	// 文件名
	//
	// 提供商支持: OpenAI
	FileName string `json:"filename,omitempty" providers:"openai"`
}

// ChatUserMsgInputVideo 输入视频
type ChatUserMsgInputVideo struct {
	// 视频文件，支持url
	//
	// 提供商支持: AliBL
	Video string `json:"video,omitempty" providers:"alibl"`
	// 视频图像列表
	//
	// 提供商支持: AliBL
	VideoImgList []string `json:"video_img_list,omitempty" providers:"alibl" copyto:"Video"`
	// 用于控制抽帧的频率，表示对视频文件每间隔 1/fps 秒抽取一帧，取值范围为 (0.1, 10)，默认值为2.0
	//
	// 提供商支持: AliBL
	FPS *float64 `json:"fps,omitempty" providers:"alibl"`
}

// ChatUserMsgVideoURL 输入视频
type ChatUserMsgVideoURL struct {
	// 视频文件，支持url
	// 提供商支持: Ark
	URL string `json:"url,omitempty" providers:"ark"`
	// 用于控制抽帧的频率，表示对视频文件每间隔 1/fps 秒抽取一帧，取值范围为 (0.2, 5)，默认值为1.0
	// 提供商支持: Ark
	FPS float64 `json:"fps,omitempty" providers:"ark"`
}

// ChatUserMsgPart 多模态内容
type ChatUserMsgPart struct {
	// 内容类型
	//
	// 提供商支持: OpenAI | Ark
	Type ChatUserMsgPartType `json:"type,omitempty" providers:"openai,ark"`
	// 文本内容
	//
	// 提供商支持: OpenAI | AliBL | Ark
	Text string `json:"text,omitempty" providers:"openai,alibl,ark"`
	// 图像URL
	//
	// 提供商支持: OpenAI | AliBL | Ark
	ImageURL *ChatUserMsgImageURL `json:"image_url,omitempty" providers:"openai,alibl,ark" flatten:"alibl"`
	// 输入音频
	//
	// 提供商支持: OpenAI | AliBL
	InputAudio *ChatUserMsgInputAudio `json:"input_audio,omitempty" providers:"openai,alibl" flatten:"alibl"`
	// 文件
	//
	// 提供商支持: OpenAI
	File *ChatUserMsgFile `json:"file,omitempty" providers:"openai"`
	// 输入视频
	//
	// 提供商支持: AliBL
	InputVideo *ChatUserMsgInputVideo `json:"input_video,omitempty" providers:"alibl" flatten:"alibl"`
	// 输入视频
	// 提供商支持: Ark
	VideoURL *ChatUserMsgVideoURL `json:"video_url,omitempty" providers:"ark"`
}

// UserMessage 用户消息，支持多模态内容
type UserMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 文本内容
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Content string `json:"content,omitempty" providers:"openai,deepseek,alibl,ark"`
	// 多模态内容
	//
	// 提供商支持: OpenAI | AliBL | Ark
	MultimodalContent []ChatUserMsgPart `json:"multimodal_content,omitempty" providers:"openai,alibl,ark" copyto:"Content"`
	// 消息角色
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL | Ark
	Role string `json:"role,omitempty" providers:"openai,deepseek,alibl,ark" default:"user"`
	// 参与者名称
	//
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty" providers:"openai,deepseek"`
}

// SetProvider 设置提供商
func (m *UserMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m UserMessage) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(m.provider).Serialize(m)
}
