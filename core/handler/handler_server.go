// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/valyala/fasthttp"
	"reflect"
)

type HandlerFunc = func(ctx *fasthttp.RequestCtx) error

type handler struct {
	rcvr   *reflect.Value
	method *reflect.Value
}

func (h *handler) ServeHTTP(ctx *fasthttp.RequestCtx, args []string) {

}
