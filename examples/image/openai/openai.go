/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-11 14:53:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 19:59:06
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/internal/utils"
	"github.com/liusuxian/go-aisdk/models"
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

func createImage(ctx context.Context, client *aisdk.SDKClient) (response models.ImageResponse, err error) {
	return client.CreateImage(ctx, models.ImageRequest{
		UserInfo: models.UserInfo{
			UserID: "123456",
		},
		Provider:     consts.OpenAI,
		Prompt:       "一间有着精致雕花窗户的花店，漂亮的深色木质门上挂着铜制把手。店内摆放着各式各样的鲜花，包括玫瑰、百合和向日葵，色彩鲜艳，生机勃勃。背景是温馨的室内场景，透过窗户可以看到街道。高清写实摄影，中景构图。",
		Model:        consts.OpenAIGPTImage1,
		N:            2,
		OutputFormat: models.ImageOutputFormatPNG,
		Quality:      models.ImageQualityHigh,
		Size:         models.ImageSize1536x1024,
	}, httpclient.WithTimeout(time.Minute*5))
}

// saveBase64Image 将base64图片数据保存为文件
func saveBase64Image(base64Data, filename string) (err error) {
	// 解码base64数据
	var imageData []byte
	if imageData, err = base64.StdEncoding.DecodeString(base64Data); err != nil {
		return fmt.Errorf("decode base64 data error: %v", err)
	}
	// 写入文件
	if err = os.WriteFile(filename, imageData, 0644); err != nil {
		return fmt.Errorf("write image file error: %v", err)
	}
	return
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
    "openai": {
      "base_url": "https://chatapi.onechats.ai/v1",
			"api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("OPENAI_API_KEYS"))
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
		log.Printf("metrics = %s\n", utils.MustString(metrics))
	}()

	ctx := context.Background()
	// 创建图像
	response1, err := createImage(ctx, client)
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
		log.Printf("createImage error = %v, request_id = %s", err, aisdk.RequestID(err))
		return
	}
	// 创建images目录来保存图片
	imagesDir := "generated_images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		log.Printf("create images directory error: %v", err)
		return
	}
	// 保存每张生成的图片
	for i, v := range response1.Data {
		if v.B64JSON != "" {
			filename := filepath.Join(imagesDir, fmt.Sprintf("image_%d_%d.png", time.Now().Unix(), i+1))
			if err := saveBase64Image(v.B64JSON, filename); err != nil {
				log.Printf("save image %d error: %v", i+1, err)
			} else {
				log.Printf("save image %d success: %s", i+1, filename)
			}
		} else {
			log.Printf("image %d base64 data is empty", i+1)
		}
	}
	log.Printf("createImage response = %+v, request_id = %s", response1, response1.RequestID())
}
