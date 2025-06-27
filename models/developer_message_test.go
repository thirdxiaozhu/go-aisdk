/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-26 17:48:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-26 18:07:50
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

func TestDeveloperMessage_MarshalJSON(t *testing.T) {
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
			name: "test", // openai basic
			fields: fields{
				provider: "openai",
				Content:  "请帮助调试这个问题",
				Role:     "", // 空值，应该使用默认值"developer"
			},
			wantB:   []byte(`{"role":"developer","content":"请帮助调试这个问题"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with name
			fields: fields{
				provider: "openai",
				Content:  "这是开发者消息",
				Role:     "", // 空值，应该使用默认值"developer"
				Name:     "debug-assistant",
			},
			wantB:   []byte(`{"role":"developer","content":"这是开发者消息","name":"debug-assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // openai only name
			fields: fields{
				provider: "openai",
				Role:     "", // 空值，应该使用默认值"developer"
				Name:     "developer-bot",
			},
			wantB:   []byte(`{"role":"developer","name":"developer-bot"}`),
			wantErr: false,
		},
		{
			name: "test", // openai empty message (only default)
			fields: fields{
				provider: "openai",
				Role:     "", // 空值，应该使用默认值"developer"
			},
			wantB:   []byte(`{"role":"developer"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with custom role
			fields: fields{
				provider: "openai",
				Content:  "测试内容",
				Role:     "custom-developer", // 用户自定义的Role值
			},
			wantB:   []byte(`{"role":"custom-developer","content":"测试内容"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with explicit empty role
			fields: fields{
				provider: "openai",
				Content:  "测试内容",
				Role:     "", // 显式设置为空，但有default标签会应用默认值
			},
			wantB:   []byte(`{"role":"developer","content":"测试内容"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider: "openai",
				Content:  "", // 空内容，会omitempty
				Role:     "", // 空值，应该使用默认值"developer"
				Name:     "", // 空名称，会omitempty
			},
			wantB:   []byte(`{"role":"developer"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek - not supported
			fields: fields{
				provider: "deepseek",
				Content:  "这是开发者消息",
				Role:     "developer",
				Name:     "assistant",
			},
			wantB:   []byte{}, // 所有字段都不支持deepseek
			wantErr: false,
		},
		{
			name: "test", // alibl - not supported
			fields: fields{
				provider: "alibl",
				Content:  "测试消息",
				Role:     "developer",
			},
			wantB:   []byte{}, // 所有字段都不支持alibl
			wantErr: false,
		},
		{
			name: "test", // unknown provider
			fields: fields{
				provider: "unknown",
				Content:  "测试消息",
				Role:     "developer",
			},
			wantB:   []byte{}, // 所有字段都不支持unknown提供商
			wantErr: false,
		},
		{
			name: "test", // empty provider
			fields: fields{
				provider: "", // 空提供商
				Content:  "测试消息",
				Role:     "developer",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "test", // openai comprehensive test
			fields: fields{
				provider: "openai",
				Content:  "完整的开发者消息测试",
				Role:     "senior-developer",
				Name:     "AI-Developer-Assistant",
			},
			wantB:   []byte(`{"role":"senior-developer","content":"完整的开发者消息测试","name":"AI-Developer-Assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // 特殊字符内容测试
			fields: fields{
				provider: "openai",
				Content:  "包含特殊字符的内容: \"引号\", \\反斜杠\\, \n换行符, \t制表符",
				Role:     "debug-specialist",
				Name:     "Special-Char-Bot",
			},
			wantB:   []byte(`{"role":"debug-specialist","content":"包含特殊字符的内容: \"引号\", \\反斜杠\\, \n换行符, \t制表符","name":"Special-Char-Bot"}`),
			wantErr: false,
		},
		{
			name: "test", // Unicode内容测试
			fields: fields{
				provider: "openai",
				Content:  "Unicode测试: 🚀 emoji, 中文, العربية, русский",
				Role:     "unicode-tester",
				Name:     "🤖AI助手",
			},
			wantB:   []byte(`{"role":"unicode-tester","content":"Unicode测试: 🚀 emoji, 中文, العربية, русский","name":"🤖AI助手"}`),
			wantErr: false,
		},
		{
			name: "test", // 长字符串测试
			fields: fields{
				provider: "openai",
				Content:  "这是一个非常长的字符串测试用例，用来验证序列化器能否正确处理长文本内容。" + strings.Repeat("重复内容", 100),
				Role:     "long-content-processor",
				Name:     "LongStringBot",
			},
			wantB:   []byte(`{"role":"long-content-processor","content":"这是一个非常长的字符串测试用例，用来验证序列化器能否正确处理长文本内容。` + strings.Repeat("重复内容", 100) + `","name":"LongStringBot"}`),
			wantErr: false,
		},
		{
			name: "test", // 提供商名称大小写测试
			fields: fields{
				provider: "OpenAI", // 大写提供商名称，应该被转换为小写
				Content:  "大小写敏感测试",
				Role:     "case-tester",
			},
			wantB:   []byte(`{"role":"case-tester","content":"大小写敏感测试"}`), // 应该正常工作
			wantErr: false,
		},
		{
			name: "test", // 提供商名称混合大小写
			fields: fields{
				provider: "OPENAI", // 全大写提供商名称，应该被转换为小写
				Content:  "全大写提供商测试",
				Role:     "uppercase-tester",
			},
			wantB:   []byte(`{"role":"uppercase-tester","content":"全大写提供商测试"}`), // 应该正常工作
			wantErr: false,
		},
		{
			name: "test", // 空白字符串内容
			fields: fields{
				provider: "openai",
				Content:  "   ", // 只包含空格的内容
				Role:     "whitespace-tester",
				Name:     "   ", // 只包含空格的名称
			},
			wantB:   []byte(`{"role":"whitespace-tester","content":"   ","name":"   "}`),
			wantErr: false,
		},
		{
			name: "test", // 单个字符内容
			fields: fields{
				provider: "openai",
				Content:  "a", // 单个字符
				Role:     "x", // 单个字符role
				Name:     "y", // 单个字符name
			},
			wantB:   []byte(`{"role":"x","content":"a","name":"y"}`),
			wantErr: false,
		},
		{
			name: "test", // 特殊提供商名称格式
			fields: fields{
				provider: "open-ai", // 带连字符的提供商名称
				Content:  "连字符提供商测试",
				Role:     "hyphen-tester",
			},
			wantB:   []byte{}, // 不匹配已知提供商，应该返回空
			wantErr: false,
		},
		{
			name: "test", // 数字内容测试
			fields: fields{
				provider: "openai",
				Content:  "12345.67890",
				Role:     "number-processor",
				Name:     "123",
			},
			wantB:   []byte(`{"role":"number-processor","content":"12345.67890","name":"123"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DeveloperMessage{
				provider: tt.fields.provider,
				Content:  tt.fields.Content,
				Role:     tt.fields.Role,
				Name:     tt.fields.Name,
			}
			gotB, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeveloperMessage.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("DeveloperMessage.MarshalJSON() content mismatch:\ngot:  %v\nwant: %v\ngot JSON:  %s\nwant JSON: %s", got, want, gotB, tt.wantB)
			}
		})
	}
}
