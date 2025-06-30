/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-11 14:53:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 12:40:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/utils"
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
			User: "123456",
		},
		Provider: consts.Ark,
		Prompt:   "一间有着精致雕花窗户的花店，漂亮的深色木质门上挂着铜制把手。店内摆放着各式各样的鲜花，包括玫瑰、百合和向日葵，色彩鲜艳，生机勃勃。背景是温馨的室内场景，透过窗户可以看到街道。高清写实摄影，中景构图。",
		Model:    consts.Doubaoseedream3,
		Size:     models.ImageSize1024x1024,
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
    "ark": {
			"base_url": "https://ark.cn-beijing.volces.com/api/v3",
      "api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("ARK_API_KEYS"))
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
	// 创建图像
	response1, err := createImage(ctx, client)
	if err != nil {
		log.Printf("createImage error = %v, request_id = %s", err, aisdk.RequestID(err))
		return
	}

	// 保存每张生成的图片
	filenames := make([]string, 0, len(response1.Data))
	for i, v := range response1.Data {
		if v.B64JSON != "" {
			if filename, err := utils.SaveBase64Image(v.B64JSON, "generated_images", fmt.Sprintf("image_%d", i+1)); err != nil {
				log.Printf("save image %d error: %v", i+1, err)
			} else {
				filenames = append(filenames, filename)
			}
		} else {

			log.Printf("Response = %v", v)
		}
	}
}
