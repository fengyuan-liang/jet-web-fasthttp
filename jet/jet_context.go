// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"jet-web/pkg/xlog"
)

// Ctx represents the Context which hold the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers and so on.
type Ctx struct {
	log xlog.Logger // log for context
}
