// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
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

func (h handler) ServeHTTP(ctx *fasthttp.RequestCtx, args []string) {
	switch string(ctx.Method()) {
	case constant.MethodGet:
		h.handleGetRequest(ctx, args)
	case constant.MethodPost, constant.MethodPut:
	case constant.MethodHead:

	}
}

func (h handler) handleGetRequest(ctx *fasthttp.RequestCtx, args []string) {
	var (
		uri        = ctx.URI().String()
		mtype      = h.method.Type
		methodArgs = []reflect.Value{*h.rcvr}
		err        error
	)
	handlerLog.Debugf("handle uri[%s]", uri)
	switch h.parametersType {
	case oneParameterAndFirstIsCtx:

	case oneParameterAndFirstNotIsCtx:
		var (
			in         = mtype.In(1)
			paramIsPtr bool
		)
		if in.Kind() == reflect.Ptr {
			in = in.Elem()
		} else {
			paramIsPtr = false
		}
		// the value is ptr
		value := reflect.New(in)
		if err = parseReqDefault(ctx, value, args); err != nil {
			return
		}
		if !paramIsPtr {
			value = value.Elem()
		}
		methodArgs = append(methodArgs, value)
	case twoParameterAndFirstIsCtx:

	case twoParameterAndSecondIsCtx:

	}
	callValues := h.method.Func.Call(methodArgs)
	switch h.returnValuesType {
	case noReturnValue:
		// noting to do
		SuccessHandler(ctx, constant.EmptyString)
	case OneReturnValueAndIsError:
		if callValues[0].Interface() != nil {
			err = callValues[0].Interface().(error)
			FailHandler(ctx, err.Error())
		} else {
			SuccessHandler(ctx, constant.EmptyString)
		}
	case OneReturnValueAndNotError:
		if callValues[0].Interface() != nil {
			data := callValues[0].Interface()
			SuccessHandler(ctx, utils.ObjToJsonStr(data))
		} else {
			SuccessHandler(ctx, constant.EmptyString)
		}
	case twoReturnValueAndFirstIsError:
		if callValues[0].Interface() != nil {
			err = callValues[0].Interface().(error)
			FailHandler(ctx, err.Error())
			return
		}
		if callValues[1].Interface() != nil {
			data := callValues[1].Interface()
			SuccessHandler(ctx, utils.ObjToJsonStr(data))
		} else {
			SuccessHandler(ctx, constant.EmptyString)
		}
	case twoReturnValueAndSecondIsError:
		if callValues[1].Interface() != nil {
			err = callValues[1].Interface().(error)
			FailHandler(ctx, err.Error())
			return
		}
		if callValues[0].Interface() != nil {
			data := callValues[0].Interface()
			SuccessHandler(ctx, utils.ObjToJsonStr(data))
		} else {
			SuccessHandler(ctx, constant.EmptyString)
		}
	}
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
