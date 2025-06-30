/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-30 19:40:44
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestUserMessage_MarshalJSON(t *testing.T) {
	type fields struct {
		provider          string
		Content           string
		MultimodalContent []ChatUserMsgPart
		Role              string
		Name              string
	}
	tests := []struct {
		name    string
		fields  fields
		wantB   []byte
		wantErr bool
	}{
		{
			name: "OpenAI完整多模态", // openai
			fields: fields{
				provider: "openai",
				Content:  "test",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "test",
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/image.png",
							EnableRotate: utils.Bool(true),
							MinPixels:    utils.Int(100),
							MaxPixels:    utils.Int(23520000),
						},
					},
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "https://example.com/audio.mp3",
							Format: ChatUserMsgInputAudioFormatMP3,
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileData: "test",
							FileID:   "test",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"text","text":"test"},{"type":"image_url","image_url":{"url":"https://example.com/image.png"}},{"type":"input_audio","input_audio":{"data":"https://example.com/audio.mp3","format":"mp3"}},{"type":"file","file":{"file_data":"test","file_id":"test"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL完整多模态", // alibl
			fields: fields{
				provider: "alibl",
				Content:  "test",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "test",
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/image.png",
							EnableRotate: utils.Bool(true),
							MinPixels:    utils.Int(100),
							MaxPixels:    utils.Int(23520000),
						},
					},
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "https://example.com/audio.mp3",
							Format: ChatUserMsgInputAudioFormatMP3,
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileData: "test",
							FileID:   "test",
						},
					},
					{
						InputVideo: &ChatUserMsgInputVideo{
							Video: "https://example.com/video.mp4",
							VideoImgList: []string{
								"https://example.com/video.mp4",
								"https://example.com/video.mp4",
							},
							FPS: utils.Float64(2.0),
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"test"},{"image":"https://example.com/image.png","enable_rotate":true,"min_pixels":100,"max_pixels":23520000},{"audio":"https://example.com/audio.mp3"},{"video":["https://example.com/video.mp4","https://example.com/video.mp4"],"fps":2}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "Ark完整多模态", // ark
			fields: fields{
				provider: "ark",
				Content:  "test",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "test",
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image.png",
							Detail: ChatUserMsgImageURLDetailHigh,
							ImagePixelLimit: &ChatUserMsgImagePixelLimit{
								MaxPixels: 10000,
								MinPixels: 9999,
							},
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"test","type":"text"},{"image_url":{"detail":"high","image_pixel_limit":{"max_pixels":10000,"min_pixels":9999},"url":"https://example.com/image.png"},"type":"image_url"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "DeepSeek纯文本",
			fields: fields{
				provider: "deepseek",
				Content:  "Hello DeepSeek",
				Role:     "user",
				Name:     "assistant",
			},
			wantB:   []byte(`{"content":"Hello DeepSeek","role":"user","name":"assistant"}`),
			wantErr: false,
		},
		{
			name: "OpenAI纯文本",
			fields: fields{
				provider: "openai",
				Content:  "Hello OpenAI",
				Role:     "",
				Name:     "test_user",
			},
			wantB:   []byte(`{"content":"Hello OpenAI","role":"user","name":"test_user"}`),
			wantErr: false,
		},
		{
			name: "空提供商默认行为",
			fields: fields{
				provider: "",
				Content:  "test content",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "不支持的提供商",
			fields: fields{
				provider: "unsupported",
				Content:  "test content",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "OpenAI图像详细配置",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image.png",
							Detail: ChatUserMsgImageURLDetailHigh,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image.png","detail":"high"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI音频WAV格式",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "base64audiodata",
							Format: ChatUserMsgInputAudioFormatWAV,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"input_audio","input_audio":{"data":"base64audiodata","format":"wav"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI文件完整信息",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileData: "base64filedata",
							FileID:   "file123",
							FileName: "document.pdf",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"file","file":{"file_data":"base64filedata","file_id":"file123","filename":"document.pdf"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL视频最小配置",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{"https://example.com/frame1.jpg"},
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"video":["https://example.com/frame1.jpg"]}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "空的多模态内容",
			fields: fields{
				provider:          "openai",
				MultimodalContent: []ChatUserMsgPart{},
			},
			wantB:   []byte(`{"role":"user"}`),
			wantErr: false,
		},
		{
			name: "nil指针字段",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type:       ChatUserMsgPartTypeImageURL,
						ImageURL:   nil,
						InputAudio: nil,
						File:       nil,
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "零值字段omitempty",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "", // 空字符串应该被omitempty跳过
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"text"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL图像零值跳过",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/image.png",
							EnableRotate: utils.Bool(false),
							MinPixels:    utils.Int(0),
							MaxPixels:    utils.Int(0),
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"image":"https://example.com/image.png","enable_rotate":false,"min_pixels":0,"max_pixels":0}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "混合提供商支持字段",
			fields: fields{
				provider: "deepseek",
				Content:  "test",
				Role:     "user",
				Name:     "test_user",
				MultimodalContent: []ChatUserMsgPart{ // deepseek不支持多模态，应该被跳过
					{
						Type: ChatUserMsgPartTypeText,
						Text: "should be skipped",
					},
				},
			},
			wantB:   []byte(`{"content":"test","role":"user","name":"test_user"}`),
			wantErr: false,
		},
		{
			name: "大小写不敏感的提供商",
			fields: fields{
				provider: "OpenAI",
				Content:  "test case sensitivity",
			},
			wantB:   []byte(`{"content":"test case sensitivity","role":"user"}`),
			wantErr: false,
		},
		{
			name: "特殊字符内容",
			fields: fields{
				provider: "openai",
				Content:  "Hello \"world\" with 'quotes' and\nnewlines\tand\ttabs",
				Name:     "user@test.com",
			},
			wantB:   []byte(`{"content":"Hello \"world\" with 'quotes' and\nnewlines\tand\ttabs","role":"user","name":"user@test.com"}`),
			wantErr: false,
		},
		{
			name: "Unicode内容",
			fields: fields{
				provider: "openai",
				Content:  "你好世界 🌍 こんにちは 🚀",
			},
			wantB:   []byte(`{"content":"你好世界 🌍 こんにちは 🚀","role":"user"}`),
			wantErr: false,
		},
		{
			name: "长内容",
			fields: fields{
				provider: "openai",
				Content:  strings.Repeat("A", 1000), // 1KB的A字符串
			},
			wantB:   []byte(`{"content":"` + strings.Repeat("A", 1000) + `","role":"user"}`),
			wantErr: false,
		},
		{
			name: "SetProvider方法测试",
			fields: fields{
				provider: "deepseek",
				Content:  "test setprovider",
			},
			wantB:   []byte(`{"content":"test setprovider","role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI图像低质量",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image.png",
							Detail: ChatUserMsgImageURLDetailLow,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image.png","detail":"low"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI图像自动质量",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image.png",
							Detail: ChatUserMsgImageURLDetailAuto,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image.png","detail":"auto"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL视频完整配置",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							Video: "https://example.com/video.mp4",
							VideoImgList: []string{
								"https://example.com/frame1.jpg",
								"https://example.com/frame2.jpg",
								"https://example.com/frame3.jpg",
							},
							FPS: utils.Float64(1.5),
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"video":["https://example.com/frame1.jpg","https://example.com/frame2.jpg","https://example.com/frame3.jpg"],"fps":1.5}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI文件只有ID",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileID: "file-abc123",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"file","file":{"file_id":"file-abc123"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI文件只有数据",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileData: "base64encodeddata",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"file","file":{"file_data":"base64encodeddata"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "OpenAI文件只有文件名",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileName: "test.pdf",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"file","file":{"filename":"test.pdf"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "混合内容类型",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "这是文本",
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL: "https://example.com/image.png",
						},
					},
					{
						Type: ChatUserMsgPartTypeText,
						Text: "这是另一个文本",
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"text","text":"这是文本"},{"type":"image_url","image_url":{"url":"https://example.com/image.png"}},{"type":"text","text":"这是另一个文本"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL所有字段最大值",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/image.png",
							EnableRotate: utils.Bool(true),
							MinPixels:    utils.Int(100),
							MaxPixels:    utils.Int(23520000),
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"image":"https://example.com/image.png","enable_rotate":true,"min_pixels":100,"max_pixels":23520000}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL视频边界FPS值",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{"https://example.com/frame.jpg"},
							FPS:          utils.Float64(0.1), // 最小值
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"video":["https://example.com/frame.jpg"],"fps":0.1}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL视频最大FPS值",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{"https://example.com/frame.jpg"},
							FPS:          utils.Float64(10.0), // 最大值
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"video":["https://example.com/frame.jpg"],"fps":10}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "空字符串字段",
			fields: fields{
				provider: "openai",
				Content:  "",
				Role:     "",
				Name:     "",
			},
			wantB:   []byte(`{"role":"user"}`), // 空字符串被omitempty跳过，但role有默认值
			wantErr: false,
		},
		{
			name: "只有Role字段",
			fields: fields{
				provider: "openai",
				Role:     "user",
			},
			wantB:   []byte(`{"role":"user"}`),
			wantErr: false,
		},
		{
			name: "自定义Role",
			fields: fields{
				provider: "openai",
				Content:  "test",
				Role:     "system", // 虽然不常见，但测试非默认值
			},
			wantB:   []byte(`{"content":"test","role":"system"}`),
			wantErr: false,
		},
		{
			name: "URL包含特殊字符",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL: "https://example.com/image with spaces & symbols?.png",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image with spaces & symbols?.png"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "Base64图像数据",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "Base64音频数据",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "data:audio/mp3;base64,SUQzAwAAAAABEFRYWFgAAAASAAADbWFqb3JfYnJhbmQAbXA0MQA=",
							Format: ChatUserMsgInputAudioFormatMP3,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"input_audio","input_audio":{"data":"data:audio/mp3;base64,SUQzAwAAAAABEFRYWFgAAAASAAADbWFqb3JfYnJhbmQAbXA0MQA=","format":"mp3"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "空的视频图像列表",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{}, // 空列表
							FPS:          utils.Float64(2.0),
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"fps":2}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "多种提供商字段混合测试",
			fields: fields{
				provider: "openai",
				Content:  "test content",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "openai text",
						ImageURL: &ChatUserMsgImageURL{ // 这个字段在openai中会被展开
							URL: "https://example.com/image.png",
						},
						InputVideo: &ChatUserMsgInputVideo{ // openai不支持，应该被忽略
							Video: "https://example.com/video.mp4",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"text","text":"openai text","image_url":{"url":"https://example.com/image.png"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "AliBL嵌套字段展开测试",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText, // alibl不支持type字段
						Text: "alibl text",
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/image.png",
							Detail:       ChatUserMsgImageURLDetailHigh, // alibl不支持detail
							EnableRotate: utils.Bool(true),
							MinPixels:    utils.Int(100),
							MaxPixels:    utils.Int(23520000),
						},
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "https://example.com/audio.mp3",
							Format: ChatUserMsgInputAudioFormatMP3, // alibl不支持format
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"alibl text","image":"https://example.com/image.png","enable_rotate":true,"min_pixels":100,"max_pixels":23520000,"audio":"https://example.com/audio.mp3"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "极值测试-最小像素",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						ImageURL: &ChatUserMsgImageURL{
							URL:       "https://example.com/image.png",
							MinPixels: utils.Int(100), // 最小值
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"image":"https://example.com/image.png","min_pixels":100}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "极值测试-最大像素",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						ImageURL: &ChatUserMsgImageURL{
							URL:       "https://example.com/image.png",
							MaxPixels: utils.Int(23520000), // 最大值
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"image":"https://example.com/image.png","max_pixels":23520000}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "所有音频格式测试",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "audio_data_mp3",
							Format: ChatUserMsgInputAudioFormatMP3,
						},
					},
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "audio_data_wav",
							Format: ChatUserMsgInputAudioFormatWAV,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"input_audio","input_audio":{"data":"audio_data_mp3","format":"mp3"}},{"type":"input_audio","input_audio":{"data":"audio_data_wav","format":"wav"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "所有图像质量测试",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image1.png",
							Detail: ChatUserMsgImageURLDetailHigh,
						},
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image2.png",
							Detail: ChatUserMsgImageURLDetailLow,
						},
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/image3.png",
							Detail: ChatUserMsgImageURLDetailAuto,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image1.png","detail":"high"}},{"type":"image_url","image_url":{"url":"https://example.com/image2.png","detail":"low"}},{"type":"image_url","image_url":{"url":"https://example.com/image3.png","detail":"auto"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "复杂嵌套结构测试",
			fields: fields{
				provider: "openai",
				Content:  "Complex nested test",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "Start text",
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "https://example.com/complex.png",
							Detail: ChatUserMsgImageURLDetailHigh,
						},
					},
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "complex_audio_data",
							Format: ChatUserMsgInputAudioFormatWAV,
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileData: "complex_file_data",
							FileID:   "complex_file_id",
							FileName: "complex.pdf",
						},
					},
					{
						Type: ChatUserMsgPartTypeText,
						Text: "End text",
					},
				},
				Role: "user",
				Name: "complex_user",
			},
			wantB:   []byte(`{"content":[{"type":"text","text":"Start text"},{"type":"image_url","image_url":{"url":"https://example.com/complex.png","detail":"high"}},{"type":"input_audio","input_audio":{"data":"complex_audio_data","format":"wav"}},{"type":"file","file":{"file_data":"complex_file_data","file_id":"complex_file_id","filename":"complex.pdf"}},{"type":"text","text":"End text"}],"role":"user","name":"complex_user"}`),
			wantErr: false,
		},
		{
			name: "AliBL复杂嵌套结构测试",
			fields: fields{
				provider: "alibl",
				Content:  "AliBL complex test",
				MultimodalContent: []ChatUserMsgPart{
					{
						Text: "Start text",
						ImageURL: &ChatUserMsgImageURL{
							URL:          "https://example.com/alibl.png",
							EnableRotate: utils.Bool(true),
							MinPixels:    utils.Int(1000),
							MaxPixels:    utils.Int(20000000),
						},
					},
					{
						InputAudio: &ChatUserMsgInputAudio{
							Data: "https://example.com/alibl_audio.mp3",
						},
					},
					{
						InputVideo: &ChatUserMsgInputVideo{
							Video: "https://example.com/alibl_video.mp4",
							VideoImgList: []string{
								"https://example.com/frame1.jpg",
								"https://example.com/frame2.jpg",
								"https://example.com/frame3.jpg",
								"https://example.com/frame4.jpg",
							},
							FPS: utils.Float64(3.5),
						},
					},
					{
						Text: "End text",
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"Start text","image":"https://example.com/alibl.png","enable_rotate":true,"min_pixels":1000,"max_pixels":20000000},{"audio":"https://example.com/alibl_audio.mp3"},{"video":["https://example.com/frame1.jpg","https://example.com/frame2.jpg","https://example.com/frame3.jpg","https://example.com/frame4.jpg"],"fps":3.5},{"text":"End text"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "空字符串和零值混合测试",
			fields: fields{
				provider: "openai",
				Content:  "", // 空内容
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeText,
						Text: "", // 空文本
					},
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL:    "", // 空URL
							Detail: "", // 空详情
						},
					},
					{
						Type: ChatUserMsgPartTypeInputAudio,
						InputAudio: &ChatUserMsgInputAudio{
							Data:   "", // 空数据
							Format: "", // 空格式
						},
					},
				},
				Role: "", // 空角色，应该使用默认值
				Name: "", // 空名称
			},
			wantB:   []byte(`{"content":[{"type":"text"},{"type":"image_url"},{"type":"input_audio"}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "单一字符测试",
			fields: fields{
				provider: "openai",
				Content:  "A",
				Role:     "u",
				Name:     "x",
			},
			wantB:   []byte(`{"content":"A","role":"u","name":"x"}`),
			wantErr: false,
		},
		{
			name: "数字和特殊符号测试",
			fields: fields{
				provider: "openai",
				Content:  "12345!@#$%^&*()_+-=[]{}|;':\",./<>?",
				Name:     "user123",
			},
			wantB:   []byte(`{"content":"12345!@#$%^&*()_+-=[]{}|;':\",./<>?","role":"user","name":"user123"}`),
			wantErr: false,
		},
		{
			name: "多语言混合测试",
			fields: fields{
				provider: "openai",
				Content:  "English 中文 日本語 한국어 العربية Русский Español Français Deutsch",
				Name:     "多语言用户",
			},
			wantB:   []byte(`{"content":"English 中文 日本語 한국어 العربية Русский Español Français Deutsch","role":"user","name":"多语言用户"}`),
			wantErr: false,
		},
		{
			name: "URL编码测试",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL: "https://example.com/image%20with%20spaces.png?param=value&other=data",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/image%20with%20spaces.png?param=value&other=data"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "文件扩展名测试",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileName: "document.pdf",
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileName: "image.jpg",
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileName: "data.csv",
						},
					},
					{
						Type: ChatUserMsgPartTypeFile,
						File: &ChatUserMsgFile{
							FileName: "code.py",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"file","file":{"filename":"document.pdf"}},{"type":"file","file":{"filename":"image.jpg"}},{"type":"file","file":{"filename":"data.csv"}},{"type":"file","file":{"filename":"code.py"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "极长URL测试",
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatUserMsgPart{
					{
						Type: ChatUserMsgPartTypeImageURL,
						ImageURL: &ChatUserMsgImageURL{
							URL: "https://example.com/very/long/path/to/image/with/many/nested/directories/and/a/very/long/filename/that/exceeds/normal/url/length/limits/image.png?param1=value1&param2=value2&param3=value3&param4=value4&param5=value5",
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"type":"image_url","image_url":{"url":"https://example.com/very/long/path/to/image/with/many/nested/directories/and/a/very/long/filename/that/exceeds/normal/url/length/limits/image.png?param1=value1&param2=value2&param3=value3&param4=value4&param5=value5"}}],"role":"user"}`),
			wantErr: false,
		},
		{
			name: "视频帧数组边界测试",
			fields: fields{
				provider: "alibl",
				MultimodalContent: []ChatUserMsgPart{
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{
								"https://example.com/frame1.jpg",
							}, // 只有一个帧
							FPS: utils.Float64(0.1), // 最小FPS
						},
					},
					{
						InputVideo: &ChatUserMsgInputVideo{
							VideoImgList: []string{
								"https://example.com/frame1.jpg",
								"https://example.com/frame2.jpg",
								"https://example.com/frame3.jpg",
								"https://example.com/frame4.jpg",
								"https://example.com/frame5.jpg",
								"https://example.com/frame6.jpg",
								"https://example.com/frame7.jpg",
								"https://example.com/frame8.jpg",
								"https://example.com/frame9.jpg",
								"https://example.com/frame10.jpg",
							}, // 多个帧
							FPS: utils.Float64(10.0), // 最大FPS
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"video":["https://example.com/frame1.jpg"],"fps":0.1},{"video":["https://example.com/frame1.jpg","https://example.com/frame2.jpg","https://example.com/frame3.jpg","https://example.com/frame4.jpg","https://example.com/frame5.jpg","https://example.com/frame6.jpg","https://example.com/frame7.jpg","https://example.com/frame8.jpg","https://example.com/frame9.jpg","https://example.com/frame10.jpg"],"fps":10}],"role":"user"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := UserMessage{
				provider:          tt.fields.provider,
				Content:           tt.fields.Content,
				MultimodalContent: tt.fields.MultimodalContent,
				Role:              tt.fields.Role,
				Name:              tt.fields.Name,
			}
			gotB, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserMessage.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantB) == 0 {
				// 如果期望结果是空字节数组，检查实际结果也应该是空的
				if len(gotB) != 0 {
					t.Errorf("Expected empty result, but got: %s", string(gotB))
				}
				return
			}
			// 解析JSON进行内容比较，而不是字节比较
			var got, want map[string]any
			if err := json.Unmarshal(gotB, &got); err != nil {
				t.Errorf("Failed to unmarshal got JSON: %v", err)
				return
			}
			if err := json.Unmarshal(tt.wantB, &want); err != nil {
				t.Errorf("Failed to unmarshal want JSON: %v", err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("UserMessage.MarshalJSON() content mismatch:\ngot:  %v\nwant: %v\ngot JSON:  %s\nwant JSON: %s", got, want, gotB, tt.wantB)
			}
		})
	}
}
