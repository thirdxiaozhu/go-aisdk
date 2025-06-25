/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-16 13:48:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 12:16:03
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package consts

// OpenAI 模型名称
const (
	// 对话模型
	OpenAIO1Mini                         = "o1-mini"                               // chat
	OpenAIO1Mini20240912                 = "o1-mini-2024-09-12"                    // chat
	OpenAIO1Preview                      = "o1-preview"                            // chat
	OpenAIO1Preview20240912              = "o1-preview-2024-09-12"                 // chat
	OpenAIO1                             = "o1"                                    // chat
	OpenAIO1_20241217                    = "o1-2024-12-17"                         // chat
	OpenAIO1Pro                          = "o1-pro"                                // chat
	OpenAIO1Pro20250319                  = "o1-pro-2025-03-19"                     // chat
	OpenAIO3                             = "o3"                                    // chat
	OpenAIO3_20250416                    = "o3-2025-04-16"                         // chat
	OpenAIO3Mini                         = "o3-mini"                               // chat
	OpenAIO3Mini20250131                 = "o3-mini-2025-01-31"                    // chat
	OpenAIO4Mini                         = "o4-mini"                               // chat
	OpenAIO4Mini20250416                 = "o4-mini-2025-04-16"                    // chat
	OpenAIGPT4_32K0613                   = "gpt-4-32k-0613"                        // chat
	OpenAIGPT4_32K0314                   = "gpt-4-32k-0314"                        // chat
	OpenAIGPT4_32K                       = "gpt-4-32k"                             // chat
	OpenAIGPT4_0613                      = "gpt-4-0613"                            // chat
	OpenAIGPT4_0314                      = "gpt-4-0314"                            // chat
	OpenAIGPT4o                          = "gpt-4o"                                // chat
	OpenAIGPT4o20240513                  = "gpt-4o-2024-05-13"                     // chat
	OpenAIGPT4o20240806                  = "gpt-4o-2024-08-06"                     // chat
	OpenAIGPT4o20241120                  = "gpt-4o-2024-11-20"                     // chat
	OpenAIChatGPT4oLatest                = "chatgpt-4o-latest"                     // chat
	OpenAIGPT4oMini                      = "gpt-4o-mini"                           // chat
	OpenAIGPT4oMini20240718              = "gpt-4o-mini-2024-07-18"                // chat
	OpenAIGPT4oSearchPreview             = "gpt-4o-search-preview"                 // chat
	OpenAIGPT4oSearchPreview20250311     = "gpt-4o-search-preview-2025-03-11"      // chat
	OpenAIGPT4oMiniSearchPreview         = "gpt-4o-mini-search-preview"            // chat
	OpenAIGPT4oMiniSearchPreview20250311 = "gpt-4o-mini-search-preview-2025-03-11" // chat
	OpenAIGPT4Turbo                      = "gpt-4-turbo"                           // chat
	OpenAIGPT4TurboPreview               = "gpt-4-turbo-preview"                   // chat
	OpenAIGPT4Turbo20240409              = "gpt-4-turbo-2024-04-09"                // chat
	OpenAIGPT4_0125Preview               = "gpt-4-0125-preview"                    // chat
	OpenAIGPT4_1106Preview               = "gpt-4-1106-preview"                    // chat
	OpenAIGPT4VisionPreview              = "gpt-4-vision-preview"                  // chat
	OpenAIGPT4                           = "gpt-4"                                 // chat
	OpenAIGPT4Dot1                       = "gpt-4.1"                               // chat
	OpenAIGPT4Dot1_20250414              = "gpt-4.1-2025-04-14"                    // chat
	OpenAIGPT4Dot1Mini                   = "gpt-4.1-mini"                          // chat
	OpenAIGPT4Dot1Mini20250414           = "gpt-4.1-mini-2025-04-14"               // chat
	OpenAIGPT4Dot1Nano                   = "gpt-4.1-nano"                          // chat
	OpenAIGPT4Dot1Nano20250414           = "gpt-4.1-nano-2025-04-14"               // chat
	OpenAIGPT4Dot5Preview                = "gpt-4.5-preview"                       // chat
	OpenAIGPT4Dot5Preview20250227        = "gpt-4.5-preview-2025-02-27"            // chat
	OpenAIGPT3Dot5Turbo0125              = "gpt-3.5-turbo-0125"                    // chat
	OpenAIGPT3Dot5Turbo1106              = "gpt-3.5-turbo-1106"                    // chat
	OpenAIGPT3Dot5Turbo0613              = "gpt-3.5-turbo-0613"                    // chat
	OpenAIGPT3Dot5Turbo0301              = "gpt-3.5-turbo-0301"                    // chat
	OpenAIGPT3Dot5Turbo16k               = "gpt-3.5-turbo-16k"                     // chat
	OpenAIGPT3Dot5Turbo16K0613           = "gpt-3.5-turbo-16k-0613"                // chat
	OpenAIGPT3Dot5Turbo                  = "gpt-3.5-turbo"                         // chat
	OpenAIGPT3Dot5TurboInstruct          = "gpt-3.5-turbo-instruct"                // chat
	OpenAIGPT3Dot5TurboInstruct0914      = "gpt-3.5-turbo-instruct-0914"           // chat
	OpenAIDavinci002                     = "davinci-002"                           // chat
	OpenAIBabbage002                     = "babbage-002"                           // chat
	// 对话 + 音频处理模型
	OpenAIGPT4oAudioPreview                = "gpt-4o-audio-preview"                    // chat, audio
	OpenAIGPT4oAudioPreview20241001        = "gpt-4o-audio-preview-2024-10-01"         // chat, audio
	OpenAIGPT4oAudioPreview20241217        = "gpt-4o-audio-preview-2024-12-17"         // chat, audio
	OpenAIGPT4oAudioPreview20250603        = "gpt-4o-audio-preview-2025-06-03"         // chat, audio
	OpenAIGPT4oRealtimePreview             = "gpt-4o-realtime-preview"                 // chat, audio
	OpenAIGPT4oRealtimePreview20241001     = "gpt-4o-realtime-preview-2024-10-01"      // chat, audio
	OpenAIGPT4oRealtimePreview20241217     = "gpt-4o-realtime-preview-2024-12-17"      // chat, audio
	OpenAIGPT4oRealtimePreview20250603     = "gpt-4o-realtime-preview-2025-06-03"      // chat, audio
	OpenAIGPT4oMiniAudioPreview            = "gpt-4o-mini-audio-preview"               // chat, audio
	OpenAIGPT4oMiniAudioPreview20241217    = "gpt-4o-mini-audio-preview-2024-12-17"    // chat, audio
	OpenAIGPT4oMiniRealtimePreview         = "gpt-4o-mini-realtime-preview"            // chat, audio
	OpenAIGPT4oMiniRealtimePreview20241217 = "gpt-4o-mini-realtime-preview-2024-12-17" // chat, audio
	// 图像生成模型
	OpenAIDallE2    = "dall-e-2"    // image
	OpenAIDallE3    = "dall-e-3"    // image
	OpenAIGPTImage1 = "gpt-image-1" // image
	// 音频处理模型
	OpenAITTS1                = "tts-1"                  // audio
	OpenAITTS1_1106           = "tts-1-1106"             // audio
	OpenAITTS1HD              = "tts-1-hd"               // audio
	OpenAITTS1HD1106          = "tts-1-hd-1106"          // audio
	OpenAIWhisper1            = "whisper-1"              // audio
	OpenAIGPT4oTranscribe     = "gpt-4o-transcribe"      // audio
	OpenAIGPT4oMiniTranscribe = "gpt-4o-mini-transcribe" // audio
	OpenAIGPT4oMiniTTS        = "gpt-4o-mini-tts"        // audio
	// 嵌入模型
	OpenAITextEmbedding3Small = "text-embedding-3-small" // embed
	OpenAITextEmbedding3Large = "text-embedding-3-large" // embed
	OpenAITextEmbeddingAda002 = "text-embedding-ada-002" // embed
	// 内容审核模型
	OpenAIOmniModerationLatest   = "omni-moderation-latest"     // moderation
	OpenAIOmniModeration20240926 = "omni-moderation-2024-09-26" // moderation
)
