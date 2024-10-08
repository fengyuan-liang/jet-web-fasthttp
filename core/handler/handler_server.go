// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"reflect"
)

type handler struct {
	rcvr             *reflect.Value
	method           *reflect.Method
	ctxType          contextType
	parametersType   parametersType
	returnValuesType returnValuesType
	hook             *hook.Hook
}

var handlerLog = xlog.NewWith("handler_log")

func (h handler) ServeHTTP(ctx *fasthttp.RequestCtx, args []string) {
	switch string(ctx.Method()) {
	case constant.MethodGet, constant.MethodPost, constant.MethodPut, constant.MethodDelete:
		h.handleRequest(ctx, args)
	case constant.MethodHead:

	}
}

func (h handler) AddHook(hooks *hook.Hook) {
	h.hook.PostParamsParseHooks = append(h.hook.PostParamsParseHooks, hooks.PostParamsParseHooks...)
	h.hook.PostMethodExecuteHooks = append(h.hook.PostMethodExecuteHooks, hooks.PostMethodExecuteHooks...)
	h.hook.PreMethodExecuteHooks = append(h.hook.PreMethodExecuteHooks, hooks.PreMethodExecuteHooks...)
}

func (h handler) handleRequest(ctx *fasthttp.RequestCtx, args []string) {
	var (
		uri         = ctx.URI().String()
		mtype       = h.method.Type
		methodArgs  = []reflect.Value{*h.rcvr}
		param       reflect.Value
		err         error
		jetCtx      = context.NewContext(ctx)
		jetCtxValue = reflect.ValueOf(jetCtx)
	)
	handlerLog.Debugf("handle uri[%s]", uri)

	// global hook
	if len(hook.PostJetCtxInitHooks) > 0 {
		for _, postJetCtxInitHook := range hook.PostJetCtxInitHooks {
			postJetCtxInitHook(jetCtx)
		}
	}

	// handle PreMethodExecuteHook
	if h.hook.HasPreMethodExecuteHooks() {
		if err = h.hook.PreMethodExecuteHook(jetCtxValue); err != nil {
			FailHandler(ctx, err.Error())
			return
		}
	}

	switch h.parametersType {
	case noParameter:
		// noting to do
	case oneParameterAndFirstIsCtx:
		methodArgs = append(methodArgs, jetCtxValue)
	case oneParameterAndFirstNotIsCtx:
		param, err = h.handleParam(ctx, args, mtype.In(1), err)
		if err != nil {
			handlerLog.Errorf("handler err: %v", err.Error())
			FailHandler(ctx, err.Error())
			return
		}
		// handle postParamsParseHook
		if err = h.hook.PostParamsParse(param); err != nil {
			FailHandler(ctx, err.Error())
			return
		}
		methodArgs = append(methodArgs, param)
	case twoParameterAndFirstIsCtx:
		// handle param
		param, err = h.handleParam(ctx, args, mtype.In(2), err)
		if err != nil {
			handlerLog.Errorf("handler err: %v", err.Error())
			FailHandler(ctx, err.Error())
			return
		}
		// handle postParamsParseHook
		if err = h.hook.PostParamsParse(param); err != nil {
			FailHandler(ctx, err.Error())
			return
		}
		methodArgs = append(methodArgs, jetCtxValue, param)
	case twoParameterAndSecondIsCtx:
		// handle param
		param, err = h.handleParam(ctx, args, mtype.In(1), err)
		if err != nil {
			handlerLog.Errorf("handler err: %v", err.Error())
			FailHandler(ctx, err.Error())
			return
		}
		// handle postParamsParseHook
		if err = h.hook.PostParamsParse(param); err != nil {
			FailHandler(ctx, err.Error())
			return
		}
		methodArgs = append(methodArgs, param, jetCtxValue)
	default:
		panic("illegal method signature")
	}

	callValues := h.method.Func.Call(methodArgs)

	switch h.returnValuesType {
	case noReturnValue:
		// noting to do
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
			SuccessHandler(ctx, fmt.Sprintf("%v", data))
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
			RestSuccessHandler(ctx, data)
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
			// handle postParamsParseHook
			postData := callValues[0].Interface()
			if h.hook.HasPostMethodExecuteHook() {
				if postData, err = h.hook.PostMethodExecuteHook(callValues[0]); err != nil {
					FailHandler(ctx, err.Error())
					return
				}
			}
			SuccessHandler(ctx, fmt.Sprintf("%v", postData))
		} else {
			SuccessHandler(ctx, constant.EmptyString)
		}
	}
	return
}

func (h handler) handleParam(ctx *fasthttp.RequestCtx, args []string, in reflect.Type, err error) (reflect.Value, error) {
	var (
		paramIsPtr bool
	)
	if in.Kind() == reflect.Ptr {
		in = in.Elem()
		paramIsPtr = true
	}
	// the value is ptr
	value := reflect.New(in)
	if err = parseReqDefault(ctx, value, args); err != nil {
		xlog.Errorf("parseReqDefault err: %v", err.Error())
		return reflect.Value{}, err
	}
	if !paramIsPtr {
		value = value.Elem()
	}
	return value, err
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
		if ctx.Request.Header.ContentLength() <= 0 || len(ctx.Request.Body()) <= 0 {
			return parseValue(param, ctx, "form")
		}
		return utils.ByteToObj(ctx.Request.Body(), param.Interface())
	} else if isFormCall(&ctx.Request) {
		return parseForm(param, ctx)
	} else {
		return parseValue(param, ctx, "form")
	}
}
