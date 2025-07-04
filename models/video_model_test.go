/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 22:43:06
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"github.com/liusuxian/go-aisdk/consts"
	"reflect"
	"testing"
)

func TestVideoModel_MarshalJSON(t *testing.T) {
	type fields struct {
		provider     string
		Model        string
		VideoContent []VideoContent
	}
	tests := []struct {
		name    string
		fields  fields
		wantB   []byte
		wantErr bool
	}{
		{
			name: "ARK文生视频", // openai
			fields: fields{
				provider: "ark",
				VideoContent: []VideoContent{
					{
						Type: VideoTypeText,
						Text: "生成一个小朋友玩玩具的视频",
						Parameters: &VideoParameters{
							Resolution:     VideoParametersResolution720P,
							Ratio:          VideoParametersRatio9_16,
							Duration:       VideoParametersDuration10,
							FramePerSecond: VideoParametersFPS16,
							Watermark:      true,
							Seed:           -1,
							CameraFixed:    true,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"生成一个小朋友玩玩具的视频 --resolution 720p --ratio 9:16 --duration 10 --framepersecond 16 --watermark true --seed -1 --camerafixed true","type":"text"}],"provider":"ark"}`),
			wantErr: false,
		},
		{
			name: "ARK首帧生成", // openai
			fields: fields{
				provider: "ark",
				VideoContent: []VideoContent{
					{
						Type: VideoTypeText,
						Text: "生成一个小朋友玩玩具的视频",
						Parameters: &VideoParameters{
							Resolution:     VideoParametersResolution720P,
							Ratio:          VideoParametersRatio9_16,
							Duration:       VideoParametersDuration10,
							FramePerSecond: VideoParametersFPS16,
							Watermark:      true,
							Seed:           -1,
							CameraFixed:    true,
						},
					},
					{
						Type: VideoTypeImageURL,
						ImageURL: &VideoContentImageURL{
							URL: "https://example.com/image.png",
						},
						Role: VideoRoleFirstFrame,
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"生成一个小朋友玩玩具的视频 --resolution 720p --ratio 9:16 --duration 10 --framepersecond 16 --watermark true --seed -1 --camerafixed true","type":"text"},{"image_url":{"url":"https://example.com/image.png"},"role":"first_frame","type":"image_url"}],"provider":"ark"}`),
			wantErr: false,
		},
		{
			name: "ARK首尾帧生成视频", // openai
			fields: fields{
				provider: "ark",
				VideoContent: []VideoContent{
					{
						Type: VideoTypeText,
						Text: "生成一个小朋友玩玩具的视频",
						Parameters: &VideoParameters{
							Resolution:     VideoParametersResolution720P,
							Ratio:          VideoParametersRatio9_16,
							Duration:       VideoParametersDuration10,
							FramePerSecond: VideoParametersFPS16,
							Watermark:      true,
							Seed:           -1,
							CameraFixed:    true,
						},
					},
					{
						Type: VideoTypeImageURL,
						ImageURL: &VideoContentImageURL{
							URL: "https://example.com/image.png",
						},
						Role: VideoRoleFirstFrame,
					},
					{
						Type: VideoTypeImageURL,
						ImageURL: &VideoContentImageURL{
							URL: "https://example.com/image2.png",
						},
						Role: VideoRoleLastFrame,
					},
				},
			},
			wantB: []byte(`{"content":[{"text":"生成一个小朋友玩玩具的视频 --resolution 720p --ratio 9:16 --duration 10 --framepersecond 16 --watermark true --seed -1 --camerafixed true","type":"text"},{"image_url":{"url":"https://example.com/image.png"},"role":"first_frame","type":"image_url"},{"image_url":{"url":"https://example.com/image2.png"},"role":"last_frame","type":"image_url"}],"provider":"ark"}
`),
			wantErr: false,
		},
		{
			name: "alibl文生视频", // openai
			fields: fields{
				provider: "alibl",
				VideoContent: []VideoContent{
					{
						Type: VideoTypeText,
						Text: "生成一个小朋友玩玩具的视频",
						Parameters: &VideoParameters{
							Resolution:     VideoParametersResolution720P,
							Ratio:          VideoParametersRatio9_16,
							Duration:       VideoParametersDuration10,
							FramePerSecond: VideoParametersFPS16,
							Watermark:      true,
							Seed:           -1,
							CameraFixed:    true,
						},
					},
				},
			},
			wantB:   []byte(`{"content":[{"text":"生成一个小朋友玩玩具的视频 --resolution 720p --ratio 9:16 --duration 10 --framepersecond 16 --watermark true --seed -1 --camerafixed true","type":"text"}],"provider":"ark"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := VideoRequest{
				Provider: consts.Provider(tt.fields.provider),
				Model:    tt.fields.Model,
				Content:  tt.fields.VideoContent,
			}
			gotB, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToolMessage.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("ToolMessage.MarshalJSON() content mismatch:\ngot:  %v\nwant: %v\ngot JSON:  %s\nwant JSON: %s", got, want, gotB, tt.wantB)
			}
		})
	}
}
