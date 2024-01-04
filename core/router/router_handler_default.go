// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/handler"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/inject"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"reflect"
	"strings"
)

var DefaultJetRouter = NewJetRouter("0")

func ServeHTTP(ctx *fasthttp.RequestCtx) {
	DefaultJetRouter.ServeHTTP(ctx)
}

func Register(rcvrs ...interface{}) {
	// provide by cmd
	for _, rcvr := range rcvrs {
		register(rcvr)
	}
}

func RegisterByInject() {
	// provide by inject
	inject.Invoke(func(controllerList inject.JetControllerList) {
		for _, controller := range controllerList.Handlers {
			register(controller)
		}
	})
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
	hooks := new(hook.Hook).GenHook(&val)
	// Install the methods
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if strings.Contains(method.Name, "Hook") {
			continue
		}
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
