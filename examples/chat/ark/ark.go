package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liusuxian/go-aisdk"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
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
	return client.ListModels(ctx, "system", consts.Ark)
}

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, "system", models.ChatRequest{
		ModelInfo: models.ModelInfo{
			Provider:  consts.Ark,
			ModelType: consts.ChatModel,
			Model:     consts.ArkThinkingVersion,
		},
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "立即输出“API”三个字母",
			},
		},
		MaxCompletionTokens: 4096,
	}, httpclient.WithTimeout(time.Minute*2))
}

func streamCallback(response models.ChatResponse) error {
	if response.Choices[0].Delta.Content == "" {
		fmt.Print(response.Choices[0].Delta.ReasoningContent)
	} else {
		fmt.Print(response.Choices[0].Delta.Content)
	}
	return nil
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (interface{}, error) {
	return client.CreateChatCompletionStream(ctx, "system", models.ChatRequest{
		ModelInfo: models.ModelInfo{
			Provider:  consts.Ark,
			ModelType: consts.ChatModel,
			Model:     consts.ArkThinkingVersion,
		},
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "请给出C语言Hello world，最简单的版本就好",
			},
		},
		Stream:              true,
		MaxCompletionTokens: 4096,
		Thinking:            &models.ChatThinking{Type: "enabled"},
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
    },
    "ark": {
      "base_url": "https://ark.cn-beijing.volces.com/api/v3",
			"api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("DEEPSEEK_API_KEYS"), getApiKeys("ARK_API_KEYS"))
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
	_, err = listModels(ctx, client)
	if err != nil {
		log.Printf("listModels error = %v\n", err)
	}

	response2, err := createChatCompletionStream(ctx, client)
	if err != nil {
		log.Fatalf("createChatCompletion error = %v", err)
		return
	}
	log.Printf("createChatCompletion response: %+v", response2)
}
