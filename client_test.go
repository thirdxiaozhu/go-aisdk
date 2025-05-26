/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-26 17:46:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk_test

import (
	"context"
	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/models"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateChatCompletion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temporary test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")
	// Create the config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	// Create an empty config file to ensure it exists
	configData := `{
  "providers": {
    "deepseek": {
			"base_url": "https://api.deepseek.com",
      "api_keys": ["test-api-key"]
    },
    "openai": {
      "base_url": "https://api.openai.com",
			"api_keys": ["test-openai-key"]
    }
  },
  "connection_options": {
    "request_timeout": "10s",
    "stream_return_interval_timeout": "20s",
    "max_retries": 3,
		"retry_delay": "1s",
		"retry_increase_delay": false
  }
}`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}

	var client *aisdk.SDKClient
	if client, err = aisdk.NewSDKClient(configPath); err != nil {
		t.Fatalf("NewSDKClient() error = %v", err)
	}

	var request models.ChatRequest
	request.ModelInfo = models.ModelInfo{
		Provider:  consts.DeepSeek,
		ModelType: consts.ChatModel,
		Model:     consts.DeepSeekReasoner,
	}
	request.Messages = []models.ChatMessage{
		&models.UserMessage{
			Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
		},
	}

	if _, err = client.CreateChatCompletion(context.Background(), request); err != nil {
		t.Fatalf("CreateChatCompletion() error = %v", err)
	}
}
