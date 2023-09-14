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
