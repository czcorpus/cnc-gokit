// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMapSetGet(t *testing.T) {
	c := NewConcurrentMap[int]()
	c.Set("foo", 100)
	assert.Equal(t, 100, c.Get("foo"))
}

func TestConcurrentMapGetWithTest(t *testing.T) {
	c := NewConcurrentMap[int]()
	c.Set("foo", 100)
	v, ok := c.GetWithTest("foo")
	assert.Equal(t, 100, v)
	assert.True(t, ok)
	v, ok = c.GetWithTest("bar")
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestConcurrentMapKeys(t *testing.T) {

	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	c := NewConcurrentMapFrom(data)
	keys := c.Keys()
	assert.Subset(t, []string{"foo", "bar", "baz"}, keys)
	assert.Equal(t, 3, len(keys))
}

func TestConcurrentMapValues(t *testing.T) {

	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	c := NewConcurrentMapFrom(data)
	values := c.Values()
	assert.Subset(t, []int{1, 2, 3}, values)
	assert.Equal(t, 3, len(values))
}

func TestConcurrentMapUpdate(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	c := NewConcurrentMapFrom(data)
	c.Update(func(k string, v int) int {
		return v + 1
	})

	for i, v := range []string{"foo", "bar", "baz"} {
		assert.Equal(t, i+2, c.Get(v))
	}
}
