/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 10:36:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 17:06:49
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package consts

// AliBL 模型名称
//
//	QwQ：基于 Qwen2.5 模型训练的 QwQ 推理模型，通过强化学习大幅度提升了模型推理能力。模型数学代码等核心指标（AIME 24/25、LiveCodeBench）以及部分通用指标（IFEval、LiveBench等）达到DeepSeek-R1 满血版水平
//
//	通义千问-Max：通义千问系列效果最好的模型，适合复杂、多步骤的任务
//
//	通义千问-Plus：能力均衡，推理效果、成本和速度介于通义千问-Max和通义千问-Turbo之间，适合中等复杂任务
//
//	通义千问-Turbo：通义千问系列速度最快、成本极低的模型，适合简单任务
//
//	通义千问-Long：通义千问系列上下文窗口最长，能力均衡且成本较低的模型，适合长文本分析、信息抽取、总结摘要和分类打标等任务
//
//	通义千问Omni：通义千问全新多模态理解生成大模型，支持文本、图像、语音与视频输入，并输出文本与音频，提供了4种自然对话音色
//
//	通义千问Omni-Realtime：相比于通义千问Omni，支持音频的流式输入，且内置 VAD（Voice Activity Detection，语音活动检测）功能，可自动检测用户语音的开始和结束
//
//	QVQ：QVQ是视觉推理模型，支持视觉输入及思维链输出，在数学、编程、视觉分析、创作以及通用任务上都表现了更强的能力
//
//	通义千问VL：通义千问VL是具有视觉（图像）理解能力的文本生成模型，不仅能进行OCR（图片文字识别），还能进一步总结和推理，例如从商品照片中提取属性，根据习题图进行解题等
//
//	通义千问OCR：通义千问OCR模型是专用于文字提取的模型。相较于通义千问VL模型，它更专注于文档、表格、试题、手写体文字等类型图像的文字提取能力。它能够识别多种语言，包括英语、法语、日语、韩语、德语、俄语和意大利语等
//
//	通义千问Audio：通义千问Audio是音频理解模型，支持输入多种音频（人类语音、自然音、音乐、歌声）和文本，并输出文本。该模型不仅能对输入的音频进行转录，还具备更深层次的语义理解、情感分析、音频事件检测、语音聊天等能力
//
//	通义千问ASR：通义千问ASR是基于Qwen-Audio训练，专用于语音识别的模型。目前支持的语言有：中文和英文
//
//	通义千问数学模型：通义千问数学模型是专门用于数学解题的语言模型
//
//	通义千问Coder：通义千问代码模型
//
//	通义千问翻译模型：基于通义千问模型优化的机器翻译大语言模型，擅长中英互译、中文与小语种互译、英文与小语种互译，小语种包括日、韩、法、西、德、葡（巴西）、泰、印尼、越、阿等26种。在多语言互译的基础上，提供术语干预、领域提示、记忆库等能力，提升模型在复杂应用场景下的翻译效果
const (
	// 对话模型
	AliBLQwqPlus                       = "qwq-plus"                            // chat
	AliBLQwqPlusLatest                 = "qwq-plus-latest"                     // chat
	AliBLQwqPlus20250305               = "qwq-plus-2025-03-05"                 // chat
	AliBLQwenMax                       = "qwen-max"                            // chat
	AliBLQwenMaxLatest                 = "qwen-max-latest"                     // chat
	AliBLQwenMax20250125               = "qwen-max-2025-01-25"                 // chat
	AliBLQwenMax20240919               = "qwen-max-2024-09-19"                 // chat
	AliBLQwenMax20240428               = "qwen-max-2024-04-28"                 // chat
	AliBLQwenMax20240403               = "qwen-max-2024-04-03"                 // chat
	AliBLQwenPlus                      = "qwen-plus"                           // chat
	AliBLQwenPlusLatest                = "qwen-plus-latest"                    // chat
	AliBLQwenPlus20250428              = "qwen-plus-2025-04-28"                // chat
	AliBLQwenPlus20250125              = "qwen-plus-2025-01-25"                // chat
	AliBLQwenPlus20250112              = "qwen-plus-2025-01-12"                // chat
	AliBLQwenPlus20241220              = "qwen-plus-2024-12-20"                // chat
	AliBLQwenPlus20241127              = "qwen-plus-2024-11-27"                // chat
	AliBLQwenPlus20241125              = "qwen-plus-2024-11-25"                // chat
	AliBLQwenPlus20240919              = "qwen-plus-2024-09-19"                // chat
	AliBLQwenPlus20240806              = "qwen-plus-2024-08-06"                // chat
	AliBLQwenPlus20240723              = "qwen-plus-2024-07-23"                // chat
	AliBLQwenTurbo                     = "qwen-turbo"                          // chat
	AliBLQwenTurboLatest               = "qwen-turbo-latest"                   // chat
	AliBLQwenTurbo20250428             = "qwen-turbo-2025-04-28"               // chat
	AliBLQwenTurbo20250211             = "qwen-turbo-2025-02-11"               // chat
	AliBLQwenTurbo20241101             = "qwen-turbo-2024-11-01"               // chat
	AliBLQwenTurbo20240919             = "qwen-turbo-2024-09-19"               // chat
	AliBLQwenTurbo20240624             = "qwen-turbo-2024-06-24"               // chat
	AliBLQwenLong                      = "qwen-long"                           // chat
	AliBLQwenLongLatest                = "qwen-long-latest"                    // chat
	AliBLQwenLong20250125              = "qwen-long-2025-01-25"                // chat
	AliBLQwenOmniTurbo                 = "qwen-omni-turbo"                     // chat
	AliBLQwenOmniTurboLatest           = "qwen-omni-turbo-latest"              // chat
	AliBLQwenOmniTurbo20250326         = "qwen-omni-turbo-2025-03-26"          // chat
	AliBLQwenOmniTurbo20250119         = "qwen-omni-turbo-2025-01-19"          // chat
	AliBLQwenOmniTurboRealtime         = "qwen-omni-turbo-realtime"            // chat
	AliBLQwenOmniTurboRealtimeLatest   = "qwen-omni-turbo-realtime-latest"     // chat
	AliBLQwenOmniTurboRealtime20250508 = "qwen-omni-turbo-realtime-2025-05-08" // chat
	AliBLQvqMax                        = "qvq-max"                             // chat
	AliBLQvqMaxLatest                  = "qvq-max-latest"                      // chat
	AliBLQvqMax20250515                = "qvq-max-2025-05-15"                  // chat
	AliBLQvqMax20250325                = "qvq-max-2025-03-25"                  // chat
	AliBLQvqPlus                       = "qvq-plus"                            // chat
	AliBLQvqPlusLatest                 = "qvq-plus-latest"                     // chat
	AliBLQvqPlus20250515               = "qvq-plus-2025-05-15"                 // chat
	AliBLQwenVlMax                     = "qwen-vl-max"                         // chat
	AliBLQwenVlMaxLatest               = "qwen-vl-max-latest"                  // chat
	AliBLQwenVlMax20250408             = "qwen-vl-max-2025-04-08"              // chat
	AliBLQwenVlMax20250402             = "qwen-vl-max-2025-04-02"              // chat
	AliBLQwenVlMax20250125             = "qwen-vl-max-2025-01-25"              // chat
	AliBLQwenVlMax20241230             = "qwen-vl-max-2024-12-30"              // chat
	AliBLQwenVlMax20241119             = "qwen-vl-max-2024-11-19"              // chat
	AliBLQwenVlMax20241030             = "qwen-vl-max-2024-10-30"              // chat
	AliBLQwenVlMax20240809             = "qwen-vl-max-2024-08-09"              // chat
	AliBLQwenVlPlus                    = "qwen-vl-plus"                        // chat
	AliBLQwenVlPlusLatest              = "qwen-vl-plus-latest"                 // chat
	AliBLQwenVlPlus20250507            = "qwen-vl-plus-2025-05-07"             // chat
	AliBLQwenVlPlus20250125            = "qwen-vl-plus-2025-01-25"             // chat
	AliBLQwenVlPlus20250102            = "qwen-vl-plus-2025-01-02"             // chat
	AliBLQwenVlPlus20240809            = "qwen-vl-plus-2024-08-09"             // chat
	AliBLQwenVlPlus20231201            = "qwen-vl-plus-2023-12-01"             // chat
	AliBLQwenVlOcr                     = "qwen-vl-ocr"                         // chat
	AliBLQwenVlOcrLatest               = "qwen-vl-ocr-latest"                  // chat
	AliBLQwenVlOcr20250413             = "qwen-vl-ocr-2025-04-13"              // chat
	AliBLQwenVlOcr20241028             = "qwen-vl-ocr-2024-10-28"              // chat
	AliBLQwenAudioTurbo                = "qwen-audio-turbo"                    // chat
	AliBLQwenAudioTurboLatest          = "qwen-audio-turbo-latest"             // chat
	AliBLQwenAudioTurbo20241204        = "qwen-audio-turbo-2024-12-04"         // chat
	AliBLQwenAudioTurbo20240807        = "qwen-audio-turbo-2024-08-07"         // chat
	AliBLQwenAudioAsr                  = "qwen-audio-asr"                      // chat
	AliBLQwenAudioAsrLatest            = "qwen-audio-asr-latest"               // chat
	AliBLQwenAudioAsr20241204          = "qwen-audio-asr-2024-12-04"           // chat
	AliBLQwenMathPlus                  = "qwen-math-plus"                      // chat
	AliBLQwenMathPlusLatest            = "qwen-math-plus-latest"               // chat
	AliBLQwenMathPlus20240919          = "qwen-math-plus-2024-09-19"           // chat
	AliBLQwenMathPlus20240816          = "qwen-math-plus-2024-08-16"           // chat
	AliBLQwenMathTurbo                 = "qwen-math-turbo"                     // chat
	AliBLQwenMathTurboLatest           = "qwen-math-turbo-latest"              // chat
	AliBLQwenMathTurbo20240919         = "qwen-math-turbo-2024-09-19"          // chat
	AliBLQwenCoderPlus                 = "qwen-coder-plus"                     // chat
	AliBLQwenCoderPlusLatest           = "qwen-coder-plus-latest"              // chat
	AliBLQwenCoderPlus20241106         = "qwen-coder-plus-2024-11-06"          // chat
	AliBLQwenCoderTurbo                = "qwen-coder-turbo"                    // chat
	AliBLQwenCoderTurboLatest          = "qwen-coder-turbo-latest"             // chat
	AliBLQwenCoderTurbo20240919        = "qwen-coder-turbo-2024-09-19"         // chat
	AliBLQwenMtPlus                    = "qwen-mt-plus"                        // chat
	AliBLQwenMtTurbo                   = "qwen-mt-turbo"                       // chat
	AliBLQwen3_235bA22b                = "qwen3-235b-a22b"                     // chat
	AliBLQwen3_32b                     = "qwen3-32b"                           // chat
	AliBLQwen3_30bA3b                  = "qwen3-30b-a3b"                       // chat
	AliBLQwen3_14b                     = "qwen3-14b"                           // chat
	AliBLQwen3_8b                      = "qwen3-8b"                            // chat
	AliBLQwen3_4b                      = "qwen3-4b"                            // chat
	AliBLQwen3_17b                     = "qwen3-1.7b"                          // chat
	AliBLQwen3_06b                     = "qwen3-0.6b"                          // chat
	AliBLQwq32b                        = "qwq-32b"                             // chat
	AliBLQwq32bPreview                 = "qwq-32b-preview"                     // chat
	AliBLQwen2Dot5_14bInstruct1m       = "qwen2.5-14b-instruct-1m"             // chat
	AliBLQwen2Dot5_7bInstruct1m        = "qwen2.5-7b-instruct-1m"              // chat
	AliBLQwen2Dot5_72bInstruct         = "qwen2.5-72b-instruct"                // chat
	AliBLQwen2Dot5_32bInstruct         = "qwen2.5-32b-instruct"                // chat
	AliBLQwen2Dot5_14bInstruct         = "qwen2.5-14b-instruct"                // chat
	AliBLQwen2Dot5_7bInstruct          = "qwen2.5-7b-instruct"                 // chat
	AliBLQwen2Dot5_3bInstruct          = "qwen2.5-3b-instruct"                 // chat
	AliBLQwen2Dot5_15bInstruct         = "qwen2.5-1.5b-instruct"               // chat
	AliBLQwen2Dot5_05bInstruct         = "qwen2.5-0.5b-instruct"               // chat
	AliBLQwen2_72bInstruct             = "qwen2-72b-instruct"                  // chat
	AliBLQwen2_57bA14bInstruct         = "qwen2-57b-a14b-instruct"             // chat
	AliBLQwen2_7bInstruct              = "qwen2-7b-instruct"                   // chat
	AliBLQwen2_15bInstruct             = "qwen2-1.5b-instruct"                 // chat
	AliBLQwen2_05bInstruct             = "qwen2-0.5b-instruct"                 // chat
	AliBLQwen1Dot5_110bChat            = "qwen1.5-110b-chat"                   // chat
	AliBLQwen1Dot5_72bChat             = "qwen1.5-72b-chat"                    // chat
	AliBLQwen1Dot5_32bChat             = "qwen1.5-32b-chat"                    // chat
	AliBLQwen1Dot5_14bChat             = "qwen1.5-14b-chat"                    // chat
	AliBLQwen1Dot5_7bChat              = "qwen1.5-7b-chat"                     // chat
	AliBLQwen1Dot5_18bChat             = "qwen1.5-1.8b-chat"                   // chat
	AliBLQwen1Dot5_05bChat             = "qwen1.5-0.5b-chat"                   // chat
	AliBLQvq72bPreview                 = "qvq-72b-preview"                     // chat
	AliBLQwen2Dot5Omni7b               = "qwen2.5-omni-7b"                     // chat
	AliBLQwen2Dot5Vl72bInstruct        = "qwen2.5-vl-72b-instruct"             // chat
	AliBLQwen2Dot5Vl32bInstruct        = "qwen2.5-vl-32b-instruct"             // chat
	AliBLQwen2Dot5Vl7bInstruct         = "qwen2.5-vl-7b-instruct"              // chat
	AliBLQwen2Dot5Vl3bInstruct         = "qwen2.5-vl-3b-instruct"              // chat
	AliBLQwen2Vl72bInstruct            = "qwen2-vl-72b-instruct"               // chat
	AliBLQwen2Vl7bInstruct             = "qwen2-vl-7b-instruct"                // chat
	AliBLQwen2Vl2bInstruct             = "qwen2-vl-2b-instruct"                // chat
	AliBLQwenVlV1                      = "qwen-vl-v1"                          // chat
	AliBLQwenVlChatV1                  = "qwen-vl-chat-v1"                     // chat
	AliBLQwen2AudioInstruct            = "qwen2-audio-instruct"                // chat
	AliBLQwenAudioChat                 = "qwen-audio-chat"                     // chat
	AliBLQwen2Dot5Math72bInstruct      = "qwen2.5-math-72b-instruct"           // chat
	AliBLQwen2Dot5Math7bInstruct       = "qwen2.5-math-7b-instruct"            // chat
	AliBLQwen2Dot5Math15bInstruct      = "qwen2.5-math-1.5b-instruct"          // chat
	AliBLQwen2Dot5Coder32bInstruct     = "qwen2.5-coder-32b-instruct"          // chat
	AliBLQwen2Dot5Coder14bInstruct     = "qwen2.5-coder-14b-instruct"          // chat
	AliBLQwen2Dot5Coder7bInstruct      = "qwen2.5-coder-7b-instruct"           // chat
	AliBLQwen2Dot5Coder3bInstruct      = "qwen2.5-coder-3b-instruct"           // chat
	AliBLQwen2Dot5Coder15bInstruct     = "qwen2.5-coder-1.5b-instruct"         // chat
	AliBLQwen2Dot5Coder05bInstruct     = "qwen2.5-coder-0.5b-instruct"         // chat
	AliBLDeepSeekR1                    = "deepseek-r1"                         // chat
	AliBLDeepSeekR1_0528               = "deepseek-r1-0528"                    // chat
	AliBLDeepSeekV3                    = "deepseek-v3"                         // chat
	AliBLDeepSeekR1DistillQwen15b      = "deepseek-r1-distill-qwen-1.5b"       // chat
	AliBLDeepSeekR1DistillQwen7b       = "deepseek-r1-distill-qwen-7b"         // chat
	AliBLDeepSeekR1DistillQwen14b      = "deepseek-r1-distill-qwen-14b"        // chat
	AliBLDeepSeekR1DistillQwen32b      = "deepseek-r1-distill-qwen-32b"        // chat
	AliBLDeepSeekR1DistillLlama8b      = "deepseek-r1-distill-llama-8b"        // chat
	AliBLDeepSeekR1DistillLlama70b     = "deepseek-r1-distill-llama-70b"       // chat
	AliBLLlama3Dot3_70bInstruct        = "llama3.3-70b-instruct"               // chat
	AliBLLlama3Dot2_3bInstruct         = "llama3.2-3b-instruct"                // chat
	AliBLLlama3Dot2_1bInstruct         = "llama3.2-1b-instruct"                // chat
	AliBLLlama3Dot1_405bInstruct       = "llama3.1-405b-instruct"              // chat
	AliBLLlama3Dot1_70bInstruct        = "llama3.1-70b-instruct"               // chat
	AliBLLlama3Dot1_8bInstruct         = "llama3.1-8b-instruct"                // chat
	AliBLLlama3_70bInstruct            = "llama3-70b-instruct"                 // chat
	AliBLLlama3_8bInstruct             = "llama3-8b-instruct"                  // chat
	AliBLLlama2_13bChatV2              = "llama2-13b-chat-v2"                  // chat
	AliBLLlama2_7bChatV2               = "llama2-7b-chat-v2"                   // chat
	AliBLLlama4Scout17b16eInstruct     = "llama-4-scout-17b-16e-instruct"      // chat
	AliBLLlama4Maverick17b128eInstruct = "llama-4-maverick-17b-128e-instruct"  // chat
	AliBLLlama3Dot2_90bVisionInstruct  = "llama3.2-90b-vision-instruct"        // chat
	AliBLLlama3Dot2_11bVision          = "llama3.2-11b-vision"                 // chat
	AliBLBaichuan2Turbo                = "baichuan2-turbo"                     // chat
	AliBLBaichuan2_13bChatV1           = "baichuan2-13b-chat-v1"               // chat
	AliBLBaichuan2_7bChatV1            = "baichuan2-7b-chat-v1"                // chat
	AliBLBaichuan7bV1                  = "baichuan-7b-v1"                      // chat
	AliBLChatglm3_6b                   = "chatglm3-6b"                         // chat
	AliBLChatglm6bV2                   = "chatglm-6b-v2"                       // chat
	AliBLYiLarge                       = "yi-large"                            // chat
	AliBLYiMedium                      = "yi-medium"                           // chat
	AliBLYiLargeRag                    = "yi-large-rag"                        // chat
	AliBLYiLargeTurbo                  = "yi-large-turbo"                      // chat
	AliBLAbab6Dot5gChat                = "abab6.5g-chat"                       // chat
	AliBLAbab6Dot5tChat                = "abab6.5t-chat"                       // chat
	AliBLAbab6Dot5sChat                = "abab6.5s-chat"                       // chat
	AliBLZiyaLlama13bV1                = "ziya-llama-13b-v1"                   // chat
	AliBLBelleLlama13b2mV1             = "belle-llama-13b-2m-v1"               // chat
	AliBLChatyuanLargeV2               = "chatyuan-large-v2"                   // chat
	AliBLBilla7bSftV1                  = "billa-7b-sft-v1"                     // chat
)
