/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 13:01:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 00:19:17
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
	"github.com/liusuxian/go-aisdk/utils"
	"log"
	"net"
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

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.AliBL,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "请简要介绍自己",
			},
		},
		Model:               consts.AliBLQwenMax,
		MaxCompletionTokens: utils.Int(4096),
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponseStream, err error) {
	return client.CreateChatCompletionStream(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.AliBL,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "请简要介绍自己",
			},
		},
		Model:               consts.AliBLQwqPlus,
		MaxCompletionTokens: utils.Int(4096),
		Stream:              utils.Bool(true),
		StreamOptions: &models.ChatStreamOptions{
			IncludeUsage: utils.Bool(true),
		},
	}, httpclient.WithTimeout(time.Minute*5), httpclient.WithStreamReturnIntervalTimeout(time.Second*5))
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
    "alibl": {
			"base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
      "api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("ALIBL_API_KEYS"))
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
		log.Printf("metrics = %+v\n", metrics)
	}()

	ctx := context.Background()
	// 创建聊天
	response1, err := createChatCompletion(ctx, client)
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
	log.Printf("createChatCompletion response = %+v, request_id = %s", response1, response1.RequestID())
	// 创建流式聊天
	response2, err := createChatCompletionStream(ctx, client)
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
	// 读取流式聊天
	log.Printf("createChatCompletionStream request_id = %s", response2.RequestID())
	if err = response2.ForEach(func(item models.ChatBaseResponse, isFinished bool) (err error) {
		if isFinished {
			return nil
		}
		log.Printf("createChatCompletionStream item = %+v", item)
		if item.Usage != nil && item.StreamStats != nil {
			log.Printf("createChatCompletionStream usage = %+v", item.Usage)
			log.Printf("createChatCompletionStream stream_stats = %+v", item.StreamStats)
		}
		return nil
	}); err != nil {
		switch {
		case errors.Is(err, aisdk.ErrTooManyEmptyStreamMessages):
			fmt.Println("ErrTooManyEmptyStreamMessages =", true)
		case errors.Is(err, aisdk.ErrStreamReturnIntervalTimeout):
			fmt.Println("ErrStreamReturnIntervalTimeout =", true)
		default:
			var netErr net.Error
			if errors.As(err, &netErr) {
				fmt.Println("net.Error =", true)
			}
		}
		log.Printf("createChatCompletionStream item error = %v", err)
		return
	}
}
