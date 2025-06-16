/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:15:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-11 14:56:11
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk/sdkerror"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	_ "github.com/liusuxian/go-aisdk/providers/deepseek"
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
	return client.ListModels(ctx, models.ListModelsRequest{
		ModelInfo: models.ModelInfo{
			Provider: consts.DeepSeek,
		},
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, models.ChatRequest{
		ModelInfo: models.ModelInfo{
			Provider:  consts.DeepSeek,
			ModelType: consts.ChatModel,
			Model:     consts.DeepSeekReasoner,
		},
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你是谁",
			},
		},
		MaxCompletionTokens: 4096,
	}, httpclient.WithTimeout(time.Minute*2))
}

func streamCallback(ctx context.Context, response models.ChatResponse) error {
	//if response.Choices[0].Delta.Content == "" {
	//	fmt.Print(response.Choices[0].Delta.ReasoningContent)
	//} else {
	//	fmt.Print(response.Choices[0].Delta.Content)
	//}
	fmt.Println(response)
	return nil
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (httpclient.Response, error) {
	return client.CreateChatCompletionStream(ctx, models.ChatRequest{
		ModelInfo: models.ModelInfo{
			Provider:  consts.DeepSeek,
			ModelType: consts.ChatModel,
			Model:     consts.DeepSeekReasoner,
		},
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "立即输出“API”三个字母",
			},
		},
		Stream:              true,
		MaxCompletionTokens: 4096,
	}, streamCallback, httpclient.WithTimeout(time.Minute*2))
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
    "deepseek": {
			"base_url": "https://api.deepseek.com",
      "api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("DEEPSEEK_API_KEYS"))
	log.Printf("configData: %s", configData)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		log.Fatalf("Failed to create empty config file: %v", err)
		return
	}

	client, err := aisdk.NewSDKClient(configPath, aisdk.WithDefaultMiddlewares())
	if err != nil {
		log.Fatalf("NewSDKClient() sdkerror = %v", err)
		return
	}

	ctx := context.Background()
	// 列出模型
	response1, err := listModels(ctx, client)
	if err != nil {
		log.Fatalf("listModels sdkerror = %v, request_id = %s", err, sdkerror.RequestID(err))
		return
	}
	log.Printf("listModels response: %+v\n", response1.Object)
	log.Printf("listModels response: %+v\n", response1.Data)

	// 创建聊天
	//response2, err := createChatCompletion(ctx, client)
	resp, err := createChatCompletionStream(ctx, client)
	if err != nil {
		originalErr := sdkerror.Unwrap(err)
		fmt.Println("originalErr =", originalErr)
		fmt.Println("Cause Error =", sdkerror.Cause(err))
		switch {
		case errors.Is(originalErr, sdkerror.ErrProviderNotSupported):
			fmt.Println("ErrProviderNotSupported =", true)
		case errors.Is(originalErr, sdkerror.ErrModelTypeNotSupported):
			fmt.Println("ErrModelTypeNotSupported =", true)
		case errors.Is(originalErr, sdkerror.ErrModelNotSupported):
			fmt.Println("ErrModelNotSupported =", true)
		case errors.Is(originalErr, sdkerror.ErrMethodNotSupported):
			fmt.Println("ErrMethodNotSupported =", true)
		case errors.Is(originalErr, sdkerror.ErrCompletionStreamNotSupported):
			fmt.Println("ErrCompletionStreamNotSupported =", true)
		case errors.Is(originalErr, context.Canceled):
			fmt.Println("context.Canceled =", true)
		case errors.Is(originalErr, context.DeadlineExceeded):
			fmt.Println("context.DeadlineExceeded =", true)
		}
		log.Fatalf("createChatCompletion sdkerror = %v, request_id = %s", err, sdkerror.RequestID(err))
		return
	}
	response2 := resp.(*models.ChatStreamResponse)
	//log.Printf("createChatCompletion response: %+v, request_id: %s", response2., response2.RequestID())
	log.Println("response2 Content", response2.Content, len(response2.Content))
	log.Println("response2 ReasoningContent", response2.ReasoningContent, len(response2.ReasoningContent))
}
