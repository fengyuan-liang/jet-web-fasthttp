// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"time"
)

var traceElapsed = func() *xlog.Logger {
	logger := xlog.NewWith("traceElapsed")
	logger.SetCalldPath(3)
	logger.SetFlags(xlog.LstdFlags | xlog.Llevel)
	return logger
}()

func TraceElapsed(start time.Time) {
	traceElapsed.Infof("trace elapsed [%v]", time.Since(start))
}

func TraceElapsedByName(start time.Time, traceName string) {
	traceElapsed.Infof(" trace[%v] elapsed[%v]", traceName, time.Since(start))
}
