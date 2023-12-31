// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import "github.com/fengyuan-liang/jet-web-fasthttp/core/context"

// Ctx is the most important part of Jet. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Ctx interface {
	context.Ctx
}

type Args struct {
	CmdArgs    []string
	FormParam1 string `json:"form_param1" form:"form_param1"`
	FormParam2 string `json:"form_param2" form:"form_param1"`
}
