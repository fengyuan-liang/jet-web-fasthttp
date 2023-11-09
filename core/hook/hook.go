// Copyright 2023 QINIU. All rights reserved
// @Description:
// @Version: 1.0.0
// @Date: 2023/11/09 16:01
// @Author: liangfengyuan@qiniu.com

package hook

import "reflect"

type Hook struct {
	PostParamsParseHooks   []*reflect.Value
	PostMethodExecuteHooks []*reflect.Value
}

func GenHook(rcvr *reflect.Value) *Hook {
	hook := new(Hook)
	if methodByName := rcvr.MethodByName("PostParamsParseHook"); methodByName.IsValid() {
		hook.PostParamsParseHooks = append(hook.PostParamsParseHooks, &methodByName)
	}

	if methodByName := rcvr.MethodByName("PostMethodExecuteHook"); methodByName.IsValid() {
		hook.PostMethodExecuteHooks = append(hook.PostMethodExecuteHooks, &methodByName)
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
