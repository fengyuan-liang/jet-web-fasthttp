# Jet 🛩

A Golang web server that is different from gin.

## Overview

- Exceptionally concise routing rules, no longer need to write tedious routes like gin, and automatically parse parameters.
- Dependency injection & inversion of control & open-closed principle.
- Integration with fasthttp.
- More fine-grained Hook support.
- DDD & Hexagonal architecture.
- CQS & Aggregate Root.
- First-level cache & Second-level cache & Anti-cache penetration (not yet implemented).
- Integration with Prometheus (not yet implemented).
- AOP integration (not yet implemented).

## usage

```go
go get github.com/fengyuan-liang/jet-web-fasthttp
```

## 使用说明

```go
// In Jet, routes are mounted on controllers, and route grouping is done through controllers.
type jetController struct{}

var bootTestLog = xlog.NewWith("boot_test_log")

func TestJetBoot(t *testing.T) {
	jet.Register(&jetController{})
	t.Logf("err:%v", jet.Run(":8080"))
}

// ------------------------------  HOOK  ----------------------------------

// After parameter parsing is completed, you can use hooks to perform parameter validation. For example, you can use validated for validation.
func (j *jetController) PostParamsParseHook(param any) error {
	if err := utils.Struct(param); err != nil {
		return errors.New(utils.ProcessErr(param, err))
	}
	return nil
}

func (j *jetController) PostMethodExecuteHook(param any) (data any, err error) {
    // You can use hooks after the execution of controller methods to handle the result in a RESTful manner.
	return utils.ObjToJsonStr(param), nil
}

// curl http://localhost:8080/v1/usage/111/week  =>  {"code":401,"data":{},"msg":"bad token"}
// if add -H "Authorization: <your_token_here>"  =>  {"code":200,"data":{},"msg":"msg"}
func (j *jetController) PreMethodExecuteHook(ctx context.Ctx) (err error) {
	if authorizationHeader := string(ctx.Request().Header.Peek("Authorization")); authorizationHeader == "" {
		ctx.Response().SetStatusCode(401)
		errInfo := map[string]any{"code": 401, "data": ctx.Keys(), "msg": "bad token"}
		err = errors.New(utils.ObjToJsonStr(errInfo))
	}
	return
}

// --------------------------------  ROUTER  -------------------------------

// We will make every effort to find the parameters you need and inject them into your struct parameters.
type req struct {
	Id   int    `json:"id" validate:"required" reg_err_info:"cannot empty"`
	Name string `json:"name" validate:"required" reg_err_info:"cannot empty"`
}

func (j *jetController) PostV1UsageContext(ctx jet.Ctx, req *req) (map[string]any, error) {
	ctx.Logger().Info("GetV1UsageContext")
	ctx.Logger().Infof("req:%v", req)
	ctx.Put("request uri", ctx.Request().URI().String())
	ctx.Put("traceId", ctx.Logger().ReqId)
	ctx.Put("req", req)
	return ctx.Keys(), nil
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
	// bootTestLog.Infof("GetV1Usage0Week %v", *args)
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

```

We noticed that the `UserController` method is quite interesting, named `GetV1UsageWeek`. In fact, this indicates that we have an endpoint `v1/usage/week` already implemented, with a `GET` request method. The requested parameters will be automatically injected into `r *Args`.

```shell
$ curl http://localhost/v1/usage/week?form_param1=1
{"request_id":"ZRgQg3Osptrx","code":200,"message":"success","data":"1"}
```

If you want to define the form `v1/usage/week/1` or `v1/usage/1/week`, you can use `0` or any other symbol as a placeholder in the route definition. For example, you can define the route as `v1/usage/:placeholder/week`, where `:placeholder` can be replaced with `0`, `1`, or any other desired value.

```go
GetV1UsageWeek0 -> v1/usage/week/1 // The position of 0 indicates that it is meant to accept a variable parameter.
GetV1Usage0Week -> v1/usage/1/week
```

The parameters will be automatically injected into `CmdArgs` by default.

```go
func (u *UserController) GetV1Usage0Week(r *Args, env *rpc.Env) (*api.Response, error) {
	return api.Success(xlog.GenReqId(), r.CmdArgs), nil
}
```

```shell
$ curl http://localhost/v1/usage/1/week
{"request_id":"H5OQ4Jg0yBtg","code":200,"message":"success","data":["1"]}
```

### example

```go
func main() {
	//jet.Register(&DemoController{})
	xlog.SetOutputLevel(xlog.Ldebug)
	jet.AddMiddleware(jet.TraceJetMiddleware)
	jet.Run(":8080")
}

func init() {
	jet.Provide(NewDemoController)
}

func NewDemoController() jet.ControllerResult {
	return jet.NewJetController(&DemoController{})
}

type BaseController struct {
	jet.IJetController
}

func (BaseController) PostParamsParseHook(param any) error {
	if err := utils.Struct(param); err != nil {
		return errors.New(utils.ProcessErr(param, err))
	}
	return nil
}

// PostMethodExecuteHook restful
func (BaseController) PostMethodExecuteHook(param any) (data any, err error) {
	// restful
	return utils.ObjToJsonStr(param), nil
}

type DemoController struct {
	BaseController
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
// 路由 get /v1/usage/{id}/week 已经可以访问了
func (j *DemoController) GetV1Usage0Week(ctx jet.Ctx, args *jet.Args) (*Person, error) {
	ctx.Logger().Infof("GetV1Usage0Week %v", *args)
	return &Person{
		Name: "张三",
		Age:  18,
	}, nil
}
```

## Update Plan

### 1. Hook

#### 1.1 Parameter-related

Support pre-parsing and custom parameter validation rules through mounted hooks (currently supported hooks include)

- [x] PostParamsParseHook
- [x] PostRouteMountHook
- [x] PostMethodExecuteHook
- [x] PreMethodExecuteHook

other

- [x] Add hook injection for custom `context` to facilitate authentication and tracing operations.

### 2. 🤡Aspect support

#### 2.1 Aspect 

- [ ] The five types of aspects are: before, after, exception, around, and finally.

### 3. router

- [ ] Customize the route prefix through the controller.

### 4. inject support

In Jet, dependency injection (inject) is a fundamental concept, and almost all functionalities in Jet are accomplished through dependency injection (Jet relies on `dig` for its underlying dependency injection implementation).

For example, you can provide `JetController` to Jet, and it will automatically detect and parse the routes.

```go
type jetController struct {
	inject.IJetController
}

func NewDemoController() inject.JetControllerResult {
	return inject.NewJetController(&jetController{})
}

func main() {
  xlog.SetOutputLevel(xlog.Ldebug)
    // inject
	jet.Provide(NewDemoController)
	jet.Run(":8080")
}
```

Jet recommends incorporating dependency injection throughout the entire development lifecycle of the program, including in the `repo`, `service`, and `controller` layers of the MVC architecture, or in the `domain` layer of the DDD architecture.

You can use the following approach, combined with the `init` method, to automatically inject dependencies into Jet and manage the lifecycle of the entire program:

```go
package main

import (
	_ "xxx/apps/xxx/internal/component"
	_ "xxx/apps/xxx/internal/controller"
	_ "xxx/apps/xxx/internal/server"
	_ "xxx/domain/repo"
)

func main() {
	jet.Run(":8080")
}
```

In other domain layers, we need to register components with Jet.

```go
// xxxController.go

func init() {
  // provide your 
  jet.Provide(NewXxxController)
}

type XxxController struct {
  xxxRepo repo.XxxRepo
}

func NewXxxController(xxxRepo repo.XxxRepo) jet.ControllerResult {
  return jet.NewJetController(&jetController{
    xxxRepo: xxxRepo
  })
}
```

### 5. middleware 

The support for middleware in Jet is straightforward, direct, and clear. When we add multiple middleware, Jet executes them from the inside out, meaning that the middleware added later will be executed first, and the ones added earlier will be executed later.

#### Logging middleware

```go
func main() {
	jet.Register(&jetController{})
	jet.AddMiddleware(TraceJetMiddleware)
	jet.Run(":8080")
}

func TraceJetMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		defer utils.TraceHttpReq(ctx, time.Now())
		next.ServeHTTP(ctx)
	}), nil
}
```

When a request is initiated

```shell
$ ➜  ~ curl http://localhost:8080/v1/usage/week/111
["111"]%
```

We can observe the output in a very intuitive manner.

```shell
2024/01/04 16:31:55.379274 [jet][INFO] | 200 | | GET | /v1/usage/week/111 | elapsed [2.00150788s]
```

When an error is returned during the invocation and subsequent execution of the middleware, the following middleware will no longer be executed.

#### recover middleware

you can use jet default middleware

```go
func main() {
  jet.AddMiddleware(RecoverJetMiddleware)
  jet.Run(":8080")
}

```

`Jet` will be return `Internal Server Error`，http code is`500`

![image-20240105110436328](https://cdn.fengxianhub.top/resources-master/image-20240105110436328.png)

Certainly, you can also customize your own middleware. However, please note that middleware is executed in the order they are added, with the later-added middleware being executed after the earlier-added ones. To avoid interference from the `recover` middleware on the logic of other middleware, Jet recommends adding your middleware in the first position.

```go
// If you return xxx, err from a middleware, the subsequent middleware will not be executed.
func RecoverJetMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				handler.FailServerInternalErrorHandler(ctx, "Internal Server Error")
				utils.PrintPanicInfo("Your server has experienced a panic, please check the stack log below")
				debug.PrintStack()
			}
		}()
		next.ServeHTTP(ctx)
	}), nil
}
```

### 6. benchmark

```shell
$ ab -c 400 -n 20000 http://localhost:8081/v1/usage/1111/week
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 2000 requests
Completed 4000 requests
Completed 6000 requests
Completed 8000 requests
Completed 10000 requests
Completed 12000 requests
Completed 14000 requests
Completed 16000 requests
Completed 18000 requests
Completed 20000 requests
Finished 20000 requests


Server Software:        JetServer
Server Hostname:        localhost
Server Port:            8081

Document Path:          /v1/usage/1111/week
Document Length:        76 bytes

Concurrency Level:      400
Time taken for tests:   1.661 seconds
Complete requests:      20000
Failed requests:        0
Total transferred:      4060000 bytes
HTML transferred:       1520000 bytes
Requests per second:    12041.08 [#/sec] (mean)
Time per request:       33.220 [ms] (mean)
Time per request:       0.083 [ms] (mean, across all concurrent requests)
Transfer rate:          2387.05 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0       1
Processing:     8   33   2.4     33      39
Waiting:        1   17   8.8     17      37
Total:          8   33   2.4     33      39

Percentage of the requests served within a certain time (ms)
  50%     33
  66%     33
  75%     34
  80%     34
  90%     35
  95%     36
  98%     37
  99%     38
 100%     39 (longest request)
```

The binary file occupies `14MB` of disk space, and during load testing, the memory usage is `6MB`.

![image-20240104182950530](https://cdn.fengxianhub.top/resources-master/image-20240104182950530.png)

![image-20240104183001418](https://cdn.fengxianhub.top/resources-master/image-20240104183001418.png)

### 其他更新

**2023/12/18**

请求计时中间件，see`jet.TraceJetMiddleware`

![image-20231218173140763](https://cdn.fengxianhub.top/resources-master/image-20231218173140763.png)
