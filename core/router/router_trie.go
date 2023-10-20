// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type ITrie[V any] interface {
	Size() int
	IsEmpty() bool
	Clear()
	Contains(path string) bool
	Get(path string) V
	GetAndArgs(path string) (value V, args []string)
	Add(path string, v V) (overwrittenValue V)
	Remove(path string) (oldValue V)
	StartWith(prefix string) bool
}

type SplitPathFunc = func(k string) []string
type SplitMethodFunc = func(k string, separator string) []string

var defaultSplitPathFunc = func(path string) []string {
	split := strings.Split(path, "/")
	if len(split) > 0 && split[0] == "" {
		return split[1:]
	}
	return split
}

var defaultSplitMethodFunc = func(path string, _ string) []string {
	split := strings.Split(path, "/")
	if len(split) > 0 && split[0] == "" {
		return split[1:]
	}
	return split
}

type Trie[V any] struct {
	size            int             // size of node
	root            *TrieNode[V]    // root of trie
	splitPathFunc   SplitPathFunc   // split func of full path
	splitMethodFunc SplitMethodFunc // split func of full path
	separator       string          // separator to split dynamic router args
	// Optimized fields
	staticRouterMap map[string]V
	regex           *regexp.Regexp
	rwLock          sync.RWMutex // read-write lock
}

type TrieNode[V any] struct {
	children map[string]*TrieNode[V] // map's k is the current sub-path
	value    V                       // value
	isEnd    bool                    // judge is end
	args     []string                // dynamic router args
}

func (n *TrieNode[V]) getOrCreateChildren() map[string]*TrieNode[V] {
	if n.children == nil {
		n.children = make(map[string]*TrieNode[V])
	}
	return n.children
}

func NewRouterTrie[V any](separator string, splitPathFuncs ...SplitPathFunc) ITrie[V] {
	splitPathFunc := defaultSplitPathFunc
	if len(splitPathFuncs) > 0 {
		splitPathFunc = splitPathFuncs[0]
	}
	return NewRouterTrieWith[V](separator, splitPathFunc, nil)
}

func NewRouterTrieWith[V any](separator string, splitPathFunc SplitPathFunc, splitMethodFunc SplitMethodFunc) ITrie[V] {
	pattern := fmt.Sprintf("[%s]", separator)
	regex, _ := regexp.Compile(pattern)
	if splitPathFunc == nil {
		splitPathFunc = defaultSplitPathFunc
	}
	if splitMethodFunc == nil {
		splitMethodFunc = splitCamelCaseFunc
	}
	return &Trie[V]{
		root:            &TrieNode[V]{},
		splitPathFunc:   splitPathFunc,
		splitMethodFunc: splitMethodFunc,
		separator:       separator,
		staticRouterMap: make(map[string]V),
		regex:           regex,
		rwLock:          sync.RWMutex{},
	}
}

func (t *Trie[V]) Contains(path string) bool {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	node := t.node(path)
	return node != nil
}

func (t *Trie[V]) Get(path string) (v V) {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	node := t.node(path)
	if node != nil {
		return node.value
	}
	return
}

func (t *Trie[V]) GetAndArgs(path string) (v V, args []string) {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	node := t.node(path)
	if node != nil {
		return node.value, node.args
	}
	return
}

func (t *Trie[V]) Add(path string, v V) (overwrittenValue V) {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	t.keyCheck(path)

	// if path is a static router
	if ok := t.regex.MatchString(path); !ok {
		if value, ok := t.staticRouterMap[path]; ok {
			overwrittenValue = value
		}
		t.staticRouterMap[path] = v
		t.size++
		return
	}

	// dynamic router
	currentNode := t.root
	for _, subPath := range t.splitMethodFunc(path, t.separator) {
		childNode, ok := currentNode.getOrCreateChildren()[subPath]
		if !ok {
			childNode = &TrieNode[V]{}
			currentNode.getOrCreateChildren()[subPath] = childNode
		}
		currentNode = childNode
	}
	if !currentNode.isEnd {
		currentNode.isEnd = true
		currentNode.value = v
		t.size++
		return
	}
	overwrittenValue = currentNode.value
	currentNode.value = v
	return
}

func (t *Trie[V]) Remove(path string) (oldValue V) {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	// static router
	if value, ok := t.staticRouterMap[path]; ok {
		oldValue = value
		delete(t.staticRouterMap, path)
		t.size--
		return
	}
	// dynamic router
	node := t.node(path)
	if node == nil {
		return
	}
	node.isEnd = false
	oldValue = node.value
	t.size--
	return
}

func (t *Trie[V]) StartWith(prefix string) bool {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	return t.node(prefix) != nil
}

func (t *Trie[V]) node(path string) *TrieNode[V] {
	t.keyCheck(path)
	// static route
	if nodeValue, ok := t.staticRouterMap[path]; ok {
		return &TrieNode[V]{value: nodeValue}
	}
	currentNode := t.root
	dynamicRoutingArgs := make([]string, 0)
	for _, subPath := range t.splitPathFunc(path) {
		childNode, ok := currentNode.children[subPath]
		if !ok {
			if v, ok := currentNode.children[t.separator]; ok {
				childNode = v
				// save args for dynamic routing
				dynamicRoutingArgs = append(dynamicRoutingArgs, subPath)
			} else {
				return nil
			}
		}
		currentNode = childNode
	}
	if currentNode.isEnd {
		// set args
		currentNode.args = dynamicRoutingArgs
		return currentNode
	}
	return nil
}

func (t *Trie[V]) Clear() {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	t.size = 0
	t.root.children = make(map[string]*TrieNode[V])
	t.staticRouterMap = make(map[string]V)
}

func (t *Trie[V]) IsEmpty() bool {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	return t.Size() == 0
}

func (t *Trie[V]) Size() int {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	return t.size
}

func (t *Trie[V]) keyCheck(k string) {
	if k == "" {
		panic("Trie[key is empty]")
	}
}
