// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"errors"
	"jet-web/pkg/constant"
	"reflect"
	"strings"
)

// ---------------------------------------------------------------------------

var ErrMethodPrefix = errors.New("invalid method name prefix")

type HandlerFactory map[string]CreatorFunc

var defaultHandlerCreator = HandlerCreator{}.New

var Factory = HandlerFactory{
	constant.MethodGet:    handlerGetCreator,
	constant.MethodPost:   defaultHandlerCreator,
	constant.MethodPut:    defaultHandlerCreator,
	constant.MethodDelete: defaultHandlerCreator,
}

func (factory HandlerFactory) Create(rcvr *reflect.Value, method *reflect.Method) (prefix string, h IHandler, err error) {
	prefix, ok := prefixOf(method.Name)
	if !ok {
		return "", nil, ErrMethodPrefix
	}
	if creatorFunc, ok := factory[strings.ToUpper(prefix)]; ok {
		h, err = creatorFunc(rcvr, method)
		if err != nil {
			return "", nil, err
		}
		return prefix, h, nil
	}
	return "", nil, ErrMethodPrefix
}

type creatorRegisterFunc func(factory HandlerFactory) (creatorFunc CreatorFunc)

// RegisterFactory
// You can replace the default factory implementation,
// but you must bear the risk for this.
// Calmly, you can also add new factories to implement your own processing logic.
//
// You can enhance the default factory or do some pre- or post-processing,
// like before- or after-RegisterHook.
func (factory HandlerFactory) RegisterFactory(prefix string, newFactoryFunc creatorRegisterFunc) {
	factory[prefix] = newFactoryFunc(factory)
}
