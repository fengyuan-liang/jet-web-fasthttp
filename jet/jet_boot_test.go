//go:build !ignore

// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"os"
	"testing"
)

type jetController struct{}

var bootTestLog = xlog.NewWith("boot_test_log")

func TestJetBoot(t *testing.T) {
	if os.Getenv("SKIP_TESTS") != "" {
		t.Skip("Skipping JetBoot test")
	}
	xlog.SetOutputLevel(xlog.Ldebug)
	Register(&jetController{})
	t.Logf("err:%v", Run(":8080"))
}

type req struct {
	Id   int    `json:"id" validate:"required" reg_err_info:"is empty"`
	Name string `json:"name" validate:"required" reg_err_info:"is empty"`
}

func (j *jetController) PostV1UsageContext(ctx Ctx, req *req) (map[string]any, error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Logger().Infof("req:%v", req)
	ctx.Put("request uri", ctx.Request().URI().String())
	ctx.Put("traceId", ctx.Logger().ReqId)
	ctx.Put("req", req)
	return ctx.Keys(), nil
}

func (j *jetController) PostParamsParseHook(param any) error {
	if err := utils.Struct(param); err != nil {
		return errors.New(utils.ProcessErr(param, err))
	}
	return nil
}

func (j *jetController) PostMethodExecuteHook(param any) (data any, err error) {
	// restful
	return utils.ObjToJsonStr(param), nil
}

func (j *jetController) GetV1UsageContext0(ctx Ctx, args *context.Args) (map[string]any, error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Put("request uri", ctx.Request().URI().String())
	ctx.Put("traceId", ctx.Logger().ReqId)
	ctx.Put("args", args)
	return map[string]any{"code": 200, "data": ctx.Keys(), "msg": "ok"}, nil
}

func (j *jetController) GetV1UsageWeek0(args *context.Args) error {
	bootTestLog.Infof("GetV1UsageWeek %v", *args)
	return errors.New(utils.ObjToJsonStr(args.CmdArgs))
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (j *jetController) GetV1Usage0Week(args *context.Args) (*Person, error) {
	//bootTestLog.Infof("GetV1Usage0Week %v", *args)
	return &Person{
		Name: "张三",
		Age:  18,
	}, nil
}

func (j *jetController) GetV1UsageWeek(args string) (map[string]string, error) {
	bootTestLog.Info("GetV1UsageWeek", args)
	return map[string]string{"args": args}, nil
}

func (j *jetController) GetV1UsageWeekk0(args *context.Args) error {
	bootTestLog.Infof("GetV1UsageWeekk0 %v", *args)
	return errors.New(utils.ObjToJsonStr(args.CmdArgs))
}
