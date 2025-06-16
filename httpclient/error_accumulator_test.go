/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:57:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-03 11:41:57
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient_test

import (
	"bytes"
	"errors"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/httpclient/test"
	"testing"
)

func TestErrorAccumulatorBytes(t *testing.T) {
	accumulator := &httpclient.DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}

	if errBytes := accumulator.Bytes(); len(errBytes) != 0 {
		t.Fatalf("did not return nil with empty bytes: %s", string(errBytes))
	}

	if err := accumulator.Write([]byte("{}")); err != nil {
		t.Fatalf("%+v", err)
	}

	if errBytes := accumulator.Bytes(); len(errBytes) == 0 {
		t.Fatalf("did not return sdkerror bytes when has sdkerror: %s", string(errBytes))
	}
}

func TestErrorByteWriteErrors(t *testing.T) {
	accumulator := &httpclient.DefaultErrorAccumulator{
		Buffer: &test.FailingErrorBuffer{},
	}

	if err := accumulator.Write([]byte("{")); !errors.Is(err, test.ErrTestErrorAccumulatorWriteFailed) {
		t.Fatalf("did not return sdkerror when write failed: %v", err)
	}
}
