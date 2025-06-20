/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-20 23:06:16
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-21 04:38:45
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageReader 图像读取器
type ImageReader struct {
	io.Reader
	filename string
}

// Name 获取文件名
func (r *ImageReader) Name() (filename string) {
	return r.filename
}

// DetectImageType 根据文件内容检测图像类型
func DetectImageType(data []byte) (imageType string) {
	if len(data) < 12 {
		return
	}
	// JPEG: FF D8 FF
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}
	// PNG: 89 50 4E 47 0D 0A 1A 0A
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "png"
	}
	// WebP: 52 49 46 46 [4 bytes] 57 45 42 50
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return "webp"
	}
	// GIF: 47 49 46 38 (GIF8)
	if len(data) >= 4 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return "gif"
	}
	// BMP: 42 4D
	if len(data) >= 2 && data[0] == 0x42 && data[1] == 0x4D {
		return "bmp"
	}
	return
}

// FileToReader 将文件转换为图像读取器
func FileToReader(filename string) (reader *ImageReader, err error) {
	// 读取文件内容
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", filename, err)
		return
	}
	// 检测文件类型
	ext := DetectImageType(data)
	switch ext {
	case "jpeg", "jpg", "png", "webp":
		base := filepath.Base(filename)
		if i := strings.LastIndexByte(base, '.'); i != -1 {
			base = base[:i]
		}
		// 创建图像读取器
		reader = &ImageReader{
			Reader:   bytes.NewReader(data),
			filename: fmt.Sprintf("%s.%s", base, ext),
		}
	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	return
}

// URLToReader 将URL转换为图像读取器
//
//	baseName 文件名，不包括扩展名
func URLToReader(url, baseName string, timeout time.Duration) (reader *ImageReader, err error) {
	// 获取URL数据
	var data []byte
	if data, err = getURLData(url, timeout); err != nil {
		return
	}
	// 检测文件类型
	ext := DetectImageType(data)
	switch ext {
	case "jpeg", "jpg", "png", "webp":
		// 创建图像读取器
		reader = &ImageReader{
			Reader:   bytes.NewReader(data),
			filename: fmt.Sprintf("%s.%s", baseName, ext),
		}
	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	return
}

// Base64ToReader 将base64数据转换为图像读取器
//
//	baseName 文件名，不包括扩展名
func Base64ToReader(b64, baseName string) (reader *ImageReader, err error) {
	b64Data := b64
	if strings.Contains(b64Data, ",") {
		parts := strings.SplitN(b64Data, ",", 2)
		if len(parts) == 2 {
			b64Data = parts[1]
		}
	}
	// 解码Base64数据
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(b64Data); err != nil {
		err = fmt.Errorf("failed to decode base64 data: %w", err)
		return
	}
	// 检测文件类型
	ext := DetectImageType(data)
	switch ext {
	case "jpeg", "jpg", "png", "webp":
		// 创建图像读取器
		reader = &ImageReader{
			Reader:   bytes.NewReader(data),
			filename: fmt.Sprintf("%s.%s", baseName, ext),
		}
	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	return
}

// FileToBase64 文件转Base64编码
func FileToBase64(filename string) (b64 string, err error) {
	// 读取文件内容
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", filename, err)
		return
	}
	// 转换为Base64编码
	return toBase64(data)
}

// URLToBase64 将URL转换为Base64编码
func URLToBase64(url string, timeout time.Duration) (b64 string, err error) {
	// 获取URL数据
	var data []byte
	if data, err = getURLData(url, timeout); err != nil {
		return
	}
	// 转换为Base64编码
	return toBase64(data)
}

// getURLData 获取URL数据
func getURLData(url string, timeout time.Duration) (data []byte, err error) {
	if timeout <= 0 {
		timeout = time.Second * 10
	}
	// 新建请求
	var (
		hc  = &http.Client{Timeout: timeout}
		req *http.Request
	)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}
	// 发送请求
	var resp *http.Response
	if resp, err = hc.Do(req); err != nil {
		err = fmt.Errorf("failed to download image from URL %s: %w", url, err)
		return
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download image from URL %s, status code: %d", url, resp.StatusCode)
		return
	}
	// 将响应体读入内存
	if data, err = io.ReadAll(resp.Body); err != nil {
		err = fmt.Errorf("failed to read image data: %w", err)
		return
	}
	return
}

// toBase64 将数据转换为Base64编码
func toBase64(data []byte) (b64 string, err error) {
	// 转换为 Base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	// 根据文件扩展名返回不同的数据 URI 前缀
	ext := DetectImageType(data)
	switch ext {
	case "jpeg", "jpg":
		b64 = fmt.Sprintf("data:image/jpeg;base64,%s", base64Str)
	case "png":
		b64 = fmt.Sprintf("data:image/png;base64,%s", base64Str)
	case "webp":
		b64 = fmt.Sprintf("data:image/webp;base64,%s", base64Str)
	case "gif":
		b64 = fmt.Sprintf("data:image/gif;base64,%s", base64Str)
	case "bmp":
		b64 = fmt.Sprintf("data:image/bmp;base64,%s", base64Str)
	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	return
}
