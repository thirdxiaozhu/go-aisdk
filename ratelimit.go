/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 23:40:48
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-08 12:15:34
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"net/http"
	"strconv"
	"time"
)

// RateLimitHeaders Openai 速率限制头信息
type RateLimitHeaders struct {
	LimitRequests     int       `json:"x-ratelimit-limit-requests"`
	LimitTokens       int       `json:"x-ratelimit-limit-tokens"`
	RemainingRequests int       `json:"x-ratelimit-remaining-requests"`
	RemainingTokens   int       `json:"x-ratelimit-remaining-tokens"`
	ResetRequests     ResetTime `json:"x-ratelimit-reset-requests"`
	ResetTokens       ResetTime `json:"x-ratelimit-reset-tokens"`
}

type ResetTime string

func (r ResetTime) String() (s string) {
	return string(r)
}

func (r ResetTime) Time() (t time.Time) {
	d, _ := time.ParseDuration(string(r))
	return time.Now().Add(d)
}

// newRateLimitHeaders 创建速率限制头信息
func newRateLimitHeaders(h http.Header) (rateLimit RateLimitHeaders) {
	limitReq, _ := strconv.Atoi(h.Get("x-ratelimit-limit-requests"))
	limitTokens, _ := strconv.Atoi(h.Get("x-ratelimit-limit-tokens"))
	remainingReq, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-requests"))
	remainingTokens, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-tokens"))
	return RateLimitHeaders{
		LimitRequests:     limitReq,
		LimitTokens:       limitTokens,
		RemainingRequests: remainingReq,
		RemainingTokens:   remainingTokens,
		ResetRequests:     ResetTime(h.Get("x-ratelimit-reset-requests")),
		ResetTokens:       ResetTime(h.Get("x-ratelimit-reset-tokens")),
	}
}
