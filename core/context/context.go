// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package context

import (
	"context"
	"github.com/valyala/fasthttp"
	"jet-web/pkg/xlog"
)

// Ctx represents the Context which holds the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers, and so on.
type Ctx struct {
	context.Context

	log  *xlog.Logger // log for context
	req  *fasthttp.Request
	resp *fasthttp.Response
}
