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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMapSetGet(t *testing.T) {
	c := NewConcurrentMap[string, int]()
	c.Set("foo", 100)
	assert.Equal(t, 100, c.Get("foo"))
}

func TestConcurrentMapGetWithTest(t *testing.T) {
	c := NewConcurrentMap[string, int]()
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

func TestConcurrentMapLenZero(t *testing.T) {
	c := NewConcurrentMap[string, int]()
	assert.Equal(t, 0, c.Len())
}

func TestConcurrentMapLen(t *testing.T) {
	c := NewConcurrentMap[string, int]()
	c.Set("foo", 100)
	assert.Equal(t, 1, c.Len())
	c.Set("bar", 101)
	assert.Equal(t, 2, c.Len())
}

func TestConcurrentMapDelete(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	c := NewConcurrentMapFrom(data)
	v, ok := c.GetWithTest("bar")
	assert.Equal(t, 2, v)
	assert.True(t, ok)
	c.Delete("bar")
	v, ok = c.GetWithTest("bar")
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestNewConcurrentMapFromJSON(t *testing.T) {
	src := `{"foo": 10, "bar": 20, "baz": 30}`
	v, err := NewConcurrentMapFromJSON[string, int]([]byte(src))
	assert.NoError(t, err)
	assert.Equal(t, 10, v.Get("foo"))
	assert.Equal(t, 20, v.Get("bar"))
	assert.Equal(t, 30, v.Get("baz"))
	assert.Equal(t, 3, v.Len())
}

func TestConcurrentMapJSONSerialization(t *testing.T) {
	v := NewConcurrentMap[string, int]()
	v.Set("foo", 10)
	v.Set("bar", 20)
	src, err := json.Marshal(v)
	assert.NoError(t, err)
	src2 := string(src)
	assert.True(t, strings.Contains(src2, `"bar":20`))
	assert.True(t, strings.Contains(src2, `"foo":10`))
}

func TestConcurrentMapFilter(t *testing.T) {
	v := NewConcurrentMapFrom[string, int](map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
		"faz": 4,
		"fuz": 5,
	})
	v = v.Filter(func(k string, v int) bool {
		return k[0] == 'f'
	})
	assert.Equal(t, 1, v.Get("foo"))
	assert.Equal(t, 4, v.Get("faz"))
	assert.Equal(t, 5, v.Get("fuz"))
	assert.Equal(t, 3, v.Len())
}
