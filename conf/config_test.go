/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 19:11:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-26 17:16:03
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package conf_test

import (
	"encoding/json"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestSDKConfigManager(t *testing.T) {
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
	if err := os.WriteFile(configPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}
	// 1. Test NewSDKConfigManager
	manager, err := conf.NewSDKConfigManager(configPath)
	if err != nil {
		t.Fatalf("NewSDKConfigManager failed: %v", err)
	}
	if manager == nil {
		t.Fatal("Config manager should not be nil")
	}
	// 2. Test default configuration values
	config := manager.GetConfig()
	if len(config.Providers) != 0 {
		t.Errorf("Providers should be empty initially, got: %d", len(config.Providers))
	}
	if config.ConnectionOptions.RequestTimeout != 10*time.Second {
		t.Errorf("Default request timeout error, got: %v, want: %v", config.ConnectionOptions.RequestTimeout, 10*time.Second)
	}
	// 3. Test SetProviderConfig and GetProviderConfig
	testProviderConfig := conf.ProviderConfig{
		BaseURL: "https://api.test.com/v1",
		APIKeys: []string{"test-key-1", "test-key-2"},
		OrgID:   "test-org-id",
		Extra: map[string]string{
			"key1": "value1",
		},
	}
	manager.SetProviderConfig(consts.OpenAI, testProviderConfig)
	providerConfig, err := manager.GetProviderConfig(consts.OpenAI)
	if err != nil {
		t.Fatalf("GetProviderConfig failed: %v", err)
	}
	// Verify configuration items
	if providerConfig.BaseURL != testProviderConfig.BaseURL {
		t.Errorf("Provider BaseURL error, got: %v, want: %v", providerConfig.BaseURL, testProviderConfig.BaseURL)
	}
	if !reflect.DeepEqual(providerConfig.APIKeys, testProviderConfig.APIKeys) {
		t.Errorf("Provider APIKeys error, got: %v, want: %v", providerConfig.APIKeys, testProviderConfig.APIKeys)
	}
	if providerConfig.OrgID != testProviderConfig.OrgID {
		t.Errorf("Provider OrgID error, got: %v, want: %v", providerConfig.OrgID, testProviderConfig.OrgID)
	}
	// Verify deep copy - modifying original data should not affect retrieved data
	testProviderConfig.APIKeys[0] = "modified-key"
	providerConfig, _ = manager.GetProviderConfig(consts.OpenAI)
	if providerConfig.APIKeys[0] == "modified-key" {
		t.Error("GetProviderConfig should return a deep copy, should not be affected by original data changes")
	}
	// 4. Test GetProviderConfig error case
	_, err = manager.GetProviderConfig("non-existent-provider")
	if err == nil {
		t.Error("GetProviderConfig should return an error for non-existent provider")
	}
	// 5. Test SetDefaultProvider and GetDefaultProvider
	testProvider := consts.Provider("test-provider")
	manager.SetProviderConfig(testProvider, testProviderConfig) // Ensure provider exists first
	// 6. Test SetConnectionOptions and GetConnectionOptions
	testConnectionOptions := conf.ConnectionOptions{
		RequestTimeout:              15 * time.Second,
		StreamReturnIntervalTimeout: 30 * time.Second,
		MaxRetries:                  5,
		RetryDelay:                  2 * time.Second,
		RetryIncreaseDelay:          true,
		RetryDelayList:              []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
	}
	manager.SetConnectionOptions(testConnectionOptions)
	connectionOptions := manager.GetConnectionOptions()
	if connectionOptions.RequestTimeout != testConnectionOptions.RequestTimeout {
		t.Errorf("ConnectionOptions RequestTimeout error, got: %v, want: %v", connectionOptions.RequestTimeout, testConnectionOptions.RequestTimeout)
	}
	if connectionOptions.MaxRetries != testConnectionOptions.MaxRetries {
		t.Errorf("ConnectionOptions MaxRetries error, got: %v, want: %v", connectionOptions.MaxRetries, testConnectionOptions.MaxRetries)
	}
	if !reflect.DeepEqual(connectionOptions.RetryDelayList, testConnectionOptions.RetryDelayList) {
		t.Errorf("ConnectionOptions RetryDelayList error, got: %v, want: %v", connectionOptions.RetryDelayList, testConnectionOptions.RetryDelayList)
	}
	// 7. Test Save and Load
	err = manager.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file should exist after save")
	}
	// Create new manager and load config
	newManager, err := conf.NewSDKConfigManager(configPath)
	if err != nil {
		t.Fatalf("NewSDKConfigManager with existing config failed: %v", err)
	}
	loadedConfig, err := newManager.GetProviderConfig(testProvider)
	if err != nil {
		t.Fatalf("GetProviderConfig after load failed: %v", err)
	}
	if loadedConfig.BaseURL != testProviderConfig.BaseURL {
		t.Errorf("Provider BaseURL after load error, got: %v, want: %v", loadedConfig.BaseURL, testProviderConfig.BaseURL)
	}
	loadedOptions := newManager.GetConnectionOptions()
	if loadedOptions.RequestTimeout != testConnectionOptions.RequestTimeout {
		t.Errorf("ConnectionOptions RequestTimeout after load error, got: %v, want: %v", loadedOptions.RequestTimeout, testConnectionOptions.RequestTimeout)
	}
	// 8. Test ConnectionOptions JSON serialization and deserialization
	testOptions := conf.ConnectionOptions{
		RequestTimeout:              20 * time.Second,
		StreamReturnIntervalTimeout: 40 * time.Second,
		MaxRetries:                  8,
		RetryDelay:                  3 * time.Second,
		RetryIncreaseDelay:          true,
		RetryDelayList:              []time.Duration{2 * time.Second, 4 * time.Second},
	}
	// Serialize
	jsonData, err := json.Marshal(testOptions)
	if err != nil {
		t.Fatalf("ConnectionOptions serialization failed: %v", err)
	}
	// Deserialize
	var unmarshaledOptions conf.ConnectionOptions
	err = json.Unmarshal(jsonData, &unmarshaledOptions)
	if err != nil {
		t.Fatalf("ConnectionOptions deserialization failed: %v", err)
	}
	// Verify deserialization results
	if unmarshaledOptions.RequestTimeout != testOptions.RequestTimeout {
		t.Errorf("RequestTimeout after deserialization error, got: %v, want: %v", unmarshaledOptions.RequestTimeout, testOptions.RequestTimeout)
	}
	if unmarshaledOptions.MaxRetries != testOptions.MaxRetries {
		t.Errorf("MaxRetries after deserialization error, got: %v, want: %v", unmarshaledOptions.MaxRetries, testOptions.MaxRetries)
	}
	if !reflect.DeepEqual(unmarshaledOptions.RetryDelayList, testOptions.RetryDelayList) {
		t.Errorf("RetryDelayList after deserialization error, got: %v, want: %v", unmarshaledOptions.RetryDelayList, testOptions.RetryDelayList)
	}
	// 9. Test GetConfig and deep copy
	fullConfig := manager.GetConfig()
	providerCfg, ok := fullConfig.Providers[testProvider]
	if !ok {
		t.Fatal("GetConfig should contain test provider configuration")
	}
	if providerCfg.BaseURL != testProviderConfig.BaseURL {
		t.Errorf("GetConfig provider BaseURL error, got: %v, want: %v", providerCfg.BaseURL, testProviderConfig.BaseURL)
	}
	// Verify deep copy - modifying returned config should not affect internal config
	providerCfg = fullConfig.Providers[testProvider]
	providerCfg.BaseURL = "modified-url"
	fullConfig.Providers[testProvider] = providerCfg
	checkConfig, _ := manager.GetProviderConfig(testProvider)
	if checkConfig.BaseURL != testProviderConfig.BaseURL {
		t.Error("GetConfig should return a deep copy, modifications should not affect internal config")
	}
}
