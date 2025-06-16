package ark

import (
	"context"
	"errors"
	"fmt"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"github.com/liusuxian/go-aisdk/sdkerror"
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
	supportedModels map[fmt.Stringer]map[string]bool // 支持的模型
	providerConfig  *conf.ProviderConfig             // 提供商配置
	hClient         *httpclient.HTTPClient           // HTTP 客户端
	lb              *loadbalancer.LoadBalancer       // 负载均衡器
}

var (
	arkService *arkProvider // Ark提供商实例
)

const (
	apiChatCompletions     = "/chat/completions"
	apiImagesGeneration    = "/images/generations"
	apiConetentsGeneration = "/contents/generations/tasks"
)

// init 包初始化时创建 arkProvider 实例并注册到工厂
func init() {

	arkService = &arkProvider{
		supportedModels: map[fmt.Stringer]map[string]bool{
			consts.ChatModel: {
				consts.ArkThinkingVersion: true,
			},
			consts.ImageModel: {
				consts.ArkTextImage: true,
			},
			consts.VideoModel: {
				consts.ArkTextVideo: true,
			},
		},
	}
	core.RegisterProvider(consts.Ark, arkService)
}

func (s *arkProvider) defaultSetters(req models.Request, setters ...httpclient.RequestOption) ([]httpclient.RequestOption, error) {
	var apiKey *loadbalancer.APIKey
	var err error
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return nil, err
	}
	retSetters := []httpclient.RequestOption{
		httpclient.WithBody(req),
		httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)),
	}

	for _, setter := range setters {
		retSetters = append(retSetters, setter)
	}

	return retSetters, nil
}

func (s *arkProvider) CheckRequestValidation(request models.Request) (err error) {
	return nil
}

// GetSupportedModels 获取支持的模型
func (s *arkProvider) GetSupportedModels() (supportedModels map[fmt.Stringer]map[string]bool) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *arkProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.hClient = httpclient.NewHTTPClient(s.providerConfig.BaseURL)
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// executeRequest 执行请求
func (s *arkProvider) executeRequest(ctx context.Context, method, apiPath string, opts []httpclient.HTTPClientOption, response httpclient.Response, reqSetters ...httpclient.RequestOption) (err error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = s.lb.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(reqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	if req, err = s.hClient.NewRequest(ctx, method, s.hClient.FullURL(apiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = s.hClient.SendRequest(req, response)
	return
}

// ListModels  列出模型
func (s *arkProvider) ListModels(ctx context.Context, opts ...httpclient.HTTPClientOption) (models.ListModelsResponse, error) {
	return models.ListModelsResponse{}, sdkerror.ErrMethodNotSupported
}

// TODO CreateChatCompletion 创建聊天
func (s *arkProvider) CreateChatCompletion(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 设置客户端选项
	err = s.executeRequest(ctx, http.MethodPost, apiChatCompletions, opts, &response, httpclient.WithBody(request))
	return
}

func (s *arkProvider) CreateChatCompletionStream(ctx context.Context, request models.Request, cb core.StreamCallback, opts ...httpclient.HTTPClientOption) (interface{}, error) {
	// 设置客户端选项
	for _, opt := range opts {
		opt(s.hClient)
	}
	var setters []httpclient.RequestOption
	var req *http.Request
	setters, err := s.defaultSetters(request)
	if req, err = s.hClient.NewRequest(ctx, http.MethodPost, s.hClient.FullURL(apiChatCompletions), setters...); err != nil {
		return nil, err
	}

	stream, err := httpclient.SendRequestStream[models.ChatResponse](s.hClient, req)

	defer func() {
		if closeErr := stream.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("stream close sdkerror: %w", closeErr)
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

func (s *arkProvider) CreateImageGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (httpclient.Response, error) {
	var resp ImageResponse
	err := s.executeRequest(ctx, http.MethodPost, apiImagesGeneration, opts, &resp, httpclient.WithBody(request))
	return &resp, err
}

func (s *arkProvider) CreateVideoGeneration(ctx context.Context, request models.Request, opts ...httpclient.HTTPClientOption) (httpclient.Response, error) {
	//var err sdkerror
	//for _, opt := range opts {
	//	opt(s.hClient)
	//}
	//// 获取一个APIKey
	//var setters []httpclient.RequestOption
	//var req *http.Request
	//setters, err = s.defaultSetters(request)
	//if req, err = s.hClient.NewRequest(ctx, http.MethodPost, s.hClient.FullURL(apiConetentsGeneration), setters...); err != nil {
	//	return nil, err
	//}
	//
	//var resp VideoResponse
	//
	//if err = s.hClient.SendRequest(req, &resp); err != nil {
	//	return nil, err
	//}
	//return &resp, err
	var resp VideoResponse
	err := s.executeRequest(ctx, http.MethodPost, apiConetentsGeneration, opts, &resp, httpclient.WithBody(request))
	return &resp, err
}
