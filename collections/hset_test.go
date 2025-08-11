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

type testItem struct {
	id   string
	data string
}

func (t testItem) ID() string {
	return t.id
}

func TestHSetSimpleInit(t *testing.T) {
	s := HSet[testItem]{}
	assert.Equal(t, 0, s.Size())

	assert.NotPanics(t, func() {
		s.Remove(testItem{id: "x"})
	})

	var i int
	assert.NotPanics(t, func() {
		s.ForEach(func(item testItem) {
			i++
		})
	})
	assert.Equal(t, 0, i)

	assert.NotPanics(t, func() {
		s.Add(testItem{id: "foo"})
	})
}

func TestHSetUniqueValues(t *testing.T) {
	s := HSet[testItem]{}
	s.Add(testItem{id: "one", data: "data1"})
	s.Add(testItem{id: "one", data: "data2"})
	assert.Equal(t, 1, s.Size())

	s.Remove(testItem{id: "one"})
	assert.Equal(t, 0, s.Size())
}

func TestHSetContains(t *testing.T) {
	s := HSet[testItem]{}
	s.Add(testItem{id: "one", data: "data1"})
	s.Add(testItem{id: "two", data: "data2"})
	assert.False(t, s.Contains(testItem{id: "three"}))
	assert.True(t, s.Contains(testItem{id: "one"}))
	assert.True(t, s.Contains(testItem{id: "two"}))
}

func TestHSetSize(t *testing.T) {
	s := HSet[testItem]{}
	assert.Equal(t, 0, s.Size())
	s.Add(testItem{id: "foo", data: "bar"})
	assert.Equal(t, 1, s.Size())
}

func TestHSetForEach(t *testing.T) {
	s := HSet[testItem]{}
	s.Add(testItem{id: "1", data: "a"})
	s.Add(testItem{id: "2", data: "b"})
	s.Add(testItem{id: "3", data: "c"})
	tst := make(map[string]bool)
	s.ForEach(func(item testItem) {
		tst[item.ID()] = true
	})
	assert.True(t, tst["1"])
	assert.True(t, tst["2"])
	assert.True(t, tst["3"])
	assert.Equal(t, 3, len(tst))
}

func TestHSetSubFromEmptySet(t *testing.T) {
	s0 := HSet[testItem]{}
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s01 := s0.Sub(&s1)
	assert.Equal(t, 0, s01.Size())
}

func TestHSetSubOfEmptySet(t *testing.T) {
	s0 := HSet[testItem]{}
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s10 := s1.Sub(&s0)
	assert.True(t, s10.Contains(testItem{id: "1"}))
	assert.True(t, s10.Contains(testItem{id: "2"}))
	assert.True(t, s10.Contains(testItem{id: "3"}))
	assert.Equal(t, 3, s10.Size())
}

func TestHSetSub(t *testing.T) {
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s2 := HSet[testItem]{}
	s2.Add(testItem{id: "2", data: "x"})
	s2.Add(testItem{id: "3", data: "y"})
	s2.Add(testItem{id: "4", data: "z"})
	s12 := s1.Sub(&s2)
	assert.True(t, s12.Contains(testItem{id: "1"}))
	assert.False(t, s12.Contains(testItem{id: "2"}))
	assert.False(t, s12.Contains(testItem{id: "3"}))
	assert.Equal(t, 1, s12.Size())
}

func TestHSetIntersectOfEmptySet(t *testing.T) {
	s0 := HSet[testItem]{}
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s10 := s1.Intersect(&s0)
	assert.Equal(t, 0, s10.Size())
}

func TestHSetIntersectFromEmptySet(t *testing.T) {
	s0 := HSet[testItem]{}
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s01 := s0.Intersect(&s1)
	assert.Equal(t, 0, s01.Size())
}

func TestHSetIntersect(t *testing.T) {
	s1 := HSet[testItem]{}
	s1.Add(testItem{id: "1", data: "a"})
	s1.Add(testItem{id: "2", data: "b"})
	s1.Add(testItem{id: "3", data: "c"})
	s2 := HSet[testItem]{}
	s2.Add(testItem{id: "2", data: "x"})
	s2.Add(testItem{id: "3", data: "y"})
	s2.Add(testItem{id: "4", data: "z"})
	s12 := s1.Intersect(&s2)
	assert.False(t, s12.Contains(testItem{id: "1"}))
	assert.True(t, s12.Contains(testItem{id: "2"}))
	assert.True(t, s12.Contains(testItem{id: "3"}))
	assert.Equal(t, 2, s12.Size())
}