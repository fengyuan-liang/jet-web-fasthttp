// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/handler"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"reflect"
)

var DefaultJetRouter = NewJetRouter("0")

func ServeHTTP(ctx *fasthttp.RequestCtx) {
	DefaultJetRouter.ServeHTTP(ctx)
}

func Register(rcvrs ...interface{}) {
	for _, rcvr := range rcvrs {
		register(rcvr)
	}
}

func register(rcvr interface{}) {
	var (
		typ = reflect.TypeOf(rcvr)
		val = reflect.ValueOf(rcvr)
	)
	if typ.Kind() != reflect.Ptr {
		xlog.Infof("receiver [%s] not pointer, Jet recommends passing a pointer as the receiver.", typ.Name())
		typ = reflect.PtrTo(typ)
		val = reflect.ValueOf(typ)
	}
	// global hook
	hooks := hook.GenHook(&val)
	// Install the methods
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		_, h, err := handler.Factory.Create(&val, &method)
		if err != nil {
			xlog.Errorf("handler.Factory.Create error:%v", err)
			continue
		}
		h.AddHook(hooks)
		DefaultJetRouter.RegisterRouter(method.Name, h)
		xlog.Debug("Install", "=>", method.Name)
	}
}
