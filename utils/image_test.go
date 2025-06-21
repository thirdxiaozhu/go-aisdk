/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-20 23:06:16
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-21 04:52:26
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils_test

import (
	"fmt"
	"github.com/liusuxian/go-aisdk/utils"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
	"time"
)

func createTempImage(ext string) (tempFile string, err error) {
	tempFile = fmt.Sprintf("test_temp.%s", ext)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	red := &image.Uniform{C: color.RGBA{255, 0, 0, 255}}
	draw.Draw(img, image.Rect(25, 25, 75, 75), red, image.Point{}, draw.Src)

	var file *os.File
	if file, err = os.Create(tempFile); err != nil {
		return
	}
	defer file.Close()

	switch ext {
	case "png":
		err = png.Encode(file, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(file, img, nil)
	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	return
}

func TestFileToReader(t *testing.T) {
	var (
		reader *utils.ImageReader
		err    error
	)
	reader, err = utils.FileToReader("")
	if err == nil {
		t.Error("FileToReader() error = nil")
	}
	if reader != nil {
		t.Error("FileToReader() reader = ", reader)
	}

	reader, err = utils.FileToReader("nonexistent.png")
	if err == nil {
		t.Error("FileToReader() error = nil")
	}
	if reader != nil {
		t.Error("FileToReader() reader = ", reader)
	}

	reader, err = utils.FileToReader("image.go")
	if err == nil {
		t.Error("FileToReader() error = nil")
	}
	if reader != nil {
		t.Error("FileToReader() reader = ", reader)
	}

	tempFile, err := createTempImage("png")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	reader, err = utils.FileToReader(tempFile)
	if err != nil {
		t.Error("FileToReader() error = ", err)
	}
	if reader == nil {
		t.Error("FileToReader() reader = nil")
	}
	os.Remove(tempFile)

	tempFile, err = createTempImage("jpg")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	reader, err = utils.FileToReader(tempFile)
	if err != nil {
		t.Error("FileToReader() error = ", err)
	}
	if reader == nil {
		t.Error("FileToReader() reader = nil")
	}
	os.Remove(tempFile)

	tempFile, err = createTempImage("jpeg")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	reader, err = utils.FileToReader(tempFile)
	if err != nil {
		t.Error("FileToReader() error = ", err)
	}
	if reader == nil {
		t.Error("FileToReader() reader = nil")
	}
	os.Remove(tempFile)
}

func TestURLToReader(t *testing.T) {
	var (
		reader *utils.ImageReader
		err    error
	)
	reader, err = utils.URLToReader("", fmt.Sprintf("%d", time.Now().Unix()), 0)
	if err == nil {
		t.Error("URLToReader() error = nil")
	}
	if reader != nil {
		t.Error("URLToReader() reader = ", reader)
	}

	reader, err = utils.URLToReader("https://www.baidu.com", fmt.Sprintf("%d", time.Now().Unix()), time.Second*10)
	if err == nil {
		t.Error("URLToReader() error = nil")
	}
	if reader != nil {
		t.Error("URLToReader() reader = ", reader)
	}

	reader, err = utils.URLToReader("data/test.jpg", fmt.Sprintf("%d", time.Now().Unix()), time.Second*10)
	if err == nil {
		t.Error("URLToReader() error = nil")
	}
	if reader != nil {
		t.Error("URLToReader() reader = ", reader)
	}

	reader, err = utils.URLToReader("https://www.gstatic.com/webp/gallery/1.webp", fmt.Sprintf("%d", time.Now().Unix()), time.Second*10)
	if err != nil {
		t.Error("URLToReader() error = ", err)
	}
	if reader == nil {
		t.Error("URLToReader() reader = nil")
	}
	t.Logf("reader.Name() = %s", reader.Name())

	reader, err = utils.URLToReader("https://dummyimage.com/600x400/000/fff.png", fmt.Sprintf("%d", time.Now().Unix()), time.Second*10)
	if err != nil {
		t.Error("URLToReader() error = ", err)
	}
	if reader == nil {
		t.Error("URLToReader() reader = nil")
	}
	t.Logf("reader.Name() = %s", reader.Name())

	reader, err = utils.URLToReader("https://httpbin.org/image/jpeg", fmt.Sprintf("%d", time.Now().Unix()), time.Second*10)
	if err != nil {
		t.Error("URLToReader() error = ", err)
	}
	if reader == nil {
		t.Error("URLToReader() reader = nil")
	}
	t.Logf("reader.Name() = %s", reader.Name())
}

func TestBase64ToReader(t *testing.T) {
	var (
		reader *utils.ImageReader
		err    error
	)
	reader, err = utils.Base64ToReader("", fmt.Sprintf("%d", time.Now().Unix()))
	if err == nil {
		t.Error("Base64ToReader() error = nil")
	}
	if reader != nil {
		t.Error("Base64ToReader() reader = ", reader)
	}

	reader, err = utils.Base64ToReader("nonexistent", fmt.Sprintf("%d", time.Now().Unix()))
	if err == nil {
		t.Error("Base64ToReader() error = nil")
	}
	if reader != nil {
		t.Error("Base64ToReader() reader = ", reader)
	}

	reader, err = utils.Base64ToReader("/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBwgHBgkIBwgKCgkLDRYPDQwMDRsUFRAWIB0iIiAdHx8kKDQsJCYxJx8fLT0tMTU3Ojo6Iys/BQYH/8QAtRABAAIB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5Pj/3", fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		t.Error("Base64ToReader() error = ", err)
	}
	if reader == nil {
		t.Error("Base64ToReader() reader = nil")
	}
	t.Logf("reader.Name() = %s", reader.Name())

	reader, err = utils.Base64ToReader("data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBwgHBgkIBwgKCgkLDRYPDQwMDRsUFRAWIB0iIiAdHx8kKDQsJCYxJx8fLT0tMTU3Ojo6Iys/BQYH/8QAtRABAAIB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5Pj/3", fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		t.Error("Base64ToReader() error = ", err)
	}
	if reader == nil {
		t.Error("Base64ToReader() reader = nil")
	}
	t.Logf("reader.Name() = %s", reader.Name())
}

func TestFileToBase64(t *testing.T) {
	var (
		b64 string
		err error
	)
	b64, err = utils.FileToBase64("")
	if err == nil {
		t.Error("FileToBase64() error = nil")
	}
	if b64 != "" {
		t.Error("FileToBase64() b64 = ", b64)
	}

	b64, err = utils.FileToBase64("nonexistent.png")
	if err == nil {
		t.Error("FileToBase64() error = nil")
	}
	if b64 != "" {
		t.Error("FileToBase64() b64 = ", b64)
	}

	b64, err = utils.FileToBase64("image.go")
	if err == nil {
		t.Error("FileToBase64() error = nil")
	}
	if b64 != "" {
		t.Error("FileToBase64() b64 = ", b64)
	}

	var tempFile string
	tempFile, err = createTempImage("png")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	b64, err = utils.FileToBase64(tempFile)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("FileToBase64() b64 = ''")
	}
	os.Remove(tempFile)

	tempFile, err = createTempImage("jpg")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	b64, err = utils.FileToBase64(tempFile)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("FileToBase64() b64 = ''")
	}
	os.Remove(tempFile)

	tempFile, err = createTempImage("jpeg")
	if err != nil {
		t.Error("createTempImage() error = ", err)
	}
	b64, err = utils.FileToBase64(tempFile)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("FileToBase64() b64 = ''")
	}
	os.Remove(tempFile)
}

func TestURLToBase64(t *testing.T) {
	var (
		b64 string
		err error
	)
	b64, err = utils.URLToBase64("https://www.gstatic.com/webp/gallery/1.webp", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://www.gstatic.com/webp/gallery/2.webp", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://www.gstatic.com/webp/gallery/3.webp", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://www.gstatic.com/webp/gallery/4.webp", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://www.gstatic.com/webp/gallery/5.webp", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://upload.wikimedia.org/wikipedia/commons/6/6a/PNG_Test.png", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://httpbin.org/image/png", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://avatars.githubusercontent.com/u/9919?v=4", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://dummyimage.com/600x400/000/fff.png", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://httpbin.org/image/jpeg", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://dummyimage.com/800x600/ff0000/ffffff.jpg", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://dummyimage.com/400x300/00ff00/000000.jpg", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://dummyimage.com/500x400/0000ff/ffffff.jpg", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])

	b64, err = utils.URLToBase64("https://dummyimage.com/300x200/ffff00/000000.jpg", time.Second*10)
	if err != nil {
		t.Error("FileToBase64() error = ", err)
	}
	if b64 == "" {
		t.Error("URLToBase64() b64 = ''")
	}
	t.Logf("b64 = %s", b64[:100])
}
