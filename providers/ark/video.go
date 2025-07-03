package ark

import (
	"context"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

const (
	apiVideoGenerations = "/contents/generations/tasks"
)

// CreateVideoTask 创建图像
func (s *arkProvider) CreateVideoTask(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = common.ExecuteRequest(ctx, http.MethodPost, s.providerConfig.BaseURL, apiImagesGenerations, opts, s.lb, nil, &response, httpclient.WithBody(request))
	return
}
