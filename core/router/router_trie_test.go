// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewRouterTrie[int]("*")
	trie.Add("/users", 1)
	trie.Add("/users/*", 2)
	trie.Add("/users/*/alice", 3)
	trie.Add("/users/*/alice/*/bob", 4)

	// Test Contains
	if !trie.Contains("/users") {
		t.Error("Expected /users to be in the trie")
	}

	// Test Get
	value := trie.Get("/users/111")
	if value != 2 {
		t.Errorf("Expected value 2, got %d", value)
	}

	// Test Get
	value = trie.Get("/users/111/alice")
	if value != 3 {
		t.Errorf("Expected value 3, got %d", value)
	}

	value, args := trie.GetAndArgs("/users/111/alice")

	if value != 3 {
		t.Errorf("Expected value 3, got %d", value)
	}

	parseInt, _ := strconv.ParseInt(args[0], 10, 32)

	if len(args) != 1 || parseInt != 111 {
		t.Errorf("Expected args 111, got %v", args)
	}

	value, args = trie.GetAndArgs("/users/111/alice/abc/bob")

	if value != 4 {
		t.Errorf("Expected value 3, got %v", value)
	}

	parseInt, _ = strconv.ParseInt(args[0], 10, 32)

	if len(args) != 2 || parseInt != 111 || args[1] != "abc" {
		t.Errorf("Expected args abc, got %v", args)
	}

	// Test StartWith
	if !trie.StartWith("/users") {
		t.Error("Expected /users to be a prefix")
	}

	// Test Remove
	removedValue := trie.Remove("/users")
	if removedValue != 1 {
		t.Errorf("Expected removed value 1, got %d", removedValue)
	}
	if trie.Contains("/users") {
		t.Error("Expected /users to be removed from the trie")
	}

	// Test Size
	size := trie.Size()
	if size != 3 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// Test Clear and IsEmpty
	trie.Clear()
	if !trie.IsEmpty() {
		t.Error("Expected trie to be empty after clearing")
	}
}

func TestTrieConcurrentAccess(t *testing.T) {
	trie := NewRouterTrie[string]("/")
	wg := sync.WaitGroup{}
	numRoutines := 100
	numOperations := 1000
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				path := fmt.Sprintf("/path-%d-%d", i, j)
				trie.Add(path, path)
				value := trie.Get(path)
				if value != path {
					t.Errorf("Incorrect value for path %s: expected %s, got %s", path, path, value)
				}
				trie.Remove(path)
				value = trie.Get(path)
				if value != "" {
					t.Errorf("Remove failed for path %s: value still exists: %s", path, value)
				}
			}
		}()
	}
	wg.Wait()
}

func TestTrieConcurrentContains(t *testing.T) {
	trie := NewRouterTrie[string]("/")
	numPaths := 100
	for i := 0; i < numPaths; i++ {
		path := fmt.Sprintf("/path-%d", i)
		trie.Add(path, path)
	}
	wg := sync.WaitGroup{}
	numRoutines := 100
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numPaths; j++ {
				path := fmt.Sprintf("/path-%d", j)
				contains := trie.Contains(path)
				if !contains {
					t.Errorf("Contains failed for path %s", path)
				}
			}
		}()
	}
	wg.Wait()
}

func TestTrieConcurrentStartWith(t *testing.T) {
	trie := NewRouterTrie[string]("/")
	numPaths := 100
	for i := 0; i < numPaths; i++ {
		path := fmt.Sprintf("/path-%d", i)
		trie.Add(path, path)
	}
	wg := sync.WaitGroup{}
	numRoutines := 100
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numPaths; j++ {
				prefix := fmt.Sprintf("/path-%d", j)
				startsWith := trie.StartWith(prefix)
				if !startsWith {
					t.Errorf("StartWith failed for prefix %s", prefix)
				}
			}
		}()
	}
	wg.Wait()
}

/* ------------------------------------------------------
	$ go test -bench=. -benchmem
	goos: windows
	goarch: amd64
	pkg: jet-web/core/router
	cpu: 13th Gen Intel(R) Core(TM) i5-13400F
	BenchmarkRouterTrie_Add-16               2824297               425.9 ns/op           113 B/op          2 allocs/op
	BenchmarkRouterTrie_Get-16              14866627                77.90 ns/op           39 B/op          2 allocs/op
	BenchmarkRouterTrie_Remove-16           13333392                84.92 ns/op           63 B/op          2 allocs/op
	PASS
	ok      jet-web/core/router     4.301s
------------------------------------------------------ */

func BenchmarkRouterTrie_Add(b *testing.B) {
	trie := NewRouterTrie[int]("/")
	for i := 0; i < b.N; i++ {
		path := "/users/" + strconv.Itoa(i)
		trie.Add(path, i)
	}
}

func BenchmarkRouterTrie_Get(b *testing.B) {
	trie := NewRouterTrie[int]("/")
	for i := 0; i < 10000; i++ {
		path := "/users/" + strconv.Itoa(i)
		trie.Add(path, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := "/users/" + strconv.Itoa(i%10000)
		trie.Get(path)
	}
}

func BenchmarkRouterTrie_Remove(b *testing.B) {
	trie := NewRouterTrie[int]("/")
	for i := 0; i < 10000; i++ {
		path := "/users/" + strconv.Itoa(i)
		trie.Add(path, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := "/users/" + strconv.Itoa(i%10000)
		trie.Remove(path)
	}
}
