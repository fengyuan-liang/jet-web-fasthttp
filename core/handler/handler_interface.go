package handler

import (
	"github.com/valyala/fasthttp"
)

type IHandler interface {
	ServeHTTP(ctx *fasthttp.RequestCtx, args []string)
}
