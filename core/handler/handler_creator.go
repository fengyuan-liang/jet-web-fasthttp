// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import "reflect"

// ---------------------------------------------------------------------------

// CreatorFunc define func of HandlerCreator
type CreatorFunc = func(rcvr *reflect.Value, method *reflect.Method) (IHandler, error)

type HandlerCreator struct {
	MethodPrefix string
	Creator      CreatorFunc
}

func (p *HandlerCreator) New(rcvr *reflect.Value, method *reflect.Method) (IHandler, error) {
	return nil, nil
}

var handlerGetCreator = func(rcvr *reflect.Value, method *reflect.Method) (IHandler, error) {
	h := new(handler)
	return h, nil
}
