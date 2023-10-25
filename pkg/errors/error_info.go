// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"runtime"
	"strconv"
	"strings"
)

const prefix = " ==> "

// --------------------------------------------------------------------

func New(msg string) error {
	return errors.New(msg)
}

// --------------------------------------------------------------------

type errorDetailer interface {
	ErrorDetail() string
}

func Detail(err error) string {
	if e, ok := err.(errorDetailer); ok {
		return e.ErrorDetail()
	}
	return prefix + err.Error()
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err  error
	Why  error
	Cmd  []interface{}
	File string
	Line int
}

func shortFile(file string) string {
	pos := strings.LastIndex(file, "/src/")
	if pos != -1 {
		return file[pos+5:]
	}
	return file
}

func Info(err error, cmd ...interface{}) *ErrorInfo {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
	}
	return &ErrorInfo{Cmd: cmd, Err: Err(err), File: file, Line: line}
}

// InfoEx file and line tracing may have problems with go1.9, see related issue: https://github.com/golang/go/issues/22916
func InfoEx(skip int, err error, cmd ...interface{}) *ErrorInfo {
	var e *ErrorInfo
	if errors.As(err, &e) {
		err = e.Err
	}
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
	}
	return &ErrorInfo{Cmd: cmd, Err: err, File: file, Line: line}
}

func (r *ErrorInfo) Cause() error {
	return r.Err
}

func (r *ErrorInfo) Error() string {
	return r.Err.Error()
}

func (r *ErrorInfo) ErrorDetail() string {
	e := prefix + shortFile(r.File) + ":" + strconv.Itoa(r.Line) + ": " + r.Err.Error() + " ~ " + fmt.Sprintln(r.Cmd...)
	if r.Why != nil {
		e += Detail(r.Why)
	} else {
		e = e[:len(e)-1]
	}
	return e
}

func (r *ErrorInfo) Detail(err error) *ErrorInfo {
	r.Why = err
	return r
}

func (r *ErrorInfo) Method() (cmd string, ok bool) {
	if len(r.Cmd) > 0 {
		if cmd, ok = r.Cmd[0].(string); ok {
			if pos := strings.Index(cmd, " "); pos > 1 {
				cmd = cmd[:pos]
			}
		}
	}
	return
}

func (r *ErrorInfo) xlogMessage() string {
	detail := r.ErrorDetail()
	if cmd, ok := r.Method(); ok {
		detail = cmd + " failed:\n" + detail
	}
	return detail
}

// Warn deprecated. please use (*ErrorInfo).xlogWarn
func (r *ErrorInfo) Warn() *ErrorInfo {
	xlog.Std.Output("", xlog.Lwarn, 2, r.xlogMessage())
	return r
}

func (r *ErrorInfo) xlogWarn(reqId string) *ErrorInfo {
	xlog.Std.Output(reqId, xlog.Lwarn, 2, r.xlogMessage())
	return r
}

func (r *ErrorInfo) xlogError(reqId string) *ErrorInfo {
	xlog.Std.Output(reqId, xlog.Lerror, 2, r.xlogMessage())
	return r
}

func (r *ErrorInfo) xlog(level int, reqId string) *ErrorInfo {
	xlog.Std.Output(reqId, level, 2, r.xlogMessage())
	return r
}

// --------------------------------------------------------------------

type causer interface {
	Cause() error
}

func Err(err error) error {
	if e, ok := err.(causer); ok {
		if diag := e.Cause(); diag != nil {
			return diag
		}
	}
	return err
}

// --------------------------------------------------------------------
