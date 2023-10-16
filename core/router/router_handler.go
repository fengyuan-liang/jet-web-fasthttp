// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/valyala/fasthttp"
	"jet-web/core/handler"
)

type IJetRouter interface {
	ServeHTTP(ctx *fasthttp.RequestCtx)
}

func NewJetRouter(separator string, f ...SplitPathFunc) IJetRouter {
	return &JetRouter{trie: NewRouterTrie[handler.IHandler](separator, f...)}
}

type JetRouter struct {
	trie ITrie[handler.IHandler]
}

func (r *JetRouter) ServeHTTP(ctx *fasthttp.RequestCtx) {
	requestURI := ctx.Request.RequestURI()
	if handler, args := r.trie.GetAndArgs(string(requestURI)); handler != nil {
		handler.ServeHTTP(ctx, args)
	} else {
		notFoundHandler(ctx)
	}
}
