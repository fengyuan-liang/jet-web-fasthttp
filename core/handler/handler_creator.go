// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"jet-web/core/context"
	"jet-web/pkg/xlog"
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
	oneParameterAndFirstIsCtx                   // only one parameter and parameter is ctx
	oneParameterAndFirstNotIsCtx                // only one parameter and parameter not is ctx
	TwoParameterAndFirstIsCtx                   // two parameter and the first parameter is ctx
	TwoParameterAndSecondIsCtx                  // two parameter and the second parameter is ctx
)

const (
	noReturnValue                  returnValuesType = iota
	OneReturnValueAndIsError                        // only return value and is error
	OneReturnValueAndNotError                       // only return value and not is error
	TwoReturnValueAndFirstIsError                   // two return value and the first  is error
	TwoReturnValueAndSecondIsError                  // two return value and the second is error
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
		methodNumIn     = mtype.NumIn()
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
		in := mtype.In(0)
		if in.Kind() == reflect.Ptr {
			in = in.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "arg type not a pointer:", in.Kind())
			return nil, syscall.EINVAL
		}
		if in.Implements(typeOfCtx) {
			ctxType++
			parameterType = oneParameterAndFirstIsCtx
		} else {
			parameterType = oneParameterAndFirstNotIsCtx
		}
	case 2:
		firstIn, secondIn := mtype.In(0), mtype.In(1)
		if firstIn.Kind() == reflect.Ptr {
			firstIn = firstIn.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "first arg type not a pointer:", firstIn.Kind())
			return nil, syscall.EINVAL
		}
		if secondIn.Kind() == reflect.Ptr {
			secondIn = secondIn.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "second arg type not a pointer:", secondIn.Kind())
			return nil, syscall.EINVAL
		}
		if firstIn.Implements(typeOfCtx) {
			ctxType++
			parameterType = TwoParameterAndFirstIsCtx
		} else if secondIn.Implements(typeOfCtx) {
			ctxType++
			parameterType = TwoParameterAndSecondIsCtx
		}
	}
	// NumOut
	switch methodNumOut {
	case 1:
		out := mtype.Out(0)
		if out.Kind() == reflect.Ptr {
			out = out.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "return value type not a pointer:", out.Kind())
			return nil, syscall.EINVAL
		}
		if out.Implements(typeOfError) {
			returnValueType = OneReturnValueAndIsError
		} else {
			returnValueType = OneReturnValueAndNotError
		}
	case 2:
		firstOut, secondOut := mtype.Out(0), mtype.Out(1)
		if firstOut.Kind() == reflect.Ptr {
			firstOut = firstOut.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "first return value type not a pointer:", firstOut.Kind())
			return nil, syscall.EINVAL
		}
		if secondOut.Kind() == reflect.Ptr {
			secondOut = secondOut.Elem()
		} else {
			handlerCreatorLog.Debug("method", methodName, "second return value type not a pointer:", secondOut.Kind())
			return nil, syscall.EINVAL
		}
		if firstOut.Implements(typeOfError) {
			returnValueType = TwoReturnValueAndFirstIsError
		} else if secondOut.Implements(typeOfError) {
			returnValueType = TwoReturnValueAndSecondIsError
		}
	}
	return &handler{
		rcvr,
		method,
		ctxType,
		parameterType,
		returnValueType,
	}, nil
}

// handlerGetCreator Creator for handling HTTP get requests
var handlerGetCreator = func(rcvr *reflect.Value, method *reflect.Method) (IHandler, error) {

	return &handler{rcvr: rcvr, method: method}, nil
}
