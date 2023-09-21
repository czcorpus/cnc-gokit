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
