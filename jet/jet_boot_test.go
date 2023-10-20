//go:build !ignore

// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"jet-web/pkg/xlog"
	"os"
	"testing"
)

type jetController struct{}

var bootTestLog = xlog.NewWith("[boot_test_log]")

func (j *jetController) GetV1UsageWeek0() error {
	bootTestLog.Info("GetV1UsageWeek")
	return nil
}

func TestJetBoot(t *testing.T) {
	if os.Getenv("SKIP_TESTS") != "" {
		t.Skip("Skipping JetBoot test")
	}
	Register(&jetController{})
	t.Logf("err:%v", Run(":8080"))
}
