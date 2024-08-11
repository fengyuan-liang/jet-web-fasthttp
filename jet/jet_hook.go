// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/hook"
)

// global hook

func AddPostJetCtxInitHook(f func(ctx context.Ctx)) {
	hook.PostJetCtxInitHooks = append(hook.PostJetCtxInitHooks, f)
}
