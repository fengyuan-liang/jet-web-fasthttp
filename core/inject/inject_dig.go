// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inject

import "go.uber.org/dig"

var container = dig.New()

func Container() *dig.Container {
	return container
}

func Invoke(i interface{}) {
	if err := container.Invoke(i); err != nil {
		panic(err)
	}
}

func Provide(constructs ...any) {
	for _, construct := range constructs {
		if err := container.Provide(construct); err != nil {
			panic(err)
		}
	}
}
