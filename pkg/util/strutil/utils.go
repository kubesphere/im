/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package strutil

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func NewString(v string) *string {
	return &v
}

func SimplifyStringList(s []string) []string {
	b := s[:0]
	for _, x := range s {
		if x := SimplifyString(x); x != "" {
			b = append(b, x)
		}
	}
	return b
}

var reMoreSpace = regexp.MustCompile(`\s+`)

// "\ta  b  c" => "a b c"
func SimplifyString(s string) string {
	return reMoreSpace.ReplaceAllString(strings.TrimSpace(s), " ")
}

func Contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}
