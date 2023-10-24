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
	RegisterRouter(path string, handler handler.IHandler)
}

func NewJetRouter(separator string, f ...SplitPathFunc) IJetRouter {
	return &JetRouter{trie: NewRouterTrie[handler.IHandler](separator, f...)}
}

type JetRouter struct {
	trie ITrie[handler.IHandler]
}

func (r *JetRouter) RegisterRouter(path string, handler handler.IHandler) {
	r.trie.Add(path, handler)
}

func (r *JetRouter) ServeHTTP(ctx *fasthttp.RequestCtx) {
	requestURI := convertToFirstLetterUpper(ctx.Method()) + string(ctx.Request.RequestURI())
	if h, args := r.trie.GetAndArgs(requestURI); h != nil {
		h.ServeHTTP(ctx, args)
	} else {
		handler.NotFoundHandler(ctx)
	}
}
