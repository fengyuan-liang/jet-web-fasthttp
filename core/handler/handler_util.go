// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/valyala/fasthttp"
	"jet-web/pkg/constant"
	"jet-web/pkg/utils"
	"reflect"
	"strings"
	"syscall"
)

func parseReqWithBody(ret *reflect.Value, ctx *fasthttp.RequestCtx) error {
	var (
		retElem = ret.Elem()
		req     = &ctx.Request
	)
	if cmdArgsBytes := req.Header.Peek("*"); cmdArgsBytes != nil {
		if field := retElem.FieldByName("CmdArgs"); field.IsValid() {
			field.Set(reflect.ValueOf([]string{string(cmdArgsBytes)}))
		}
	}
	if isJsonCall(req) {
		if req.Header.ContentLength() == 0 {
			return nil
		}
		return utils.Decode(req.BodyStream(), ret.Interface())
	}
	return syscall.EINVAL
}

func isJsonCall(req *fasthttp.Request) bool {
	var ct string

	if ctBytes := req.Header.Peek(constant.HeaderContentType); ctBytes == nil {
		return false
	} else {
		ct = string(ctBytes)
	}

	return ct == "application/json" || strings.HasPrefix(ct, "application/json;")
}
