// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package hook

type PreRouteSetupHook interface {
	PreRouteSetup()
}

// PostParamsParseHook Hook triggered after parameter parsing is complete
type PostParamsParseHook interface {
	PostParamsParseHook(param any) error
}

// PostRouteMountHook Hook triggered after route is mounted
type PostRouteMountHook interface {
	PostRouteMount()
}

// PostMethodExecuteHook Hook triggered after method execution but before returning
type PostMethodExecuteHook interface {
	PostMethodExecuteHook(param any) (data any, err error)
}
