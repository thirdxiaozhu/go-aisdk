/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 20:54:21
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-10 14:13:38
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk_test

import (
	"github.com/liusuxian/aisdk"
	"net/http"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	var (
		baseURL   = "https://api.openai.com/v1"
		authToken = "test-token"
		config    = aisdk.DefaultConfig(baseURL, authToken)
	)

	if config.BaseURL != baseURL {
		t.Errorf("Expected base URL to be %s, but got %s", baseURL, config.BaseURL)
	}

	if config.AuthToken != authToken {
		t.Errorf("Expected auth token to be %s, but got %s", authToken, config.AuthToken)
	}

	if config.HTTPClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if config.EmptyMessagesLimit != 300 {
		t.Errorf("Expected empty messages limit to be %d, but got %d", 300, config.EmptyMessagesLimit)
	}
}

func TestClientConfig(t *testing.T) {
	config := aisdk.ClientConfig{
		BaseURL:            "https://custom-api.com",
		AuthToken:          "custom-token",
		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: 500,
	}

	if config.BaseURL != "https://custom-api.com" {
		t.Errorf("Expected base URL to be https://custom-api.com, but got %s", config.BaseURL)
	}

	if config.AuthToken != "custom-token" {
		t.Errorf("Expected auth token to be custom-token, but got %s", config.AuthToken)
	}

	if config.HTTPClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if config.EmptyMessagesLimit != 500 {
		t.Errorf("Expected empty messages limit to be 500, but got %d", config.EmptyMessagesLimit)
	}
}
