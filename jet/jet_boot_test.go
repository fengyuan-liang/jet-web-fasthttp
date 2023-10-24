//go:build !ignore

// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"errors"
	"jet-web/core/context"
	"jet-web/pkg/utils"
	"jet-web/pkg/xlog"
	"os"
	"testing"
)

type jetController struct{}

var bootTestLog = xlog.NewWith("boot_test_log")

func (j *jetController) GetV1UsageWeek0(args *context.Args) error {
	bootTestLog.Infof("GetV1UsageWeek %v", *args)
	return errors.New(utils.ObjToJsonStr(args.CmdArgs))
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (j *jetController) GetV1Usage0Week(args *context.Args) (*Person, error) {
	bootTestLog.Infof("GetV1Usage0Week %v", *args)
	return &Person{
		Name: "张三",
		Age:  18,
	}, nil
}

func (j *jetController) GetV1UsageWeek(args string) error {
	bootTestLog.Info("GetV1UsageWeek", args)
	return nil
}

func (j *jetController) GetV1UsageWeekk0(args *context.Args) error {
	bootTestLog.Infof("GetV1UsageWeekk0 %v", *args)
	return errors.New(utils.ObjToJsonStr(args.CmdArgs))
}

func TestJetBoot(t *testing.T) {
	if os.Getenv("SKIP_TESTS") != "" {
		t.Skip("Skipping JetBoot test")
	}
	xlog.SetOutputLevel(xlog.Ldebug)
	Register(&jetController{})
	t.Logf("err:%v", Run(":8080"))
}
