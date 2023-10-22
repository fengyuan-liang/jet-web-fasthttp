// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/valyala/fasthttp"
	"jet-web/pkg/constant"
	"jet-web/pkg/utils"
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

func (h *handler) handleGetRequest(ctx *fasthttp.RequestCtx, args []string) (err error) {
	var (
		uri        = ctx.URI().String()
		mtype      = h.method.Type
		methodArgs = []reflect.Value{*h.rcvr}
	)
	handlerLog.Debugf("handle uri[%s]", uri)
	switch h.parametersType {
	case oneParameterAndFirstIsCtx:

	case oneParameterAndFirstNotIsCtx:
		println(mtype.In(1).Name())
		value := reflect.New(mtype.In(1).Elem())
		if err = parseReqDefault(ctx, value, args); err != nil {
			return
		}
		methodArgs = append(methodArgs, value)
	case twoParameterAndFirstIsCtx:

	case twoParameterAndSecondIsCtx:

	}
	h.method.Func.Call(methodArgs)
	return
}

func setCtx(ctx *fasthttp.RequestCtx) {

}

func parseReqDefault(ctx *fasthttp.RequestCtx, param reflect.Value, args []string) (err error) {
	// query path
	if len(args) > 0 {
		v := param.Elem().FieldByName("CmdArgs")
		if v.IsValid() {
			v.Set(reflect.ValueOf(args))
			if param.Elem().NumField() == 1 {
				return
			}
		}
	}
	if isJsonCall(&ctx.Request) {
		if ctx.Request.Header.ContentLength() == 0 {
			return
		}
		return utils.Decode(ctx.Request.BodyStream(), param.Interface())
	} else {
		return parseValue(param, ctx, "form")
	}
}
