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

func TestSliceFilter(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	ans := SliceFilter(s, func(v int, i int) bool { return v > 4 })
	assert.Equal(t, []int{5, 6}, ans)
}

func TestSliceFilterRetEmpty(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	ans := SliceFilter(s, func(v int, i int) bool { return v > 6 })
	assert.Equal(t, []int{}, ans)
}

// TestSliceShuffle
// note: Due to randomness in the function, we
// just test whether any position of the tested
// slice does not exhibit too different behavior
func TestSliceShuffle(t *testing.T) {
	results := make([][]int, 0, 10000)

	genShuffle := func() []int {
		s := []int{0, 1, 2, 3, 4, 5, 6, 7}
		SliceShuffle(s)
		return s
	}

	for i := 0; i < 10000; i++ {
		results = append(results, genShuffle())
	}

	sums := make(map[int]int)
	for i := 0; i < 10000; i++ {
		for j := 0; j < 8; j++ {
			sums[j] += results[i][j]
		}
	}
	var avg int
	for _, v := range sums {
		avg += v
	}
	avg /= len(sums)
	for _, v := range sums {
		assert.InDelta(t, avg, v, 1000)
	}
}

func TestRandomSample(t *testing.T) {
	rnd := &mockrnd{sequence: []int{3, 2, 1, 1, 4, 0, 1, 3, 2, 1}}
	ans := sliceSample([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 5, rnd)
	// -) [3]  0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	// 0) [2]  0, 1, 2, 9, 4, 5, 6, 7, 8, 3
	// 1) [1]  0, 1, 8, 9, 4, 5, 6, 7, 2, 3
	// 2) [1]  0, 7, 8, 9, 4, 5, 6, 1, 2, 3
	// 3) [4]  0, 6, 8, 9, 4, 5, 7, 1, 2, 3
	// 4) -    0, 6, 8, 9, 5, 4, 7, 1, 2, 3
	assert.Equal(t, []int{4, 7, 1, 2, 3}, ans)
}

func TestRandomSampleZeroSample(t *testing.T) {
	rnd := &mockrnd{sequence: []int{3, 2, 1, 1, 4, 0, 1, 3, 2, 1}}
	ans := sliceSample([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 0, rnd)
	assert.Equal(t, []int{}, ans)
}

func TestRandomSampleTooBigSample(t *testing.T) {
	rnd := &mockrnd{sequence: []int{3, 2, 1, 1, 4, 0, 1, 3, 2, 1}}
	assert.Panics(t, func() {
		sliceSample([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 20, rnd)
	})
}

func TestRandomSampleMaxSample(t *testing.T) {
	rnd := &mockrnd{sequence: []int{3, 2, 1, 1, 4, 0, 1, 3, 2, 1}}
	ans := sliceSample([]int{0, 1, 2, 3}, 4, rnd)
	// -) [3]  0, 1, 2, 3
	// 0) [2]  0, 1, 2, 3
	// 1) [1]  0, 1, 2, 3
	// 2) [1]  0, 1, 2, 3
	// 3)      0, 1, 2, 3
	assert.Equal(t, []int{0, 1, 2, 3}, ans)
}
