/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 19:36:52
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-02 04:27:24
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpClient_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/liusuxian/go-aisdk/httpClient"
	"net/http"
	"reflect"
	"testing"
)

var errTestMarshallerFailed = errors.New("test marshaller failed")

type failingMarshaller struct{}

func (fm *failingMarshaller) Marshal(_ any) ([]byte, error) {
	return []byte{}, errTestMarshallerFailed
}

func TestRequestBuilderReturnsMarshallerErrors(t *testing.T) {
	builder := httpClient.NewRequestBuilder(httpClient.WithMarshaller(&failingMarshaller{}))
	if _, err := builder.Build(context.Background(), "", "", struct{}{}, nil); !errors.Is(err, errTestMarshallerFailed) {
		t.Fatalf("did not return error when marshaller failed: %v", err)
	}
}

func TestRequestBuilderReturnsRequest(t *testing.T) {
	var (
		ctx         = context.Background()
		method      = http.MethodPost
		url         = "/foo"
		request     = map[string]string{"foo": "bar"}
		builder     = httpClient.NewRequestBuilder()
		reqBytes, _ = (&httpClient.JSONMarshaller{}).Marshal(request)
		want, _     = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBytes))
	)

	got, _ := builder.Build(ctx, method, url, request, nil)
	if !reflect.DeepEqual(got.Body, want.Body) ||
		!reflect.DeepEqual(got.URL, want.URL) ||
		!reflect.DeepEqual(got.Method, want.Method) {
		t.Errorf("Build() got = %v, want = %v", got, want)
	}
}

func TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil(t *testing.T) {
	var (
		ctx     = context.Background()
		method  = http.MethodGet
		url     = "/foo"
		want, _ = http.NewRequestWithContext(ctx, method, url, nil)
		builder = httpClient.NewRequestBuilder()
	)
	got, _ := builder.Build(ctx, method, url, nil, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() got = %v, want = %v", got, want)
	}
}
