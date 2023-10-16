// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/valyala/fasthttp"
	"unicode"
)

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBodyString("404 Not Found")
}

var splitCamelCaseFunc = func(s string) []string {
	var result []string
	start := 0

	for i := 1; i < len(s); i++ {
		if unicode.IsUpper(rune(s[i])) {
			result = append(result, s[start:i])
			start = i
		}
	}

	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}
