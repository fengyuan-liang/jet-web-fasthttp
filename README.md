# Jet 🛩

一款和gin不太一样的golang web服务器

## Overview

- 异常简洁的路由规则，再也不用像gin一样写繁琐的路由，并且自动解析参数
- 依赖注入 & 依赖倒置 & 开闭原则
- 集成 fasthttp
- HTTP/H2C Server & Client
- 集成普罗米修斯
- AOP Worker & 无侵入 Context
- 可扩展组件 Infrastructure
- DDD & 六边形架构
- 领域事件 & 消息队列组件
- CQS & 聚合根
- CRUD & PO Generate
- 一级缓存 & 二级缓存 & 防击穿

## usage

```go
// 在Jet中 路由是挂载在Controller上的，通过Controller进行路由分组
type jetController struct{}

var bootTestLog = xlog.NewWith("boot_test_log")

func TestJetBoot(t *testing.T) {
	jet.Register(&jetController{})
	t.Logf("err:%v", jet.Run(":8080"))
}

// ----------------------------------------------------------------------

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

// ----------------------------------------------------------------------

// 我们会尽可能的找到您需要的参数并将参数注入到您的结构体参数中
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
// hook
func (j *jetController) PostParamsParseHook(param any) error {
    // 可以通过参数注入完后的hook对参数进行校验，比如使用`validated`库进行校验
	if err := utils.Struct(param); err != nil {
		return errors.New(utils.ProcessErr(param, err))
	}
	return nil
}
// hook
func (j *jetController) PostMethodExecuteHook(param any) (data any, err error) {
	// 你可以通过controller方法执行完后的hook来restful方式的处理返回结果
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

```

我们注意到，`UserController`的方法比较有意思，叫`GetV1UsageWeek`，其实这代表着我们有一个接口`v1/usage/week`已经写好了，请求方式为`Get`，我们请求的参数会自动注入到`r *Args`中

```shell
$ curl http://localhost/v1/usage/week?form_param1=1
{"request_id":"ZRgQg3Osptrx","code":200,"message":"success","data":"1"}
```

如果想要定义`v1/usage/week/1`的形式，或者`v1/usage/1/week`，我们可以使用`0`或其他符号填充名字

```go
GetV1UsageWeek0 -> v1/usage/week/1 // 0的位置表示要接受一个可变的参数
GetV1Usage0Week -> v1/usage/1/week
```

参数会默认注入到`CmdArgs`中

```go
func (u *UserController) GetV1Usage0Week(r *Args, env *rpc.Env) (*api.Response, error) {
	return api.Success(xlog.GenReqId(), r.CmdArgs), nil
}
```

```shell
$ curl http://localhost/v1/usage/1/week
{"request_id":"H5OQ4Jg0yBtg","code":200,"message":"success","data":["1"]}
```

## 更新计划

### 1. Hook

#### 1.1 参数相关

- [x] 支持通过挂载hook对参数进行预解析、自定义参数校验规则（目前支持hook有）
  - PostParamsParseHook
  - PostRouteMountHook
  - PostMethodExecuteHook
  - PreMethodExecuteHooks
- [x] 添加hook注入自定义的`context`，便于进行鉴权以及链路追踪等操作

### 2. 🤡Aspect（切面）支持

#### 2.1 常规切面

- [ ] 前置、后置、异常、环绕、最终五种切面

### 3. 路由策略

- [ ] 通过controller自定义路由前缀

### 4. 依赖注入支持

