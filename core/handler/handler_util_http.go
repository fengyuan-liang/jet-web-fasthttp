// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
)

const (
	Ok = "Ok"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Response.Header.SetServer("JetServer")
	ctx.SetBodyString("404 Not Found")
}

func SuccessHandler(ctx *fasthttp.RequestCtx, data string) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetServer("JetServer")
	ctx.Response.Header.Set("Content-Type", constant.MIMEApplicationJSON)
	ctx.SetBodyString(data)
}

func RestSuccessHandler(ctx *fasthttp.RequestCtx, data any) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetServer("JetServer")
	ctx.SetContentType(constant.MIMEApplicationJSON)
	ctx.SetBodyString(utils.ObjToJsonStr(data))
}

func FailHandler(ctx *fasthttp.RequestCtx, data string) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.Response.Header.SetServer("JetServer")
	ctx.SetBodyString(data)
}
