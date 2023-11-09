// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"reflect"
	"syscall"
)

// ---------------------------------------------------------------------------

// CreatorFunc define func of HandlerCreator
type CreatorFunc = func(rcvr *reflect.Value, method *reflect.Method) (IHandler, error)

type HandlerCreator struct {
	MethodPrefix string
	Creator      CreatorFunc
}

var unusedError *error
var unusedCtx *context.Ctx
var typeOfError = reflect.TypeOf(unusedError).Elem()
var typeOfCtx = reflect.TypeOf(unusedCtx).Elem()

type contextType int
type parametersType int
type returnValuesType int

const (
	noCtx contextType = iota
	defaultCtx
	customCtx
)

const (
	noParameter                  parametersType = iota
	oneParameterAndFirstIsCtx                   // only one parameter and parameter are ctx
	oneParameterAndFirstNotIsCtx                // only one parameter and parameter not are ctx
	twoParameterAndFirstIsCtx                   // two parameters and the first parameter is ctx
	twoParameterAndSecondIsCtx                  // two parameters and the second parameter are ctx
)

const (
	noReturnValue                  returnValuesType = iota
	OneReturnValueAndIsError                        // only return value and is error
	OneReturnValueAndNotError                       // only return value and not an error
	twoReturnValueAndFirstIsError                   // two return value and the first is error
	twoReturnValueAndSecondIsError                  // two return value and the second is error
)

var handlerCreatorLog = xlog.NewWith("handle_create_log")

// New common creator
// Method spec:
//
//	(rcvr *XXXX) YYYY(ctx jet.Ctx, req ZZZZ) (err error)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx, req ZZZZ) (ret RRRR, err error)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx, req ZZZZ)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx, req ZZZZ) (ret RRRR, err error)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx, req ZZZZ) (err error)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx) (err error)
//	(rcvr *XXXX) YYYY(ctx jet.Ctx) (ret RRRR, err error)
//	(rcvr *XXXX) YYYY() (err error)
//	(rcvr *XXXX) YYYY() (ret RRRR, err error)
func (p HandlerCreator) New(rcvr *reflect.Value, method *reflect.Method) (IHandler, error) {
	var (
		mtype           = method.Type
		methodNumIn     = mtype.NumIn() - 1
		methodNumOut    = mtype.NumOut()
		methodName      = method.Name
		ctxType         = noCtx
		parameterType   = noParameter
		returnValueType = noReturnValue
	)
	// only allow up to two parameters
	if methodNumIn > 2 {
		handlerCreatorLog.Errorf("too many parameters, Method [%s] exceeds the default limit", methodName)
		return nil, syscall.EINVAL
	}
	// only allow up to two return values allowed
	if methodNumOut > 2 {
		handlerCreatorLog.Errorf("too many return values, Method [%s] exceeds the default limit", methodName)
		return nil, syscall.EINVAL
	}
	switch methodNumIn {
	case 1:
		in := mtype.In(1)
		if in.Kind() == reflect.Ptr {
			if in.Implements(typeOfCtx) {
				ctxType++
				parameterType = oneParameterAndFirstIsCtx
			} else {
				parameterType = oneParameterAndFirstNotIsCtx
			}
		} else if in.Kind() == reflect.Struct || in.Kind() == reflect.Interface {
			if in == typeOfCtx || in.Implements(typeOfCtx) {
				parameterType = oneParameterAndFirstIsCtx
			}
		} else {
			parameterType = oneParameterAndFirstNotIsCtx
		}
	case 2:
		firstIn, secondIn := mtype.In(1), mtype.In(2)
		if firstIn.Kind() == reflect.Ptr {
			firstIn = firstIn.Elem()
		}
		if secondIn.Kind() == reflect.Ptr {
			secondIn = secondIn.Elem()
		}
		if firstIn.Implements(typeOfCtx) {
			ctxType++
			parameterType = twoParameterAndFirstIsCtx
		} else if secondIn.Implements(typeOfCtx) {
			ctxType++
			parameterType = twoParameterAndSecondIsCtx
		}
	}
	// NumOut
	switch methodNumOut {
	case 1:
		out := mtype.Out(0)
		if out.Kind() == reflect.Ptr {
			if out.Implements(typeOfError) {
				returnValueType = OneReturnValueAndIsError
			} else {
				returnValueType = OneReturnValueAndNotError
			}
		} else if out == typeOfError {
			returnValueType = OneReturnValueAndIsError
		} else {
			returnValueType = OneReturnValueAndNotError
		}
	case 2:
		firstOut, secondOut := mtype.Out(0), mtype.Out(1)
		if firstOut.Kind() == reflect.Ptr {
			if firstOut.Implements(typeOfError) {
				returnValueType = twoReturnValueAndFirstIsError
			}
		} else if firstOut == typeOfError {
			returnValueType = twoReturnValueAndFirstIsError
		}
		if secondOut.Kind() == reflect.Ptr {
			if firstOut.Implements(typeOfError) {
				returnValueType = twoReturnValueAndSecondIsError
			}
		} else if secondOut == typeOfError {
			returnValueType = twoReturnValueAndSecondIsError
		}
	}
	return &handler{
		rcvr:             rcvr,
		method:           method,
		ctxType:          ctxType,
		parametersType:   parameterType,
		returnValuesType: returnValueType,
		hook:             new(hook.Hook),
	}, nil
}
