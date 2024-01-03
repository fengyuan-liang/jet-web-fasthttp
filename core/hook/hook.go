// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package hook

import (
	"reflect"
)

type Hook struct {
	PostParamsParseHooks   []*reflect.Value
	PostMethodExecuteHooks []*reflect.Value
	PreMethodExecuteHooks  []*reflect.Value
}

// GenHook see handle#AddHook
func GenHook(rcvr *reflect.Value) *Hook {
	hook := new(Hook)
	if methodByName := rcvr.MethodByName("PostParamsParseHook"); methodByName.IsValid() {
		hook.PostParamsParseHooks = append(hook.PostParamsParseHooks, &methodByName)
	}

	if methodByName := rcvr.MethodByName("PostMethodExecuteHook"); methodByName.IsValid() {
		hook.PostMethodExecuteHooks = append(hook.PostMethodExecuteHooks, &methodByName)
	}

	if methodByName := rcvr.MethodByName("PreMethodExecuteHook"); methodByName.IsValid() {
		hook.PreMethodExecuteHooks = append(hook.PreMethodExecuteHooks, &methodByName)
	}
	return hook
}

func (hook *Hook) PostParamsParse(param reflect.Value) (err error) {
	for _, h := range hook.PostParamsParseHooks {
		// call
		result := h.Call([]reflect.Value{param})

		// handle error
		if !result[0].IsNil() {
			err = result[0].Interface().(error)
			return
		}
	}
	return
}

func (hook *Hook) HasMethodExecuteHook() bool {
	return len(hook.PostMethodExecuteHooks) != 0
}

func (hook *Hook) PostMethodExecuteHook(param reflect.Value) (data any, err error) {
	for _, h := range hook.PostMethodExecuteHooks {
		// call
		result := h.Call([]reflect.Value{param})
		// handle error
		if !result[1].IsNil() {
			err = result[1].Interface().(error)
			return
		}
		// handle data
		if !result[0].IsNil() {
			data = result[0].Interface()
		}
	}
	return
}

func (hook *Hook) HasPreMethodExecuteHooks() bool {
	return len(hook.PostMethodExecuteHooks) != 0
}

func (hook *Hook) PreMethodExecuteHook(ctx reflect.Value) (err error) {
	for _, h := range hook.PreMethodExecuteHooks {
		// call
		result := h.Call([]reflect.Value{ctx})
		// handle error
		if !result[0].IsNil() {
			err = result[0].Interface().(error)
			return
		}
	}
	return
}
