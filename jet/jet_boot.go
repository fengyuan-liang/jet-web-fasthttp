// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/inject"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/commands"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"go.uber.org/dig"
	"time"
)

type Server struct{}

var (
	jetLog         = xlog.NewWith("jet_log")
	startTime      time.Time
	localAddr      string
	fastHttpServer *fasthttp.Server
)

func SetFastHttpServer(server *fasthttp.Server) {
	fastHttpServer = server
}

func (j *Server) Initialize() (err error) {
	if localAddr == "" {
		err = errors.New("addr is empty")
		return
	}
	return
}

func (j *Server) RunLoop() {
	jetLog.Infof("jet server start on [%s] elapsed [%v]", localAddr, time.Since(startTime))
	if fastHttpServer == nil {
		jetLog.Errorf("%v", fasthttp.ListenAndServe(localAddr, router.ServeHTTP))
	} else {
		jetLog.Errorf("%v", fastHttpServer.ListenAndServe(localAddr))
	}
}

func (j *Server) Destroy() {
	jetLog.Info("Jet server Destroy...")
}

func NewByInject(jetControllerList inject.JetControllerList) commands.MainInstance {
	// provide by inject
	if len(jetControllerList.Handlers) != 0 {
		for _, handler := range jetControllerList.Handlers {
			router.Register(handler)
		}
	}
	return &Server{}
}

func New() *Server {
	return &Server{}
}

func Run(addr string) {
	startTime = time.Now()
	localAddr = addr
	inject.Provide(NewByInject)
	inject.Invoke(func(srv commands.MainInstance) {
		// add middleware
		if len(middlewares) != 0 {
			for _, middleware := range middlewares {
				if nextRouter, err := middleware(router.DefaultJetRouter); err == nil {
					router.DefaultJetRouter = nextRouter
				}
			}
		}
		commands.Run(srv)
	})
}

func Register(rcvrs ...any) {
	router.Register(rcvrs...)
}

func Invoke(i any) {
	inject.Invoke(i)
}

func Provide(constructs ...any) {
	inject.Provide(constructs...)
}

type ControllerResult struct {
	dig.Out
	Handler inject.IJetController `group:"server"`
}

type IJetController interface {
	inject.IJetController
}

func NewJetController(controller IJetController) ControllerResult {
	return ControllerResult{
		Handler: controller,
	}
}

// ----------------------------------------------------------------------

// BaseJetController Provide some basic hooks, such as parameter validation and restful style returns
type BaseJetController struct {
	IJetController
}

func (BaseJetController) PostParamsParseHook(param any) (err error) {
	if err = utils.Struct(param); err != nil {
		err = errors.New(utils.ObjToJsonStr(map[string]any{"code": 400, "message": "bad request", "data": utils.ProcessErr(param, err)}))
	}
	return
}

// PostMethodExecuteHook restful
func (BaseJetController) PostMethodExecuteHook(param any) (data any, err error) {
	// restful
	return utils.ObjToJsonStr(param), nil
}
