package hook

import "github.com/fengyuan-liang/jet-web-fasthttp/core/context"

// ======================================================

var (
	PostJetCtxInitHooks = make([]func(ctx context.Ctx), 0)
)
