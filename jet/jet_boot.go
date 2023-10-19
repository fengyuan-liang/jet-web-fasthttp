// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"github.com/valyala/fasthttp"
	"jet-web/core/router"
)

func Run(addr string) error {
	return fasthttp.ListenAndServe(addr, router.ServeHTTP)
}

func Register(controller any) {

}
