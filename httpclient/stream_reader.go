/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 18:00:38
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-19 16:46:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`data: {"error":`)
)

// StreamDataHandler 流式数据处理函数
type StreamDataHandler[T Streamable] func(response T, isFinished bool) (err error)

// Streamable 可流式传输的类型
type Streamable any

// StreamReader 流读取器
type StreamReader[T Streamable] struct {
	emptyMessagesLimit          uint
	isFinished                  bool
	reader                      *bufio.Reader
	response                    *http.Response
	streamReturnIntervalTimeout time.Duration
	streamReturnIntervalTimer   *time.Timer
	errAccumulator              ErrorAccumulator
	unmarshaler                 Unmarshaler
	// 统计字段
	startTime  time.Time
	chunkCount int
	// 响应头
	HttpHeader
}

// StreamStatsReceiver 流式传输统计信息接收器
type StreamStatsReceiver interface {
	SetStreamStats(stats StreamStats) // 设置流式传输统计信息
}

// StreamStats 流式传输统计信息
type StreamStats struct {
	TotalDurationMs int64     `json:"total_duration_ms"` // 传输总耗时（持续更新）
	DurationMs      int64     `json:"duration_ms"`       // 单次传输耗时
	ChunkCount      int       `json:"chunk_count"`       // 传输的 chunk 数量（持续更新）
	StartTime       time.Time `json:"start_time"`        // 传输开始时间
	EndTime         time.Time `json:"end_time"`          // 传输结束时间（持续更新）
}

// ForEach 循环处理流式数据，对每个数据项调用处理函数
func (stream *StreamReader[T]) ForEach(handler StreamDataHandler[T]) (err error) {
	if stream.streamReturnIntervalTimer == nil {
		stream.streamReturnIntervalTimer = time.NewTimer(stream.streamReturnIntervalTimeout)
	}
	// 在单独的 goroutine 中处理流
	var (
		lineChan = make(chan T, 1)
		errChan  = make(chan error, 1)
		done     = make(chan struct{})
	)
	defer stream.Close()
	defer close(done)
	defer stream.streamReturnIntervalTimer.Stop()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				resp, finished, e := stream.Recv()
				if e != nil {
					errChan <- e
					return
				}
				if finished {
					errChan <- nil
					return
				}
				lineChan <- resp
			}
		}
	}()

	for {
		select {
		case <-stream.streamReturnIntervalTimer.C:
			return ErrStreamReturnIntervalTimeout
		case err = <-errChan:
			if err == nil {
				var empty T
				return handler(empty, true)
			}
			return
		case line := <-lineChan:
			if err = handler(line, false); err != nil {
				return
			}
			stream.streamReturnIntervalTimer.Reset(stream.streamReturnIntervalTimeout)
		}
	}
}

// Recv 接收数据
func (stream *StreamReader[T]) Recv() (response T, isFinished bool, err error) {
	var (
		processingStartTime = time.Now()
		rawLine             []byte
	)
	if rawLine, err = stream.processLines(); err != nil {
		if stream.isFinished {
			isFinished = true
			err = nil
		}
		return
	}
	// 解析数据
	if err = stream.unmarshaler.Unmarshal(rawLine, &response); err != nil {
		return
	}
	// 更新统计信息
	stream.chunkCount++
	if statsReceiver, ok := Streamable(&response).(StreamStatsReceiver); ok {
		now := time.Now()
		stats := StreamStats{
			TotalDurationMs: now.Sub(stream.startTime).Milliseconds(),
			DurationMs:      now.Sub(processingStartTime).Milliseconds(),
			ChunkCount:      stream.chunkCount,
			StartTime:       stream.startTime,
			EndTime:         now,
		}
		statsReceiver.SetStreamStats(stats)
	}
	return
}

// RecvRaw 接收原始数据
func (stream *StreamReader[T]) RecvRaw() (b []byte, err error) {
	return stream.processLines()
}

// processLines 处理行数据
func (stream *StreamReader[T]) processLines() (b []byte, err error) {
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
func (stream *StreamReader[T]) unmarshalError() (errResp map[string]any) {
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
func (stream *StreamReader[T]) Close() (err error) {
	return stream.response.Body.Close()
}
