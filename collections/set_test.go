// Copyright 2022 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2022 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
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

func TestSetSimpleInit(t *testing.T) {
	s := Set[string]{}
	assert.Equal(t, 0, s.Size())

	assert.NotPanics(t, func() {
		s.Remove("x")
	})

	var i int
	assert.NotPanics(t, func() {
		s.ForEach(func(item string) {
			i++
		})
	})
	assert.Equal(t, 0, i)

	assert.NotPanics(t, func() {
		s.Add("foo")
	})
}

func TestSetUniqueValues(t *testing.T) {
	s := Set[string]{}
	s.Add("one")
	s.Add("one")
	assert.Equal(t, 1, s.Size())

	s.Remove("one")
	assert.Equal(t, 0, s.Size())
}

func TestSetContains(t *testing.T) {
	s := Set[string]{}
	s.Add("one")
	s.Add("two")
	assert.False(t, s.Contains("three"))
	assert.True(t, s.Contains("one"))
	assert.True(t, s.Contains("two"))
}

func TestSetSize(t *testing.T) {
	s := Set[string]{}
	assert.Equal(t, 0, s.Size())
	s.Add("foo")
	assert.Equal(t, 1, s.Size())
}

func TestSetForEach(t *testing.T) {
	s := Set[string]{}
	s.Add("1")
	s.Add("2")
	s.Add("3")
	tst := make(map[string]bool)
	s.ForEach(func(item string) {
		tst[item] = true
	})
	assert.True(t, tst["1"])
	assert.True(t, tst["2"])
	assert.True(t, tst["3"])
	assert.Equal(t, 3, len(tst))
}
