// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/handler"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
	"runtime/debug"
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

func AddMiddleware(jetMiddlewareList ...JetMiddleware) {
	middlewares = append(middlewares, jetMiddlewareList...)
}

func TraceJetMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		defer utils.TraceHttpReqByCtx(ctx)()
		next.ServeHTTP(ctx)
	}), nil
}

func RecoverJetMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				handler.FailServerInternalErrorHandler(ctx, "Internal Server Error")
				utils.PrintPanicInfo("Your server has experienced a panic, please check the stack log below")
				fmt.Printf("Panic: %v\n", err)
				fmt.Printf("stack info\n")
				debug.PrintStack()
			}
		}()
		next.ServeHTTP(ctx)
	}), nil
}
