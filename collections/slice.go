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
	"math/rand"
	"time"
)

func SliceContains[T comparable](data []T, value T) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func SliceFindIndex[T any](data []T, cmp func(v T) bool) int {
	for i, v := range data {
		if cmp(v) {
			return i
		}
	}
	return -1
}

func SliceMap[T any, U any](data []T, mapFn func(v T, i int) U) []U {
	ans := make([]U, len(data))
	for i, item := range data {
		ans[i] = mapFn(item, i)
	}
	return ans
}

func SliceReduce[T any, U any](data []T, reduceFn func(acc U, curr T, i int) U, initial U) U {
	ans := initial
	for i, item := range data {
		ans = reduceFn(ans, item, i)
	}
	return ans
}

func SliceFilter[T any](data []T, filterFn func(v T, i int) bool) []T {
	ans := make([]T, 0, len(data))
	for i, v := range data {
		if filterFn(v, i) {
			ans = append(ans, v)
		}
	}
	return ans
}

// SliceShuffle shuffles a slice in place
func SliceShuffle[T any](data []T) {
	for i := 0; i < len(data); i++ {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

type randomSource interface {
	Intn(n int) int
}

func sliceSample[T any](data []T, sampleSize int, rnd randomSource) []T {
	if sampleSize > len(data) {
		panic("SliceSample - the sampleSize must be at most the length of the original data")
	}
	tmp := make([]T, len(data))
	copy(tmp, data)
	for i := 0; i < sampleSize; i++ {
		j := rnd.Intn(len(tmp) - i)
		tmp[len(tmp)-1-i], tmp[j] = tmp[j], tmp[len(tmp)-1-i]
	}
	return tmp[len(tmp)-sampleSize:]
}

// SliceSample creates a uniform sample of size given
// by the sampleSize argument. The function allocates
// a copy of the input data so it should be taken into
// account when dealing with large slices.
// Please note that the randomness used by the
// function is not cryptographically secure.
// A zero-size sample is accepted.
func SliceSample[T any](data []T, sampleSize int) []T {
	return sliceSample(data, sampleSize, rand.New(rand.NewSource(time.Now().Unix())))
}

// SliceGroupBy takse a slice and groups its items based on
// how function `key` associates string values to each individual
// item.
// The order of groups is not guaranteed (it comes from how
// internally used map works).
func SliceGroupBy[T any](items []T, key func(v T) string) [][]T {
	tmp := make(map[string][]T)
	numGroups := 0
	for _, item := range items {
		k := key(item)
		curr, ok := tmp[k]
		if !ok {
			numGroups++
			curr = make([]T, 0, len(items)/4)
		}
		curr = append(curr, item)
		tmp[k] = curr
	}
	ans := make([][]T, numGroups)
	i := 0
	for _, v := range tmp {
		ans[i] = v
		i++
	}
	return ans
}
