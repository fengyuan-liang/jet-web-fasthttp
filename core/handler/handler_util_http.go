// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/valyala/fasthttp"
	"jet-web/pkg/utils"
)

const (
	Ok = "Ok"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBodyString("404 Not Found")
}

func SuccessHandler(ctx *fasthttp.RequestCtx, data string) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(data)
}

func FailHandler(ctx *fasthttp.RequestCtx, data any) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(utils.ObjToJsonStr(data))
}
