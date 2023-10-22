// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import "testing"

func TestSplitCamelCaseFunc(t *testing.T) {
	var str = "PostV1Usage"
	t.Logf("%v", splitCamelCaseFunc(str, "*"))
}
