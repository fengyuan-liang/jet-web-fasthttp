// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"jet-web/pkg/constant"
	"reflect"
)

// ---------------------------------------------------------------------------

type HandlerFactory []*HandlerCreator

var defaultHandlerCreator = HandlerCreator{}.New

var Factory = HandlerFactory{
	{constant.MethodGet, handlerGetCreator},
	{constant.MethodPost, defaultHandlerCreator},
	{constant.MethodPut, defaultHandlerCreator},
	{constant.MethodDelete, defaultHandlerCreator},
}

func AddHandlerCreator(methodPrefix string, creatorFunc CreatorFunc) {
	Factory = append(Factory, &HandlerCreator{methodPrefix, creatorFunc})
}

func (r *HandlerFactory) CreateHandlerFunc(rcvr *reflect.Value, method *reflect.Value) (string, IHandler, error) {
	return "", nil, nil
}
