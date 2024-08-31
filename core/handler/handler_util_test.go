package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrefixOf(t *testing.T) {
	// 测试数据
	testData := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"Get", "", false},
		{"GetUser", "Get", true},
		{"GetUserList", "Get", true},
	}

	for _, td := range testData {
		prefix, ok := prefixOf(td.input)
		assert.Equal(t, td.expected, prefix, "For input '%s', expected prefix '%s', got '%s'", td.input, td.expected, prefix)
		assert.Equal(t, td.ok, ok, "For input '%s', expected ok '%t', got '%t'", td.input, td.ok, ok)
	}
}
