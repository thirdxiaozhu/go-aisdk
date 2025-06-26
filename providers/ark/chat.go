package ark

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiChatCompletions = "/chat/completions"
)

// CreateChatCompletion 创建聊天
func (s *arkProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, apiChatCompletions, opts, s.lb, nil, &response, withRequestOptions(request)...)
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *arkProvider) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	var stream *httpclient.StreamReader[models.ChatBaseResponse]
	if stream, err = common.ExecuteStreamRequest[models.ChatBaseResponse](ctx, http.MethodPost, s.providerConfig.BaseURL, apiChatCompletions, opts, s.lb, withRequestOptions(request)...); err != nil {
		return
	}
	response = models.ChatResponseStream{
		StreamReader: stream,
	}
	return
}

// withRequestOptions 添加请求选项
func withRequestOptions(request models.ChatRequest) (reqSetters []httpclient.RequestOption) {
	if request.XDashScopeDataInspection != "" {
		reqSetters = []httpclient.RequestOption{
			httpclient.WithKeyValue("X-DashScope-DataInspection", request.XDashScopeDataInspection),
			httpclient.WithBody(request),
		}
	} else {
		reqSetters = []httpclient.RequestOption{
			httpclient.WithBody(request),
		}
	}
	return
}
