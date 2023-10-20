// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"github.com/valyala/fasthttp"
	"jet-web/core/router"
	"jet-web/pkg/xlog"
)

func Run(addr string) error {
	xlog.NewWith("jet_log").Infof("jet server start on [%s]", addr)
	return fasthttp.ListenAndServe(addr, router.ServeHTTP)
}

func Register(rcvrs ...interface{}) {
	router.Register(rcvrs...)
}
