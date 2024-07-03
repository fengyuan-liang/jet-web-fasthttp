//go:build !ignore

// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"errors"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"os"
	"testing"
	"time"
)

type jetController struct {
	BaseJetController
}

func NewDemoController() ControllerResult {
	return NewJetController(&jetController{})
}

var bootTestLog = xlog.NewWith("boot_test_log")

func TestJetBoot(t *testing.T) {
	if os.Getenv("SKIP_TESTS") != "" {
		t.Skip("Skipping JetBoot test")
	}
	xlog.SetOutputLevel(xlog.Ldebug)
	//Register(&jetController{})
	AddMiddleware(RecoverJetMiddleware, TraceJetMiddleware)
	Provide(NewDemoController)
	Run(":8080")
}

// ---------------------------  hooks  ----------------------------------

// curl http://localhost:8080/v1/usage/111/week  =>  {"code":401,"data":{},"msg":"bad token"}
// if add -H "Authorization: <your_token_here>"  =>  {"code":200,"data":{},"msg":"msg"}
//func (j *jetController) PreMethodExecuteHook(ctx context.Ctx) (err error) {
//	if authorizationHeader := string(ctx.Request().Header.Peek("Authorization")); authorizationHeader == "" {
//		ctx.Response().SetStatusCode(401)
//		errInfo := map[string]any{"code": 401, "data": ctx.Keys(), "msg": "bad token"}
//		err = errors.New(utils.ObjToJsonStr(errInfo))
//	}
//	return
//}

// ----------------------------------------------------------------------

type req struct {
	Id   int    `json:"id" form:"id" validate:"required" reg_err_info:"is empty"`
	Name string `json:"name" form:"name" validate:"required" reg_err_info:"is empty"`
}

func (j *jetController) PostV1UsageContext(ctx Ctx, req *req) (maps.IMap[string, any], error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Logger().Infof("req:%v", req)
	ctx.Put("request uri", ctx.Request().URI().String())
	ctx.Put("traceId", ctx.Logger().ReqId)
	ctx.Put("req", req)
	value := ctx.FastHttpCtx().UserValue(11)
	ctx.Logger().Infof("value:%v", value)

	return ctx.Keys(), nil
}

type Response struct {
	RequestId string `json:"request_id,omitempty"` //请求ID
	Code      int    `json:"code"`                 //错误码，200 成功，其他失败
	Message   string `json:"message,omitempty"`    //错误信息
	Data      any    `json:"data,omitempty"`
}

func (j *jetController) GetV1UsageContext(ctx Ctx, req *req) (*Response, error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Logger().Infof("req:%v", req)
	ctx.Put("traceId", ctx.Logger().ReqId)
	return &Response{Data: ctx.Keys()}, nil
}

func (j *jetController) GetV1UsageContext0(ctx Ctx, args *context.Args) (map[string]any, error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Put("request uri", ctx.Request().URI().String())
	ctx.Put("traceId", ctx.Logger().ReqId)
	ctx.Put("args", args)
	return map[string]any{"code": 200, "data": ctx.Keys(), "msg": "ok"}, nil
}

func (j *jetController) GetV1UsageWeek0(args *context.Args) error {
	time.Sleep(time.Second * 2)
	bootTestLog.Infof("GetV1UsageWeek %v", *args)
	return errors.New(utils.ObjToJsonStr(args.CmdArgs))
}

type Person struct {
	Name string `json:"name" form:"name"`
	Age  int    `json:"age" form:"age"`
}

func (j *jetController) GetV1Usage0Week() (*Person, error) {
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

func (j *jetController) PostFormCallTest(p *Person) (*Person, error) {
	bootTestLog.Infof("GetV1UsageWeekk0 %v", *p)
	return p, nil
}
