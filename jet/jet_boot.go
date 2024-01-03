// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jet

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/inject"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/commands"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"go.uber.org/dig"
	"time"
)

type Server struct {
	handlers []inject.IJetController
}

var (
	jetLog    = xlog.NewWith("jet_log")
	startTime time.Time
	localAddr string
)

func (j *Server) Initialize() (err error) {
	if localAddr == "" {
		err = errors.New("addr is empty")
		return
	}
	return
}

func (j *Server) RunLoop() {
	// provide by inject
	if len(j.handlers) != 0 {
		for _, handler := range j.handlers {
			router.Register(handler)
		}
	}
	jetLog.Infof("jet server start on [%s] elapsed [%v]", localAddr, time.Since(startTime))
	jetLog.Errorf("%v", fasthttp.ListenAndServe(localAddr, router.ServeHTTP))
}

func (j *Server) Destroy() {
	jetLog.Info("lego server Destroy...")
}

func NewByInject(jetControllerList inject.JetControllerList) commands.MainInstance {
	return &Server{handlers: jetControllerList.Handlers}
}

func New() *Server {
	return &Server{}
}

func Run(addr string) {
	startTime = time.Now()
	localAddr = addr
	inject.Provide(NewByInject)
	inject.Invoke(func(srv commands.MainInstance) {
		commands.Run(srv)
	})
}

func Register(rcvrs ...interface{}) {
	router.Register(rcvrs...)
}

func Invoke(i interface{}) {
	inject.Invoke(i)
}

func Provide(constructs ...any) {
	for _, construct := range constructs {
		inject.Provide(construct)
	}
}

type ControllerResult struct {
	dig.Out
	Handler inject.IJetController `group:"server"`
}

func NewJetController(controller inject.IJetController) ControllerResult {
	return ControllerResult{
		Handler: controller,
	}
}
