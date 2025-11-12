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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapUpdate(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"b": 10, "c": 100, "d": 1000}
	MapUpdate(m1, m2)
	assert.Equal(t, map[string]int{"a": 1, "b": 10, "c": 100, "d": 1000}, m1)
	assert.Equal(t, map[string]int{"b": 10, "c": 100, "d": 1000}, m2)
}

func TestMapUpdateByEmpty(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{}
	MapUpdate(m1, m2)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m1)
	assert.Equal(t, map[string]int{}, m2)
}

func TestMapUpdateEmpty(t *testing.T) {
	m1 := map[string]int{}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	MapUpdate(m1, m2)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m1)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m2)
}

func TestMapToEntries(t *testing.T) {
	m1 := map[string]int{"foo": 10, "bar": 20, "baz": 30}
	items := MapToEntries(m1)
	assert.Equal(t, 3, len(items))
	hits := 0
	for _, entry := range items {
		fmt.Println("ENTRY: ", entry)
		if entry.K == "foo" && entry.V == 10 {
			hits++
		}
		if entry.K == "bar" && entry.V == 20 {
			hits++
		}
		if entry.K == "baz" && entry.V == 30 {
			hits++
		}
	}
	assert.Equal(t, 3, hits)
}

func TestMapToEntriesEmpty(t *testing.T) {
	m1 := map[string]int{}
	items := MapToEntries(m1)
	assert.Equal(t, 0, len(items))
}

func TestMapToEntriesSorted(t *testing.T) {
	m1 := map[string]int{"foo": 10, "bar": 20, "baz": 30}
	items := MapToEntriesSorted(m1, func(a, b MapEntry[string, int]) int {
		if a.V > b.V {
			return -1
		}
		if a.V < b.V {
			return 1
		}
		return 0
	})
	assert.Equal(t, 3, len(items))
	assert.Equal(t, "baz", items[0].K)
	assert.Equal(t, "bar", items[1].K)
	assert.Equal(t, "foo", items[2].K)
}
