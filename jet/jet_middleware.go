// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/handler"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
	"time"
)

type JetHandlerFunc func(ctx *fasthttp.RequestCtx)

func (f JetHandlerFunc) ServeHTTP(ctx *fasthttp.RequestCtx) {
	f(ctx)
}

func (f JetHandlerFunc) RegisterRouter(path string, handler handler.IHandler) {
	// noting to do
}

type JetMiddleware func(next router.IJetRouter) (router.IJetRouter, error)

var middlewares []JetMiddleware

func AddMiddleware(jetMiddleware JetMiddleware) {
	middlewares = append(middlewares, jetMiddleware)
}

func TraceJetMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		defer utils.TraceHttpReq(ctx, time.Now())
		next.ServeHTTP(ctx)
	}), nil
}
