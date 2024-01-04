// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

// --------------------------------------------------------------------

var httpTraceLog = func() *xlog.Logger {
	logger := xlog.NewWith("jet")
	logger.SetCalldPath(3)
	logger.SetFlags(xlog.LstdFlags | xlog.Llevel)
	return logger
}()

func httpTrace(start time.Time, traceName string) {
	httpTraceLog.Infof(" %v | elapsed [%v]", traceName, time.Since(start))
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

// isTerminal Check if the output destination is a terminal
var isTerminal = func() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}()

func TraceHttpReq(ctx *fasthttp.RequestCtx, start time.Time) {
	var (
		statusCode = ctx.Response.StatusCode()
		method     = string(ctx.Method())
		path       = string(ctx.Path())
	)
	if isTerminal {
		// status method path
		httpTrace(start,
			fmt.Sprintf("|%s %3d %s| |%s %s %s| %v",
				colorForStatus(statusCode), statusCode, reset,
				colorForMethod(method), method, reset,
				path,
			),
		)
	} else {
		// status method path
		httpTrace(start,
			fmt.Sprintf("%v | %v | %v",
				statusCode,
				method,
				path,
			),
		)
	}
}
