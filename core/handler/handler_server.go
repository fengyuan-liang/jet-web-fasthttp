// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/valyala/fasthttp"
	"jet-web/pkg/constant"
	"jet-web/pkg/xlog"
	"reflect"
)

type HandlerFunc = func(ctx *fasthttp.RequestCtx) error

type handler struct {
	rcvr             *reflect.Value
	method           *reflect.Method
	ctxType          contextType
	parametersType   parametersType
	returnValuesType returnValuesType
}

var handlerLog = xlog.NewWith("handler_log")

func (h *handler) ServeHTTP(ctx *fasthttp.RequestCtx, args []string) {
	switch string(ctx.Method()) {
	case constant.MethodGet:
		h.handleGetRequest(ctx, args)
	case constant.MethodPost, constant.MethodPut:
	case constant.MethodHead:

	}
}

func (h *handler) handleGetRequest(ctx *fasthttp.RequestCtx, args []string) {
	var (
		uri = ctx.URI().String()
	)
	handlerLog.Infof("handle uri[%s]", uri)
}
