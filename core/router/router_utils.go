// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"bytes"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/valyala/fasthttp"
	"strings"
	"unicode"
)

// if sep is [_]
// AppleBanana => ["apple", "banana"]
// Apple_Banana => ["apple", "*", "banana"]
// AppleBanana_ => ["apple", "banana", "*"]
// Apple_Banana_ => ["apple", "*", "banana", "*"]
// ...
func splitCamelCaseFunc(method string, sep string) (pattern []string) {
	defer func() {
		if len(pattern) != 0 {
			for index, sub := range pattern {
				if sub == sep {
					continue
				}
				pattern[index] = strings.ToLower(sub[:1]) + sub[1:]
			}
		}
	}()
	for method != "" {
		pos := strings.Index(method, sep)
		if pos == -1 {
			return appendPattern(pattern, method)
		}
		if pos > 0 {
			pattern = appendPattern(pattern, method[:pos])
		}
		pattern = append(pattern, sep)
		method = method[pos+len(sep):]
	}
	return
}

func appendPattern(pattern []string, method string) []string {
	var i, last int
	for i = 1; i < len(method); i++ {
		c := method[i]
		if c >= 'A' && c <= 'Z' {
			pattern = append(pattern, method[last:i])
			last = i
		}
	}
	return append(pattern, method[last:i])
}

func convertToCamelCase(path []byte) string {
	components := strings.Split(string(path), "/")

	var result string
	for _, component := range components {
		if component != "" {
			result += strings.Title(component)
		}
	}

	return result
}

func convertToFirstLetterUpper(method []byte) string {
	return strings.ToLower(string(method))
}

// ConvertToURL Convert uppercase-separated string to URL format
// like GetV1UsageWeek => get/v1/usage/week
func ConvertToURL(input string) string {
	var buffer bytes.Buffer

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				buffer.WriteRune('/')
			}
			buffer.WriteString(strings.ToLower(string(r)))
		} else {
			buffer.WriteRune(unicode.ToLower(r))
		}
	}

	return buffer.String()
}

func FirstRuneIsUp(input string) bool {
	return len(input) > 0 && unicode.IsUpper(rune(input[0]))
}

func ParseUrl(ctx *fasthttp.RequestCtx) (queryStringMap maps.IMap[string, string]) {
	queryStringMap = maps.NewLinkedHashMap[string, string]()
	// Parse the query parameters
	ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryStringMap.Put(string(key), string(value))
	})
	return
}
