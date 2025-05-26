/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:09:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-26 20:46:37
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk_test

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/models"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getApiKeys(envKey string) (apiKeys string) {
	list := strings.Split(os.Getenv(envKey), ",")
	for i, v := range list {
		if i == 0 {
			apiKeys = fmt.Sprintf(`"%s"`, v)
		} else {
			apiKeys = fmt.Sprintf(`%s,"%s"`, apiKeys, v)
		}
	}
	return
}

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
      "api_keys": [%v]
    },
    "openai": {
      "base_url": "https://newapi.binggokeji.fun/v1",
			"api_keys": [%v]
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
	configData = fmt.Sprintf(configData, getApiKeys("DEEPSEEK_API_KEYS"), getApiKeys("OPENAI_API_KEYS"))
	t.Logf("configData: %s", configData)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}

	var client *aisdk.SDKClient
	if client, err = aisdk.NewSDKClient(configPath); err != nil {
		t.Fatalf("NewSDKClient() error = %v", err)
	}

	var (
		request = models.ChatRequest{
			ModelInfo: models.ModelInfo{
				Provider:  consts.DeepSeek,
				ModelType: consts.ChatModel,
				Model:     consts.DeepSeekChat,
			},
			Messages: []models.ChatMessage{
				&models.UserMessage{
					Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
				},
			},
			MaxCompletionTokens: 4096,
		}
		response models.ChatResponse
	)
	if response, err = client.CreateChatCompletion(context.Background(), request); err != nil {
		t.Fatalf("CreateChatCompletion() error = %v", err)
	}
	t.Logf("response: %+v", response)
}
