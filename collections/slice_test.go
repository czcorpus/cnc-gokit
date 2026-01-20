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
	"fmt"

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

type groupable struct {
	Type    string
	Enabled bool
	ID      int
}

func TestSliceGroupBy(t *testing.T) {
	key := func(v groupable) string {
		return fmt.Sprintf("%t-%s", v.Enabled, v.Type)
	}
	items := []groupable{
		{Type: "A", Enabled: true, ID: 1},
		{Type: "B", Enabled: true, ID: 2},
		{Type: "C", Enabled: true, ID: 3},
		{Type: "A", Enabled: false, ID: 4},
		{Type: "B", Enabled: false, ID: 5},
		{Type: "C", Enabled: true, ID: 6},
		{Type: "A", Enabled: true, ID: 7},
		{Type: "B", Enabled: true, ID: 8},
	}
	groups := SliceGroupBy(items, key)

	assert.Equal(t, 5, len(groups))

	var grAt []groupable // should be 1, 7
	var grAf []groupable // should be 4
	var grBt []groupable // should be 2, 8
	var grBf []groupable // should be 5
	var grCt []groupable // should be 3, 6
	for _, group := range groups {
		if group[0].ID == 1 && group[1].ID == 7 {
			grAt = group

		} else if group[0].ID == 4 {
			grAf = group

		} else if group[0].ID == 2 && group[1].ID == 8 {
			grBt = group

		} else if group[0].ID == 5 {
			grBf = group

		} else if group[0].ID == 3 && group[1].ID == 6 {
			grCt = group
		}
	}
	assert.Equal(t, 2, len(grAt))
	assert.Equal(t, 1, len(grAf))
	assert.Equal(t, 2, len(grBt))
	assert.Equal(t, 1, len(grBf))
	assert.Equal(t, 2, len(grCt))

}

func TestSliceGroupByEmpty(t *testing.T) {
	key := func(v groupable) string {
		return v.Type
	}
	items := []groupable{}
	groups := SliceGroupBy(items, key)
	assert.Equal(t, 0, len(groups))
}

func TestSliceGroupBySingleGroup(t *testing.T) {
	key := func(v int) string {
		return "all"
	}
	items := []int{1, 2, 3, 4, 5}
	groups := SliceGroupBy(items, key)
	assert.Equal(t, 1, len(groups))
	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, groups[0])
}
