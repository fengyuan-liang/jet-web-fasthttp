// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package context

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"os"
)

type Args struct {
	CmdArgs    []string
	FormParam1 string `json:"form_param1" form:"form_param1"`
	FormParam2 string `json:"form_param2" form:"form_param1"`
}

type Ctx interface {
	Logger() *xlog.Logger
	Request() *fasthttp.Request
	Response() *fasthttp.Response
	Get(key string) (value any, exists bool)
	MustGet(key string) (value any)
	Put(key string, value any)
	Keys() map[string]any
}

func NewContext(ctx *fasthttp.RequestCtx, logs ...*xlog.Logger) Ctx {
	var log *xlog.Logger
	if logs != nil && len(logs) > 0 {
		log = logs[0]
	} else {
		log = xlog.New(os.Stderr, "", xlog.Ldefault, xlog.GenReqId())
	}
	return &Context{
		keys: maps.NewLinkedHashMap[string, any](),
		log:  log,
		req:  &ctx.Request,
		resp: &ctx.Response,
	}
}

// Context represents the Context which holds the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers, and so on.
type Context struct {
	// keys is a key/value pair exclusively for the context of each request.
	// default maps.LinkedHashMap
	keys maps.IMap[string, any]
	log  *xlog.Logger // log for context
	req  *fasthttp.Request
	resp *fasthttp.Response
}

func (c *Context) Logger() *xlog.Logger {
	return c.log
}

func (c *Context) Request() *fasthttp.Request {
	return c.req
}

func (c *Context) Response() *fasthttp.Response {
	return c.resp
}

func (c *Context) Get(key string) (value any, exists bool) {
	return c.keys.Get(key)
}

func (c *Context) MustGet(key string) (value any) {
	return c.keys.MustGet(key)
}

func (c *Context) Keys() (rawMap map[string]any) {
	return c.keys.RawMap()
}

func (c *Context) Put(key string, value any) {
	if c.keys == nil {
		c.keys = maps.NewLinkedHashMap[string, any]()
	}
	c.keys.Put(key, value)
}
