package ark

import "encoding/json"

type ResponseFormat string

const (
	ResponseFormatURL ResponseFormat = "url"
	ResponseFormatB64 ResponseFormat = "b64_json"
)

type ImageSize string

const (
	_1024x1024 ImageSize = "1024x1024"
	_864x1152  ImageSize = "864x1152"
	_1152x864  ImageSize = "1152x864"
	_1280x720  ImageSize = "1280x720"
	_720x1280  ImageSize = "720x1280"
	_832x1248  ImageSize = "832x1248"
	_1248x832  ImageSize = "1248x832"
	_1512x648  ImageSize = "1512x648"
)

type ImageRequest struct {
	Prompt         string         `json:"prompt"` // 提示词
	Model          string         `json:"model"`  // 模型名称
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
	Size           ImageSize      `json:"size,omitempty"`
	Seed           int            `json:"seed,omitempty"`
	GuidanceScale  float64        `json:"guidance_scale,omitempty"`
	Watermark      bool           `json:"watermark,omitempty"`
}

func (i *ImageRequest) MarshalJSON() (b []byte, err error) {
	return json.Marshal(i)
}
