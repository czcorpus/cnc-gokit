// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Institute of the Czech National Corpus,
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

type someStruct struct {
	Name string
}

type otherStruct struct {
	ID string
}

type numValue struct {
	V int
}

func TestSliceContains(t *testing.T) {
	s := []string{"foo", "bar", "baz", "foo"}
	ans := SliceContains(s, "baz")
	assert.True(t, ans)
}

func TestSliceContainsOnEmpty(t *testing.T) {
	s := []string{}
	ans := SliceContains(s, "baz")
	assert.False(t, ans)
}

func TestSliceContainsRefRespect(t *testing.T) {
	v1 := someStruct{"foo"}
	v2 := someStruct{"bar"}
	v3 := someStruct{"baz"}
	v4 := someStruct{"foo"}
	s := []*someStruct{&v1, &v2, &v3}
	ans := SliceContains(s, &v4)
	assert.False(t, ans)
}

func TestSliceFindIndex(t *testing.T) {
	s := []string{"bar", "foo", "baz", "foo"}
	ans := SliceFindIndex(s, func(v string) bool { return v == "foo" })
	assert.Equal(t, 1, ans)
}

func TestSliceFindIndexOnEmpty(t *testing.T) {
	s := []string{}
	ans := SliceFindIndex(s, func(v string) bool { return v == "foo" })
	assert.Equal(t, -1, ans)
}

func TestSliceFindIndexNotFound(t *testing.T) {
	s := []string{"foo", "bar", "baz", "foo"}
	ans := SliceFindIndex(s, func(v string) bool { return v == "xxx" })
	assert.Equal(t, -1, ans)
}

func TestSliceMap(t *testing.T) {
	s := []someStruct{{"foo"}, {"bar"}, {"baz"}, {"foo"}}
	ans := SliceMap(s, func(v someStruct, i int) otherStruct { return otherStruct{ID: v.Name} })
	exp := []otherStruct{{"foo"}, {"bar"}, {"baz"}, {"foo"}}
	assert.Equal(t, exp, ans)
}

func TestSliceReduce(t *testing.T) {
	s := []numValue{{V: 1}, {V: 2}, {V: 3}, {V: 4}}
	ans := SliceReduce(
		s,
		func(acc int, curr numValue, i int) int {
			return acc + curr.V
		},
		70,
	)
	assert.Equal(t, 80, ans)
}
