package models

import (
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/internal/utils"
)

// VideoType 输入内容的类型
type VideoType string

const (
	// VideoTypeText 输入给模型的文本内容
	//
	// 提供商支持: Ark
	VideoTypeText VideoType = "text"
	// VideoTypeImageURL 输入给模型的图片对象
	//
	// 提供商支持: Ark
	VideoTypeImageURL VideoType = "image_url"
)

// VideoRole 图片的位置或用途
type VideoRole string

const (
	// VideoRoleFirstFrame 首帧图片
	//
	// 提供商支持: Ark
	VideoRoleFirstFrame VideoRole = "first_frame"
	// VideoRoleLastFrame 尾帧图片
	//
	// 提供商支持: Ark
	VideoRoleLastFrame VideoRole = "last_frame"
)

// VideoContentImageURL 输入给模型的图片对象
type VideoContentImageURL struct {
	// 图片信息
	//
	// 提供商支持: Ark
	URL string `json:"url" providers:"ark"`
}

// VideoParametersResolution 视频分辨率
type VideoParametersResolution string

const (
	// VideoParametersResolution480P 分辨率480P
	//
	// 提供商支持: Ark
	VideoParametersResolution480P VideoParametersResolution = "480p"
	// VideoParametersResolution720P 分辨率720P
	//
	// 提供商支持: Ark
	VideoParametersResolution720P VideoParametersResolution = "720p"
	// VideoParametersResolution1080P 分辨率1080P
	//
	// 提供商支持: Ark
	VideoParametersResolution1080P VideoParametersResolution = "1080p"
)

// VideoParametersRatio 生成视频的宽高比例
type VideoParametersRatio string

const (
	// VideoParametersRatio21_9 宽高比21:9
	//
	// 提供商支持: Ark
	VideoParametersRatio21_9 VideoParametersRatio = "21:9"
	// VideoParametersRatio16_9 宽高比16:9
	//
	// 提供商支持: Ark
	VideoParametersRatio16_9 VideoParametersRatio = "16:9"
	// VideoParametersRatio4_3 宽高比4:3
	//
	// 提供商支持: Ark
	VideoParametersRatio4_3 VideoParametersRatio = "4:3"
	// VideoParametersRatio1_1 宽高比1:1
	//
	// 提供商支持: Ark
	VideoParametersRatio1_1 VideoParametersRatio = "1:1"
	// VideoParametersRatio3_4 宽高比3:4
	//
	// 提供商支持: Ark
	VideoParametersRatio3_4 VideoParametersRatio = "3:4"
	// VideoParametersRatio9_16 宽高比9:16
	//
	// 提供商支持: Ark
	VideoParametersRatio9_16 VideoParametersRatio = "9:16"
	// VideoParametersRatio9_21 宽高比9:21
	//
	// 提供商支持: Ark
	VideoParametersRatio9_21 VideoParametersRatio = "9:21"
)

// VideoParametersDuration 生成视频时长
type VideoParametersDuration int

const (
	// VideoParametersDuration5  5秒
	//
	// 提供商支持: Ark
	VideoParametersDuration5 VideoParametersDuration = 5
	// VideoParametersDuration10  10秒
	//
	// 提供商支持: Ark
	VideoParametersDuration10 VideoParametersDuration = 10
)

// VideoParametersFPS 帧率
type VideoParametersFPS int

const (
	// VideoParametersFPS16   16帧
	//
	// 提供商支持: Ark
	VideoParametersFPS16 VideoParametersFPS = 16
	// VideoParametersFPS24   24帧
	//
	// 提供商支持: Ark
	VideoParametersFPS24 VideoParametersFPS = 24
)

// VideoParameters 视频参数
type VideoParameters struct {
	// 分辨率
	//
	// 提供商支持: Ark
	Resolution VideoParametersResolution `json:"resolution,omitempty" providers:"ark"`
	// 宽高比
	//
	// 提供商支持: Ark
	Ratio VideoParametersRatio `json:"ratio,omitempty" providers:"ark"`
	// 视频时长
	//
	// 提供商支持: Ark | Alibl
	Duration VideoParametersDuration `json:"duration,omitempty" providers:"ark"`
	// 视频帧率
	//
	// 提供商支持: Ark
	FramePerSecond VideoParametersFPS `json:"framepersecond,omitempty" providers:"ark"`
	// 是否带有水印
	//
	// 提供商支持: Ark
	Watermark bool `json:"watermark,omitempty" providers:"ark"`
	// 种子值
	//
	// 提供商支持: Ark | Alibl
	Seed int `json:"seed,omitempty" providers:"ark,alibl" default:"-1"`
	// 相机是否固定
	//
	// 提供商支持: Ark
	CameraFixed bool `json:"camerafixed,omitempty" providers:"ark"`
	// 是否开启prompt智能改写
	//
	// 提供商支持: Alibl
	PromptExtend bool `json:"prompt_extend,omitempty" providers:"alibl"`
}

type VideoContent struct {
	// 输入内容的类型
	//
	// 提供商支持: Ark
	Type VideoType `json:"type" providers:"ark"`
	// 输入给模型的文本内容
	//
	// 提供商支持: Ark
	Text string `json:"text,omitempty" providers:"ark" parameters:"ark:Parameters"`
	// 模型文本命令
	//
	// 提供商支持: Ark
	Parameters *VideoParameters `json:"-" providers:"ark"`
	// 输入给模型的图片对象
	//
	// 提供商支持: Ark
	ImageURL *VideoContentImageURL `json:"image_url,omitempty" providers:"ark"`
	// 图片的位置或用途
	//
	// 提供商支持: Ark
	Role VideoRole `json:"role,omitempty" providers:"ark"`
}

type VideoRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 模型名称
	//
	// 提供商支持: Ark
	Model string `json:"model,omitempty" providers:"ark"`
	// 图像处理参数
	//
	// 提供商支持: Alibl
	Parameters VideoParameters `json:"parameters,omitempty" providers:"alibl"`
	// 输入给模型，生成视频的信息
	//
	// 提供商支持: Ark
	Content []VideoContent `json:"content,omitempty" providers:"ark"`
	// 回调通知地址
	//
	// 提供商支持: Ark
	CallBackURL string `json:"callback_url,omitempty" providers:"ark"`
}

// MarshalJSON 序列化JSON
func (r VideoRequest) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(r.Provider.String()).Serialize(r)
}

type VideoCreateResponse struct {
	ID string `json:"id,omitempty"` // 任务ID
	httpclient.HttpHeader
}
