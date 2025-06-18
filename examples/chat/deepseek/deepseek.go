/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:15:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-18 15:09:48
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
	"github.com/liusuxian/go-aisdk/internal/utils"
	"github.com/liusuxian/go-aisdk/models"
	"io"
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
	return client.ListModels(ctx, models.ListModelsRequest{
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Provider: consts.DeepSeek,
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Provider: consts.DeepSeek,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
			},
		},
		Model:               consts.DeepSeekChat,
		MaxCompletionTokens: 4096,
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponseStream, err error) {
	return client.CreateChatCompletionStream(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Provider: consts.DeepSeek,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
			},
		},
		Model:               consts.DeepSeekReasoner,
		MaxCompletionTokens: 4096,
		Stream:              true,
	}, httpclient.WithTimeout(time.Minute*5))
}

func main() {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		log.Printf("Failed to create temporary test directory: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v", err)
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
		log.Printf("Failed to create empty config file: %v", err)
		return
	}

	client, err := aisdk.NewSDKClient(configPath, aisdk.WithDefaultMiddlewares())
	if err != nil {
		log.Printf("NewSDKClient() error = %v", err)
		return
	}
	defer func() {
		metrics := client.GetMetrics()
		log.Printf("metrics: %s\n", utils.MustString(metrics))
	}()

	ctx := context.Background()
	// 列出模型
	response1, err := listModels(ctx, client)
	if err != nil {
		log.Printf("listModels error = %v, request_id = %s", err, aisdk.RequestID(err))
		return
	}
	log.Printf("listModels response: %+v, request_id: %s", response1, response1.RequestID())
	// 创建聊天
	response2, err := createChatCompletion(ctx, client)
	if err != nil {
		originalErr := aisdk.Unwrap(err)
		fmt.Println("originalErr =", originalErr)
		fmt.Println("Cause Error =", aisdk.Cause(err))
		switch {
		case errors.Is(originalErr, aisdk.ErrProviderNotSupported):
			fmt.Println("ErrProviderNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrModelTypeNotSupported):
			fmt.Println("ErrModelTypeNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrModelNotSupported):
			fmt.Println("ErrModelNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrMethodNotSupported):
			fmt.Println("ErrMethodNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrCompletionStreamNotSupported):
			fmt.Println("ErrCompletionStreamNotSupported =", true)
		case errors.Is(originalErr, context.Canceled):
			fmt.Println("context.Canceled =", true)
		case errors.Is(originalErr, context.DeadlineExceeded):
			fmt.Println("context.DeadlineExceeded =", true)
		}
		log.Printf("createChatCompletion error = %v, request_id = %s", err, aisdk.RequestID(err))
		return
	}
	log.Printf("createChatCompletion response: %+v, request_id: %s", response2, response2.RequestID())
	log.Printf("createChatCompletion content: %s", response2.Choices[0].Message.Content)
	// 创建流式聊天
	response3, err := createChatCompletionStream(ctx, client)
	if err != nil {
		originalErr := aisdk.Unwrap(err)
		fmt.Println("originalErr =", originalErr)
		fmt.Println("Cause Error =", aisdk.Cause(err))
		switch {
		case errors.Is(originalErr, aisdk.ErrProviderNotSupported):
			fmt.Println("ErrProviderNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrModelTypeNotSupported):
			fmt.Println("ErrModelTypeNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrModelNotSupported):
			fmt.Println("ErrModelNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrMethodNotSupported):
			fmt.Println("ErrMethodNotSupported =", true)
		case errors.Is(originalErr, aisdk.ErrCompletionStreamNotSupported):
			fmt.Println("ErrCompletionStreamNotSupported =", true)
		case errors.Is(originalErr, context.Canceled):
			fmt.Println("context.Canceled =", true)
		case errors.Is(originalErr, context.DeadlineExceeded):
			fmt.Println("context.DeadlineExceeded =", true)
		}
		log.Printf("createChatCompletionStream error = %v, request_id = %s", err, aisdk.RequestID(err))
		return
	}
	defer response3.Close()
	// 读取流式聊天
	usage := models.Usage{}
	for {
		var item models.ChatBaseResponse
		if item, err = response3.StreamReader.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			log.Printf("createChatCompletionStream error = %v, request_id = %s", err, aisdk.RequestID(err))
			break
		}
		usage = item.Usage
		if item.Choices[0].Delta.ReasoningContent != "" {
			log.Printf("createChatCompletionStream reasoning_content: %+v", item.Choices[0].Delta.ReasoningContent)
		}
		if item.Choices[0].Delta.Content != "" {
			log.Printf("createChatCompletionStream content: %+v", item.Choices[0].Delta.Content)
		}
	}
	log.Printf("createChatCompletionStream usage: %+v", usage)
}
