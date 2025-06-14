/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-11 14:53:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-13 19:30:03
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/sdkerrors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func listModels(ctx context.Context, client *aisdk.SDKClient) (response models.ListModelsResponse, err error) {
	return client.ListModels(ctx, "system", consts.OpenAI)
}

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, "system", models.ChatRequest{
		ModelInfo: models.ModelInfo{
			Provider:  consts.OpenAI,
			ModelType: consts.ChatModel,
			Model:     consts.OpenAIGPT4oAudioPreview,
		},
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
			},
		},
		MaxCompletionTokens: 4096,
	}, httpclient.WithTimeout(time.Minute*2))
}

func main() {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		log.Fatalf("Failed to create temporary test directory: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
		return
	}
	configData := `{
  "providers": {
    "openai": {
      "base_url": "https://chatapi.onechats.ai/v1",
			"api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("OPENAI_API_KEYS"))
	log.Printf("configData: %s", configData)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		log.Fatalf("Failed to create empty config file: %v", err)
		return
	}

	client, err := aisdk.NewSDKClient(configPath, aisdk.WithDefaultMiddlewares())
	if err != nil {
		log.Fatalf("NewSDKClient() error = %v", err)
		return
	}

	ctx := context.Background()
	// 列出模型
	response1, err := listModels(ctx, client)
	if err != nil {
		log.Fatalf("listModels error = %v, request_id = %s", err, sdkerrors.RequestID(err))
		return
	}
	log.Printf("listModels response: %+v, request_id: %s", response1, response1.RequestID())
	// 创建聊天
	response2, err := createChatCompletion(ctx, client)
	if err != nil {
		fmt.Println("Cause Error =", sdkerrors.Cause(err))
		originalErr := sdkerrors.Unwrap(err)
		switch {
		case errors.Is(originalErr, sdkerrors.ErrProviderNotSupported):
			fmt.Println("ErrProviderNotSupported =", true)
		case errors.Is(originalErr, sdkerrors.ErrModelTypeNotSupported):
			fmt.Println("ErrModelTypeNotSupported =", true)
		case errors.Is(originalErr, sdkerrors.ErrModelNotSupported):
			fmt.Println("ErrModelNotSupported =", true)
		case errors.Is(originalErr, sdkerrors.ErrMethodNotSupported):
			fmt.Println("ErrMethodNotSupported =", true)
		case errors.Is(originalErr, sdkerrors.ErrCompletionStreamNotSupported):
			fmt.Println("ErrCompletionStreamNotSupported =", true)
		case errors.Is(originalErr, context.Canceled):
			fmt.Println("context.Canceled =", true)
		case errors.Is(originalErr, context.DeadlineExceeded):
			fmt.Println("context.DeadlineExceeded =", true)
		}
		log.Fatalf("createChatCompletion error = %v, request_id = %s", err, sdkerrors.RequestID(err))
		return
	}
	log.Printf("createChatCompletion response: %+v, request_id: %s", response2, response2.RequestID())
}
