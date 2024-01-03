// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import "github.com/valyala/fasthttp"

type JetHandlerFunc func(ctx *fasthttp.RequestCtx) error

func (f JetHandlerFunc) ServeHTTP(ctx *fasthttp.RequestCtx) error {
	return f(ctx)
}

var JetMiddlewareList = make([]JetHandlerFunc, 0)

func AddJetMiddleware(jetMiddleware JetHandlerFunc) {
	JetMiddlewareList = append(JetMiddlewareList, jetMiddleware)
}
