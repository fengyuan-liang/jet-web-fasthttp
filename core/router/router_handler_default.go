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

var DefaultJetRouter = NewJetRouter("0", splitCamelCaseFunc)

func ServeHTTP(ctx *fasthttp.RequestCtx) {
	DefaultJetRouter.ServeHTTP(ctx)
}

func Register(rcvr interface{}) {
	var (
		typ    = reflect.TypeOf(rcvr)
		val    = reflect.ValueOf(rcvr)
		method reflect.Method
	)
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
