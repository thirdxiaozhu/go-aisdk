package ark

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/models"
	"reflect"
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

func (i ImageRequest) GetUserInfo() models.UserInfo {
	return models.UserInfo{
		UserID: "undefined",
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

type ImageResponseData struct {
	URL string `json:"url"`
	B64 string `json:"b64_json"`
}

type ImageResponseUsgae struct {
	GeneratedImages int `json:"generated_images"`
}

type ImageResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ImageResponse struct {
	Model   string              `json:"model"`
	Created int                 `json:"created"`
	Data    []ImageResponseData `json:"data"`
	Usage   ImageResponseUsgae  `json:"usage"`
	Error   ImageResponseError  `json:"error"`
	httpclient.HttpHeader
}

type VideoRequest struct {
	models.ModelInfo
	Content     []VideoContent `json:"content"`
	CallbackURL string         `json:"callback_url,omitempty"`
}

func (v VideoRequest) GetModelInfo() models.ModelInfo {
	return v.ModelInfo
}

func (v VideoRequest) GetUserInfo() models.UserInfo {
	return models.UserInfo{
		UserID: "undefined",
	}
}

func (v VideoRequest) MarshalJSON() (b []byte, err error) {
	// 创建一个别名结构体
	type Alias VideoRequest
	temp := struct {
		Model string `json:"model"`
		Alias
	}{
		Model: v.GetModelInfo().Model,
		Alias: Alias(v),
	}

	//marsheld, err := json.Marshal(temp)
	//fmt.Println(string(marsheld))
	return json.Marshal(temp)
}

type VideoImageRole string

const (
	VideoImageFirstFrame VideoImageRole = "first_frame"
	VideoImageLastFrame  VideoImageRole = "last_frame"
)

type VideoResolution string

const (
	Resolution_480P  VideoResolution = "480p"
	Resolution_720P  VideoResolution = "720p"
	Resolution_1080P VideoResolution = "1080p"
)

type VideoRatio string

const (
	Ratio_21_9       VideoRatio = "21:9"
	Ratio_16_9       VideoRatio = "16:9"
	Ratio_4_3        VideoRatio = "4_3"
	Ratio_1_1        VideoRatio = "1:1"
	Ratio_3_4        VideoRatio = "3:4"
	Ratio_9_16       VideoRatio = "9:16"
	Ratio_9_21       VideoRatio = "9:21"
	Ratio_KEEP_RATIO VideoRatio = "keep_ratio"
	Ratio_ADAPTIVE   VideoRatio = "adaptive"
)

type VideoDuration int

const (
	Duration_5  VideoDuration = 5
	Duration_10 VideoDuration = 10
)

type VideoFramePerSecond int

const (
	VideoFrame_16 VideoFramePerSecond = 16
	VideoFrame_24 VideoFramePerSecond = 24
)

type VideoParameters struct {
	Resolution     VideoResolution     `json:"resolution"`
	Ratio          VideoRatio          `json:"ratio"`
	Duration       VideoDuration       `json:"duration"`
	FramePerSecond VideoFramePerSecond `json:"framepersecond"`
	WaterMark      bool                `json:"watermark"`
	Seed           int                 `json:"seed"`
	CameraFixed    bool                `json:"camerafixed"`
}

type VideoContent struct {
	Type       models.ChatUserMsgPartType `json:"type"`
	Text       string                     `json:"text"`                 // 文本内容
	Parameters VideoParameters            `json:"parameters,omitempty"` // 多模态内容
	ImageURL   ImageResponseData          `json:"image_url,omitempty"`
	Role       VideoImageRole             `json:"role, omitempty"`
}

func MarshelVideoParas(val reflect.Value) string {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	paras := ""

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		switch field.Kind() {
		case reflect.String:
			para := field.String()
			if para != "" {
				paras = fmt.Sprintf("%s --%s %v", paras, fieldType.Tag.Get("json"), para)
			}
		case reflect.Bool:
			paras = fmt.Sprintf("%s --%s %v", paras, fieldType.Tag.Get("json"), field.Bool())
		case reflect.Int:
			para := field.Int()
			if para != 0 {
				paras = fmt.Sprintf("%s --%s %v", paras, fieldType.Tag.Get("json"), para)
			}
		default:
			return ""
		}

	}
	return paras
}

func (v VideoContent) MarshalJSON() (b []byte, err error) {
	if v.Type == models.ChatUserMsgPartTypeText {
		temp := struct {
			Type models.ChatUserMsgPartType `json:"type"`
			Text string                     `json:"text"`
		}{
			Type: v.Type,
			Text: v.Text + " " + MarshelVideoParas(reflect.ValueOf(v.Parameters)),
		}
		var marsheld []byte
		marsheld, err = json.Marshal(temp)
		fmt.Println(string(marsheld))
		return json.Marshal(temp)
	} else {
		//temp := struct {
		//	Type models.ChatUserMsgPartType `json:"type"`
		//	Text string                     `json:"text"`
		//}{
		//	Type: v.Type,
		//	Text: v.Text,
		//}
		return json.Marshal(v)
	}
	return
}

type VideoResponse struct {
	ID string `json:"id"`
	httpclient.HttpHeader
}
