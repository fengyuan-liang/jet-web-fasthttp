// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inject

import "go.uber.org/dig"

// IJetController By marking it with IJetController,
// this indicates that it is a Controller that needs to be registered with Jet.
// Jet will automatically handle the mounted routes and resolve them.
type IJetController interface {
}

type JetControllerResult struct {
	dig.Out
	Handler IJetController `group:"server"`
}

type JetControllerList struct {
	dig.In
	Handlers []IJetController `group:"server"`
}

func NewJetController(controller IJetController) JetControllerResult {
	return JetControllerResult{
		Handler: controller,
	}
}
