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
	"testing"

	. "kubesphere.io/im/pkg/util/assert"
)

func TestSimplifyStringList(t *testing.T) {
	s0 := []string{"a", "", "c", "  ", " d "}
	s1 := SimplifyStringList(s0)

	Assert(t, len(s1) == 3)
	Assert(t, s1[0] == "a")
	Assert(t, s1[1] == "c")
	Assert(t, s1[2] == "d")
}

func TestSimplifyString(t *testing.T) {
	var tests = []struct{ s, expect string }{
		{s: "\ta  b  c", expect: "a b c"},
		{s: "a b c", expect: "a b c"},
		{s: "abc", expect: "abc"},
	}
	for _, v := range tests {
		got := SimplifyString(v.s)
		Assertf(t, got == v.expect, "expect = %q, got = %q", v.expect, got)
	}
}
