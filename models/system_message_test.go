/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 20:15:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 22:57:25
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSystemMessage_MarshalJSON(t *testing.T) {
	type fields struct {
		provider string
		Content  string
		Role     string
		Name     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantB   []byte
		wantErr bool
	}{
		{
			name: "test", // openai
			fields: fields{
				provider: "openai",
				Content:  "test",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "test",
			},
			wantB:   []byte(`{"role":"system","content":"test","name":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider: "openai",
				Content:  "",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "test",
			},
			wantB:   []byte(`{"role":"system","name":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider: "openai",
				Content:  "test",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "",
			},
			wantB:   []byte(`{"role":"system","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider: "openai",
				Content:  "",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "",
			},
			wantB:   []byte(`{"role":"system"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek
			fields: fields{
				provider: "deepseek",
				Content:  "test",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "test",
			},
			wantB:   []byte(`{"role":"system","content":"test","name":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider: "deepseek",
				Content:  "",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "test",
			},
			wantB:   []byte(`{"role":"system","name":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider: "deepseek",
				Content:  "test",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "",
			},
			wantB:   []byte(`{"role":"system","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider: "deepseek",
				Content:  "",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "",
			},
			wantB:   []byte(`{"role":"system"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl
			fields: fields{
				provider: "alibl",
				Content:  "test",
				Role:     "",     // 空值，应该使用默认值"system"
				Name:     "test", // alibl不支持name字段，不会序列化
			},
			wantB:   []byte(`{"role":"system","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider: "alibl",
				Content:  "",
				Role:     "",     // 空值，应该使用默认值"system"
				Name:     "test", // alibl不支持name字段，不会序列化
			},
			wantB:   []byte(`{"role":"system"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider: "alibl",
				Content:  "test",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "", // alibl不支持name字段，不会序列化
			},
			wantB:   []byte(`{"role":"system","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider: "alibl",
				Content:  "",
				Role:     "", // 空值，应该使用默认值"system"
				Name:     "", // alibl不支持name字段，不会序列化
			},
			wantB:   []byte(`{"role":"system"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with custom role
			fields: fields{
				provider: "openai",
				Content:  "You are a helpful assistant",
				Role:     "assistant", // 用户自定义的Role值
				Name:     "AI Helper",
			},
			wantB:   []byte(`{"role":"assistant","content":"You are a helpful assistant","name":"AI Helper"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with custom role
			fields: fields{
				provider: "deepseek",
				Content:  "You are an expert programmer",
				Role:     "expert", // 用户自定义的Role值
				Name:     "CodeBot",
			},
			wantB:   []byte(`{"role":"expert","content":"You are an expert programmer","name":"CodeBot"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with custom role
			fields: fields{
				provider: "alibl",
				Content:  "你是一个智能助手",
				Role:     "智能助手",  // 用户自定义的Role值（中文）
				Name:     "不会序列化", // alibl不支持name字段
			},
			wantB:   []byte(`{"role":"智能助手","content":"你是一个智能助手"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with custom role and omitempty
			fields: fields{
				provider: "openai",
				Content:  "",              // 空内容，会omitempty
				Role:     "custom_system", // 用户自定义的Role值
				Name:     "Helper",
			},
			wantB:   []byte(`{"role":"custom_system","name":"Helper"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek providers support test
			fields: fields{
				provider: "deepseek",
				Content:  "Test content",
				Role:     "",        // 默认值
				Name:     "TestBot", // deepseek支持name字段
			},
			wantB:   []byte(`{"role":"system","content":"Test content","name":"TestBot"}`),
			wantErr: false,
		},
		//{
		//	name: "test", // deepseek providers support test
		//	fields: fields{
		//		provider: "ark",
		//		Content:  "Test content",
		//		Role:     "", // 默认值
		//	},
		//	wantB:   []byte(`{"role":"system","content":"Test content","name":"TestBot"}`),
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SystemMessage{
				provider: tt.fields.provider,
				Content:  tt.fields.Content,
				Role:     tt.fields.Role,
				Name:     tt.fields.Name,
			}
			gotB, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("SystemMessage.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("SystemMessage.MarshalJSON() content mismatch:\ngot:  %v\nwant: %v\ngot JSON:  %s\nwant JSON: %s", got, want, gotB, tt.wantB)
			}
		})
	}
}
