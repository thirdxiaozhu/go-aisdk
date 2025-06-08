package ark

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/models"
)

type ResponseFormat string

const (
	ResponseFormatURL ResponseFormat = "url"
	ResponseFormatB64 ResponseFormat = "b64_json"
)

type ImageSize string

const (
	Size1024x1024 ImageSize = "1024x1024"
	Size864x1152  ImageSize = "864x1152"
	Size1152x864  ImageSize = "1152x864"
	Size1280x720  ImageSize = "1280x720"
	Size720x1280  ImageSize = "720x1280"
	Size832x1248  ImageSize = "832x1248"
	Size1248x832  ImageSize = "1248x832"
	Size1512x648  ImageSize = "1512x648"
)

type ImageRequest struct {
	Prompt string `json:"prompt"` // 提示词
	//Model          string         `json:"model"`  // 模型名称
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
	Size           ImageSize      `json:"size,omitempty"`
	Seed           int            `json:"seed,omitempty"`
	GuidanceScale  float64        `json:"guidance_scale,omitempty"`
	Watermark      bool           `json:"watermark,omitempty"`
}

func (i ImageRequest) GetModelInfo() models.ModelInfo {
	return models.ModelInfo{
		Provider:  consts.Ark,
		ModelType: consts.ImageModel,
		Model:     consts.ArkTextImage,
	}
}

func (i ImageRequest) MarshalJSON() (b []byte, err error) {
	//provider := i.ModelInfo.Provider
	// 创建一个别名结构体
	type Alias ImageRequest
	temp := struct {
		Model string `json:"model"`
		Alias
	}{
		Model: i.GetModelInfo().Model,
		Alias: Alias(i),
	}

	marsheld, err := json.MarshalIndent(temp, "", "  ")
	fmt.Println(string(marsheld))
	return json.Marshal(temp)
}
