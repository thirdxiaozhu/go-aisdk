/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 22:28:26
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-09 10:57:52
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"bufio"
	"bytes"
	"fmt"
	utils "github.com/liusuxian/go-openai/internal"
	"io"
	"net/http"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`data: {"error":`)
)

// streamable 可流式传输的类型
type streamable interface {
}

// streamReader 流读取器
type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool
	reader             *bufio.Reader
	response           *http.Response
	errAccumulator     utils.ErrorAccumulator
	unmarshaler        utils.Unmarshaler
	httpHeader
}

// Recv 接收数据
func (stream *streamReader[T]) Recv() (response T, err error) {
	var rawLine []byte
	if rawLine, err = stream.RecvRaw(); err != nil {
		return
	}

	if err = stream.unmarshaler.Unmarshal(rawLine, &response); err != nil {
		return
	}
	return
}

// RecvRaw 接收数据
func (stream *streamReader[T]) RecvRaw() (b []byte, err error) {
	if stream.isFinished {
		return nil, io.EOF
	}

	return stream.processLines()
}

// processLines 处理行数据
func (stream *streamReader[T]) processLines() (b []byte, err error) {
	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			if respErr := stream.unmarshalError(); respErr != nil {
				return nil, fmt.Errorf("error, %v", respErr)
			}
			return nil, readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			if writeErr := stream.errAccumulator.Write(noSpaceLine); writeErr != nil {
				return nil, writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return nil, ErrTooManyEmptyStreamMessages
			}
			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return nil, io.EOF
		}

		return noPrefixLine, nil
	}
}

// unmarshalError 解析错误响应数据
func (stream *streamReader[T]) unmarshalError() (errResp map[string]any) {
	var errBytes []byte
	if errBytes = stream.errAccumulator.Bytes(); len(errBytes) == 0 {
		return
	}

	if err := stream.unmarshaler.Unmarshal(errBytes, &errResp); err != nil {
		errResp = nil
		return
	}
	return
}

// Close 关闭流
func (stream *streamReader[T]) Close() (err error) {
	return stream.response.Body.Close()
}
