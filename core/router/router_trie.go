// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"regexp"
	"strings"
)

type ITrie[V any] interface {
	Size() int // total node
	IsEmpty() bool
	Clear()
	Contains(path string) bool
	Get(path string) V
	Add(path string, v V) (overwrittenValue V)
	Remove(path string) (oldValue V)
	StartWith(prefix string) bool
}

func NewRouterTrie[V any](separator string, f ...SplitPathFunc) ITrie[V] {
	pattern := `[*:]`
	regex, _ := regexp.Compile(pattern)
	splitPathFunc := defaultSplitPathFunc
	if len(f) > 0 {
		splitPathFunc = f[0]
	}
	return &Trie[V]{
		root:            &TrieNode[V]{},
		splitFunc:       splitPathFunc,
		separator:       separator,
		staticRouterMap: map[string]V{},
		regex:           regex,
	}
}

type SplitPathFunc = func(k string) []string

var defaultSplitPathFunc = func(path string) []string {
	return strings.Split(path, "/")
}

type Trie[V any] struct {
	size      int           // size of node
	root      *TrieNode[V]  // root of trie
	splitFunc SplitPathFunc // split func of full path
	separator string        //

	// Optimized fields
	staticRouterMap map[string]V
	regex           *regexp.Regexp
}

type TrieNode[V any] struct {
	children map[string]*TrieNode[V] // k is the current sub-path
	value    V                       // value
	isEnd    bool                    // judge is end
}

func (n *TrieNode[V]) getOrCreateChildren() map[string]*TrieNode[V] {
	if n.children == nil {
		n.children = make(map[string]*TrieNode[V])
	}
	return n.children
}

func (t *Trie[V]) Contains(path string) bool {
	node := t.node(path)
	return node != nil
}

func (t *Trie[V]) Get(path string) (v V) {
	node := t.node(path)
	if node != nil {
		return node.value
	}
	return
}

func (t *Trie[V]) Add(path string, v V) (overwrittenValue V) {
	// if path is static router
	if ok := t.regex.MatchString(path); !ok {
		if value, ok := t.staticRouterMap[path]; ok {
			overwrittenValue = value
		}
		t.staticRouterMap[path] = v
		t.size++
		return
	}
	t.keyCheck(path)
	currentNode := t.root
	for _, subPath := range t.splitFunc(path)[1:] {
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
	return t.node(prefix) != nil
}

func (t *Trie[V]) node(path string) *TrieNode[V] {
	t.keyCheck(path)
	// static route
	if nodeValue, ok := t.staticRouterMap[path]; ok {
		return &TrieNode[V]{value: nodeValue}
	}
	currentNode := t.root
	for _, subPath := range t.splitFunc(path)[1:] {
		childNode, ok := currentNode.children[subPath]
		if !ok {
			if v, ok := currentNode.children[t.separator]; ok {
				childNode = v
			} else {
				return nil
			}
		}
		currentNode = childNode
	}
	if currentNode.isEnd {
		return currentNode
	}
	return nil
}

func (t *Trie[V]) Clear() {
	t.size = 0
	t.root.children = make(map[string]*TrieNode[V])
}

func (t *Trie[V]) IsEmpty() bool {
	return t.Size() == 0
}

func (t *Trie[V]) Size() int {
	return t.size
}

func (t *Trie[V]) keyCheck(k string) {
	if k == "" {
		panic("Trie[key is empty]")
	}
}
