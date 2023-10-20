// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/valyala/fasthttp"
	"jet-web/core/handler"
	"jet-web/pkg/xlog"
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
		typ    = reflect.TypeOf(rcvr)
		val    = reflect.ValueOf(rcvr)
		method reflect.Method
	)
	if typ.Kind() != reflect.Ptr {
		xlog.Infof("receiver [%s] not pointer, Jet recommends passing a pointer as the receiver.", typ.Name())
		typ = reflect.PtrTo(typ)
		val = reflect.ValueOf(typ)
	}
	// Install the methods
	for i := 0; i < typ.NumMethod(); i++ {
		method = typ.Method(i)
		_, h, err := handler.Factory.Create(&val, &method)
		if err != nil {
			continue
		}
		DefaultJetRouter.RegisterRouter(method.Name, h)
		xlog.Debug("Install", "=>", method.Name)
	}
}
