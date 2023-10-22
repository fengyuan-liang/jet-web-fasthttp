// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package context

import (
	"context"
	"github.com/valyala/fasthttp"
	"jet-web/pkg/xlog"
)

type Ctx interface {
	Logger() *xlog.Logger
}

// Context represents the Context which holds the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers, and so on.
type Context struct {
	context.Context

	log  *xlog.Logger // log for context
	req  *fasthttp.Request
	resp *fasthttp.Response
}

func (c *Context) Logger() *xlog.Logger {
	return c.log
}

type Args struct {
	CmdArgs    []string
	FormParam1 string `json:"form_param1" form:"form_param1"`
	FormParam2 string `json:"form_param2" form:"form_param1"`
}
