package handler

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
	"github.com/valyala/fasthttp"
)

type IHandler interface {
	ServeHTTP(ctx *fasthttp.RequestCtx, args []string)
	AddHook(hooks *hook.Hook)
}
