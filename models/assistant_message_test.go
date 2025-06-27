/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 22:58:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-27 10:10:36
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

func TestAssistantMessage_MarshalJSON(t *testing.T) {
	type fields struct {
		provider          string
		Audio             *ChatAssistantMsgAudio
		Content           string
		MultimodalContent []ChatAssistantMsgPart
		Role              string
		Name              string
		Refusal           string
		ToolCalls         []ToolCalls
		Prefix            bool
		ReasoningContent  string
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
				Content:  "Hello, how can I help you?",
				Role:     "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant","content":"Hello, how can I help you?"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with name
			fields: fields{
				provider: "openai",
				Content:  "I'm here to help",
				Role:     "", // 空值，应该使用默认值"assistant"
				Name:     "AI Assistant",
			},
			wantB:   []byte(`{"role":"assistant","content":"I'm here to help","name":"AI Assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with audio
			fields: fields{
				provider: "openai",
				Audio: &ChatAssistantMsgAudio{
					ID: "audio_123",
				},
				Content: "This is a voice message",
				Role:    "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant","audio":{"id":"audio_123"},"content":"This is a voice message"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with refusal
			fields: fields{
				provider: "openai",
				Content:  "I can help with that",
				Role:     "", // 空值，应该使用默认值"assistant"
				Refusal:  "I cannot assist with this request",
			},
			wantB:   []byte(`{"role":"assistant","content":"I can help with that","refusal":"I cannot assist with this request"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with tool calls
			fields: fields{
				provider: "openai",
				Content:  "Let me call a function",
				Role:     "", // 空值，应该使用默认值"assistant"
				ToolCalls: []ToolCalls{
					{
						ID:   "call_123",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "get_weather",
							Arguments: `{"location": "New York"}`,
						},
					},
				},
			},
			wantB:   []byte(`{"role":"assistant","content":"Let me call a function","tool_calls":[{"id":"call_123","type":"function","function":{"name":"get_weather","arguments":"{\"location\": \"New York\"}"}}]}`),
			wantErr: false,
		},
		{
			name: "test", // openai with multimodal content - copyto功能测试
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "This is multimodal text content",
					},
				},
				Role: "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text","text":"This is multimodal text content"}]}`),
			wantErr: false,
		},
		{
			name: "test", // openai with custom role
			fields: fields{
				provider: "openai",
				Content:  "Custom role test",
				Role:     "custom_assistant", // 用户自定义的Role值
			},
			wantB:   []byte(`{"role":"custom_assistant","content":"Custom role test"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with omitempty
			fields: fields{
				provider: "openai",
				Content:  "", // 空内容，会omitempty
				Role:     "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek basic
			fields: fields{
				provider: "deepseek",
				Content:  "DeepSeek response",
				Role:     "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant","content":"DeepSeek response"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with name
			fields: fields{
				provider: "deepseek",
				Content:  "DeepSeek with name",
				Role:     "",
				Name:     "DeepSeek Bot",
			},
			wantB:   []byte(`{"role":"assistant","content":"DeepSeek with name","name":"DeepSeek Bot"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with prefix
			fields: fields{
				provider: "deepseek",
				Content:  "Prefix test",
				Role:     "",
				Prefix:   true,
			},
			wantB:   []byte(`{"role":"assistant","content":"Prefix test","prefix":true}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with reasoning content
			fields: fields{
				provider:         "deepseek",
				Content:          "Final answer",
				Role:             "",
				Prefix:           true,
				ReasoningContent: "Let me think about this step by step...",
			},
			wantB:   []byte(`{"role":"assistant","content":"Final answer","prefix":true,"reasoning_content":"Let me think about this step by step..."}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with omitempty
			fields: fields{
				provider: "deepseek",
				Content:  "",
				Role:     "",
				Prefix:   false, // false值会omitempty
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl basic
			fields: fields{
				provider: "alibl",
				Content:  "AliBL response",
				Role:     "", // 空值，应该使用默认值"assistant"
			},
			wantB:   []byte(`{"role":"assistant","content":"AliBL response"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with prefix (mapping: partial)
			fields: fields{
				provider: "alibl",
				Content:  "AliBL prefix test",
				Role:     "",
				Prefix:   true, // 在alibl中映射为"partial"
			},
			wantB:   []byte(`{"role":"assistant","content":"AliBL prefix test","partial":true}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with name (不支持，不会序列化)
			fields: fields{
				provider: "alibl",
				Content:  "AliBL with name",
				Role:     "",
				Name:     "Won't appear", // alibl不支持name字段
			},
			wantB:   []byte(`{"role":"assistant","content":"AliBL with name"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with omitempty
			fields: fields{
				provider: "alibl",
				Content:  "",
				Role:     "",
				Prefix:   false,
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // openai complex case
			fields: fields{
				provider: "openai",
				Audio: &ChatAssistantMsgAudio{
					ID: "audio_456",
				},
				Content: "Complex test case",
				Role:    "helper",
				Name:    "Advanced AI",
				Refusal: "",
				ToolCalls: []ToolCalls{
					{
						ID:   "call_456",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "complex_function",
							Arguments: `{"param1": "value1", "param2": 42}`,
						},
					},
				},
			},
			wantB:   []byte(`{"role":"helper","audio":{"id":"audio_456"},"content":"Complex test case","name":"Advanced AI","tool_calls":[{"id":"call_456","type":"function","function":{"name":"complex_function","arguments":"{\"param1\": \"value1\", \"param2\": 42}"}}]}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek complex case
			fields: fields{
				provider:         "deepseek",
				Content:          "Complex DeepSeek case",
				Role:             "thinking_assistant",
				Name:             "DeepThink",
				Prefix:           true,
				ReasoningContent: "Complex reasoning process...",
			},
			wantB:   []byte(`{"role":"thinking_assistant","content":"Complex DeepSeek case","name":"DeepThink","prefix":true,"reasoning_content":"Complex reasoning process..."}`),
			wantErr: false,
		},
		{
			name: "test", // multimodal with refusal
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type:    ChatAssistantMsgPartTypeRefusal,
						Refusal: "I cannot process this request",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"refusal","refusal":"I cannot process this request"}]}`),
			wantErr: false,
		},
		{
			name: "test", // mixed multimodal content
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "First part",
					},
					{
						Type:    ChatAssistantMsgPartTypeRefusal,
						Refusal: "Cannot show second part",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text","text":"First part"},{"type":"refusal","refusal":"Cannot show second part"}]}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek with multimodal content (不支持，应该被忽略)
			fields: fields{
				provider: "deepseek",
				Content:  "DeepSeek normal content",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "This should be ignored",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":"DeepSeek normal content"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl with multimodal content (不支持，应该被忽略)
			fields: fields{
				provider: "alibl",
				Content:  "AliBL normal content",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "This should be ignored",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":"AliBL normal content"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with multimodal + original content (multimodal应该覆盖)
			fields: fields{
				provider: "openai",
				Content:  "Original content",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Multimodal content should override",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text","text":"Multimodal content should override"}]}`),
			wantErr: false,
		},
		{
			name: "test", // openai with empty multimodal content
			fields: fields{
				provider:          "openai",
				Content:           "Should keep original content",
				MultimodalContent: []ChatAssistantMsgPart{}, // 空数组
				Role:              "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Should keep original content"}`),
			wantErr: false,
		},
		{
			name: "test", // openai with nil multimodal content
			fields: fields{
				provider:          "openai",
				Content:           "Should keep original content",
				MultimodalContent: nil, // nil切片
				Role:              "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Should keep original content"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek complex with unsupported fields
			fields: fields{
				provider: "deepseek",
				Audio: &ChatAssistantMsgAudio{ // deepseek不支持audio
					ID: "should_be_ignored",
				},
				Content: "DeepSeek content",
				MultimodalContent: []ChatAssistantMsgPart{ // deepseek不支持multimodal
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Should be ignored",
					},
				},
				Role:    "",
				Name:    "DeepSeek Bot",
				Refusal: "Should be ignored", // deepseek不支持refusal
				ToolCalls: []ToolCalls{ // deepseek不支持tool_calls
					{
						ID:   "should_be_ignored",
						Type: "function",
					},
				},
				Prefix:           true,
				ReasoningContent: "DeepSeek reasoning",
			},
			wantB:   []byte(`{"role":"assistant","content":"DeepSeek content","name":"DeepSeek Bot","prefix":true,"reasoning_content":"DeepSeek reasoning"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl complex with unsupported fields
			fields: fields{
				provider: "alibl",
				Audio: &ChatAssistantMsgAudio{ // alibl不支持audio
					ID: "should_be_ignored",
				},
				Content: "AliBL content",
				MultimodalContent: []ChatAssistantMsgPart{ // alibl不支持multimodal
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Should be ignored",
					},
				},
				Role:    "",
				Name:    "Should be ignored", // alibl不支持name
				Refusal: "Should be ignored", // alibl不支持refusal
				ToolCalls: []ToolCalls{ // alibl不支持tool_calls
					{
						ID:   "should_be_ignored",
						Type: "function",
					},
				},
				Prefix:           true,                // alibl支持，但映射为partial
				ReasoningContent: "Should be ignored", // alibl不支持reasoning_content
			},
			wantB:   []byte(`{"role":"assistant","content":"AliBL content","partial":true}`),
			wantErr: false,
		},
		{
			name: "test", // 不支持的提供商
			fields: fields{
				provider: "unsupported_provider",
				Audio: &ChatAssistantMsgAudio{
					ID: "should_be_ignored",
				},
				Content: "Should be ignored",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Should be ignored",
					},
				},
				Role:             "Should be ignored",
				Name:             "Should be ignored",
				Refusal:          "Should be ignored",
				Prefix:           true,
				ReasoningContent: "Should be ignored",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "test", // openai omitempty comprehensive test
			fields: fields{
				provider:          "openai",
				Audio:             nil, // nil指针，会omitempty
				Content:           "",  // 空字符串，会omitempty
				MultimodalContent: nil, // nil切片，不会copyto
				Role:              "",  // 空字符串，使用默认值
				Name:              "",  // 空字符串，会omitempty
				Refusal:           "",  // 空字符串，会omitempty
				ToolCalls:         nil, // nil切片，会omitempty
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek omitempty comprehensive test
			fields: fields{
				provider:         "deepseek",
				Content:          "",    // 空字符串，会omitempty
				Role:             "",    // 空字符串，使用默认值
				Name:             "",    // 空字符串，会omitempty
				Prefix:           false, // false值，会omitempty
				ReasoningContent: "",    // 空字符串，会omitempty
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // alibl omitempty comprehensive test
			fields: fields{
				provider: "alibl",
				Content:  "",    // 空字符串，会omitempty
				Role:     "",    // 空字符串，使用默认值
				Prefix:   false, // false值，会omitempty
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // openai 极限case - multimodal parts with empty fields
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "", // 空文本
					},
					{
						Type:    ChatAssistantMsgPartTypeRefusal,
						Refusal: "", // 空拒绝消息
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text"},{"type":"refusal"}]}`),
			wantErr: false,
		},
		{
			name: "test", // openai multimodal parts with zero values
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{}, // 空结构体
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant"}`),
			wantErr: false,
		},
		{
			name: "test", // openai ToolCalls with nil function
			fields: fields{
				provider: "openai",
				Content:  "Tool call with nil function",
				ToolCalls: []ToolCalls{
					{
						ID:       "call_789",
						Type:     "function",
						Function: nil, // nil函数指针
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Tool call with nil function","tool_calls":[{"id":"call_789","type":"function"}]}`),
			wantErr: false,
		},
		{
			name: "test", // openai ToolCalls with empty function
			fields: fields{
				provider: "openai",
				Content:  "Tool call with empty function",
				ToolCalls: []ToolCalls{
					{
						ID:   "call_890",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "", // 空函数名
							Arguments: "", // 空参数
						},
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Tool call with empty function","tool_calls":[{"id":"call_890","type":"function"}]}`),
			wantErr: false,
		},
		{
			name: "test", // edge case - 空provider
			fields: fields{
				provider: "", // 空提供商
				Content:  "Content with empty provider",
				Role:     "custom_role",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "test", // edge case - provider with special characters
			fields: fields{
				provider: "provider-with-dashes_and_underscores", // 不匹配任何已知提供商
				Content:  "Should be ignored",
				Role:     "Should be ignored",
			},
			wantB:   []byte{},
			wantErr: false,
		},
		{
			name: "test", // openai audio with empty id
			fields: fields{
				provider: "openai",
				Audio: &ChatAssistantMsgAudio{
					ID: "", // 空ID
				},
				Content: "Audio with empty ID",
				Role:    "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Audio with empty ID"}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek prefix true but no reasoning content
			fields: fields{
				provider: "deepseek",
				Content:  "Prefix true without reasoning",
				Role:     "",
				Prefix:   true,
				// ReasoningContent 为空
			},
			wantB:   []byte(`{"role":"assistant","content":"Prefix true without reasoning","prefix":true}`),
			wantErr: false,
		},
		{
			name: "test", // deepseek reasoning without prefix
			fields: fields{
				provider:         "deepseek",
				Content:          "Reasoning without prefix",
				Role:             "",
				Prefix:           false, // false会被omitempty
				ReasoningContent: "Reasoning content without prefix",
			},
			wantB:   []byte(`{"role":"assistant","content":"Reasoning without prefix","reasoning_content":"Reasoning content without prefix"}`),
			wantErr: false,
		},
		{
			name: "test", // copyto edge case - multimodal with both text and refusal in same part (invalid but test robustness)
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type:    ChatAssistantMsgPartTypeText,
						Text:    "Text content",
						Refusal: "Refusal content", // 同时设置两个字段
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text","text":"Text content","refusal":"Refusal content"}]}`),
			wantErr: false,
		},
		{
			name: "test", // large tool calls array
			fields: fields{
				provider: "openai",
				Content:  "Multiple tool calls",
				ToolCalls: []ToolCalls{
					{
						ID:   "call_1",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "func1",
							Arguments: `{"param": "value1"}`,
						},
					},
					{
						ID:   "call_2",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "func2",
							Arguments: `{"param": "value2"}`,
						},
					},
					{
						ID:   "call_3",
						Type: "function",
						Function: &ToolCallsFunction{
							Name:      "func3",
							Arguments: `{"param": "value3"}`,
						},
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":"Multiple tool calls","tool_calls":[{"id":"call_1","type":"function","function":{"name":"func1","arguments":"{\"param\": \"value1\"}"}},{"id":"call_2","type":"function","function":{"name":"func2","arguments":"{\"param\": \"value2\"}"}},{"id":"call_3","type":"function","function":{"name":"func3","arguments":"{\"param\": \"value3\"}"}}]}`),
			wantErr: false,
		},
		{
			name: "test", // multimodal with large array
			fields: fields{
				provider: "openai",
				MultimodalContent: []ChatAssistantMsgPart{
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Part 1",
					},
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Part 2",
					},
					{
						Type:    ChatAssistantMsgPartTypeRefusal,
						Refusal: "Cannot show part 3",
					},
					{
						Type: ChatAssistantMsgPartTypeText,
						Text: "Part 4",
					},
				},
				Role: "",
			},
			wantB:   []byte(`{"role":"assistant","content":[{"type":"text","text":"Part 1"},{"type":"text","text":"Part 2"},{"type":"refusal","refusal":"Cannot show part 3"},{"type":"text","text":"Part 4"}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AssistantMessage{
				provider:          tt.fields.provider,
				Audio:             tt.fields.Audio,
				Content:           tt.fields.Content,
				MultimodalContent: tt.fields.MultimodalContent,
				Role:              tt.fields.Role,
				Name:              tt.fields.Name,
				Refusal:           tt.fields.Refusal,
				ToolCalls:         tt.fields.ToolCalls,
				Prefix:            tt.fields.Prefix,
				ReasoningContent:  tt.fields.ReasoningContent,
			}
			gotB, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssistantMessage.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("AssistantMessage.MarshalJSON() content mismatch:\ngot:  %v\nwant: %v\ngot JSON:  %s\nwant JSON: %s", got, want, gotB, tt.wantB)
			}
		})
	}
}
