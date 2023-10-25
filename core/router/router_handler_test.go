// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"testing"
)

type Controller struct{}

func (c *Controller) GetV1UsageWeek() error {
	return nil
}

func (c *Controller) GetV1UsageWeek0() error {
	return nil
}

func (c *Controller) PublicGetV1UsageWeek0() error {
	return nil
}

func TestJetRouter_RegisterRouter(t *testing.T) {
	xlog.SetOutputLevel(xlog.Ldebug)
	Register(Controller{}, &Controller{})
}
