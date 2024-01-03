// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
)

var TraceJetMiddleware = func(ctx *fasthttp.RequestCtx) error {
	defer utils.TraceHttpReq(ctx)
	return nil
}
