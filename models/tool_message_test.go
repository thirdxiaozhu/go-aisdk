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
	"reflect"
	"testing"
)

func TestToolMessage_MarshalJSON(t *testing.T) {
	type fields struct {
		provider   string
		Content    string
		Role       string
		ToolCallID string
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
				provider:   "openai",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","content":"test","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider:   "openai",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider:   "openai",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider:   "openai",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek
			fields: fields{
				provider:   "deepseek",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","content":"test","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider:   "deepseek",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider:   "deepseek",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider:   "deepseek",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl
			fields: fields{
				provider:   "alibl",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","content":"test","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider:   "alibl",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "test",
			},
			wantB:   []byte(`{"role":"tool","tool_call_id":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider:   "alibl",
				Content:    "test",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool","content":"test"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider:   "alibl",
				Content:    "",
				Role:       "", // 空值，应该使用默认值"tool"
				ToolCallID: "",
			},
			wantB:   []byte(`{"role":"tool"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with custom role
			fields: fields{
				provider:   "openai",
				Content:    "test content",
				Role:       "custom_tool", // 用户自定义的Role值
				ToolCallID: "call_123",
			},
			wantB:   []byte(`{"role":"custom_tool","content":"test content","tool_call_id":"call_123"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with custom role
			fields: fields{
				provider:   "deepseek",
				Content:    "result",
				Role:       "function_result", // 用户自定义的Role值
				ToolCallID: "ds_call_456",
			},
			wantB:   []byte(`{"role":"function_result","content":"result","tool_call_id":"ds_call_456"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with custom role and omitempty
			fields: fields{
				provider:   "alibl",
				Content:    "",               // 空内容，会omitempty
				Role:       "assistant_tool", // 用户自定义的Role值
				ToolCallID: "alibl_789",
			},
			wantB:   []byte(`{"role":"assistant_tool","tool_call_id":"alibl_789"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ToolMessage{
				provider:   tt.fields.provider,
				Content:    tt.fields.Content,
				Role:       tt.fields.Role,
				ToolCallID: tt.fields.ToolCallID,
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
