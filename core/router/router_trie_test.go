package router

import (
	"strconv"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewRouterTrie[int]("*")
	trie.Add("/users", 1)
	trie.Add("/users/*", 2)

	// Test Contains
	if !trie.Contains("/users") {
		t.Error("Expected /users to be in the trie")
	}

	// Test Get
	value := trie.Get("/users/:id")
	if value != 2 {
		t.Errorf("Expected value 2, got %d", value)
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
	if size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// Test Clear and IsEmpty
	trie.Clear()
	if !trie.IsEmpty() {
		t.Error("Expected trie to be empty after clearing")
	}
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
