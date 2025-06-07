package ark

import (
	"context"
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"io"
	"net/http"

	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
)

// arkProvider Ark提供商
type arkProvider struct {
	supportedModels map[consts.ModelType][]string // 支持的模型
	providerConfig  *conf.ProviderConfig          // 提供商配置
	hClient         *httpclient.HTTPClient        // HTTP 客户端
	lb              *loadbalancer.LoadBalancer    // 负载均衡器
}

var (
	arkService *arkProvider // Ark提供商实例
)

const (
	apiChatCompletions = "/chat/completions"
)

// init 包初始化时创建 arkProvider 实例并注册到工厂
func init() {
	arkService = &arkProvider{
		supportedModels: map[consts.ModelType][]string{
			consts.ChatModel: {
				consts.ArkThinkingVersion,
			},
		},
	}
	core.RegisterProvider(consts.Ark, arkService)
}

// GetSupportedModels 获取支持的模型
func (s *arkProvider) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *arkProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.hClient = httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// TODO ListModels 列出模型
func (s *arkProvider) ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (models.ListModelsResponse, error) {

	return models.ListModelsResponse{}, consts.ErrFunctionNotSupported
}

// TODO CreateChatCompletion 创建聊天
func (s *arkProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	var (
		setters = []httpclient.RequestOption{
			httpclient.WithBody(request),
			httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
		}
		req *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, http.MethodPost, s.hClient.FullURL(apiChatCompletions), setters...); err != nil {
		return
	}
	err = s.hClient.SendRequest(req, &response)
	return
}

func (s *arkProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (interface{}, error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var err error
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return nil, err
	}
	var (
		setters = []httpclient.RequestOption{
			httpclient.WithBody(request),
			httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
		}
		req *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, http.MethodPost, s.hClient.FullURL(apiChatCompletions), setters...); err != nil {
		return nil, err
	}

	stream, err := httpclient.SendRequestStream[models.ChatResponse](s.hClient, req)

	defer func() {
		if closeErr := stream.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("stream close error: %w", closeErr)
		}
	}()

	for {
		select {
		case <-ctx.Done(): // 支持上下文取消
			return nil, ctx.Err()
		default:
			var msg models.ChatResponse
			msg, err = stream.Recv()
			switch {
			case errors.Is(err, io.EOF):
				return nil, nil // 正常结束
			case err != nil:
				return nil, err // 错误处理
			default:
				// 使用回调处理消息
				if err = cb(msg); err != nil {
					return nil, err
				}
			}
		}
	}
}
