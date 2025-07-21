package ark

import (
	"context"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiImagesGenerations = "/images/generations"
)

// CreateImage 创建图像
func (s *arkProvider) CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Provider: consts.Ark,
		Method:   http.MethodPost,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiImagesGenerations,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
		ReqSetters: []httpclient.RequestOption{
			httpclient.WithBody(request),
		},
	})
	return
}
