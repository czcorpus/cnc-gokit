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
	"cmp"
	"sort"
)

// Set is a set implementation for ordered value types
type Set[T cmp.Ordered] struct {
	data map[T]bool
}

func (set *Set[T]) testAndInit() {
	if set.data == nil {
		set.data = make(map[T]bool)
	}
}

func (set *Set[T]) Add(value T) {
	set.testAndInit()
	set.data[value] = true
}

func (set *Set[T]) Remove(value T) {
	set.testAndInit()
	delete(set.data, value)
}

func (set *Set[T]) Contains(value T) bool {
	set.testAndInit()
	_, ok := set.data[value]
	return ok
}

func (set *Set[T]) ToSlice() []T {
	set.testAndInit()
	ans := make([]T, 0, len(set.data))
	for k := range set.data {
		ans = append(ans, k)
	}
	return ans
}

func (set *Set[T]) ToOrderedSlice() []T {
	ans := set.ToSlice()
	sort.Slice(
		ans,
		func(i, j int) bool {
			return ans[i] < ans[j]
		},
	)
	return ans
}

func (set *Set[T]) ForEach(fn func(item T)) {
	for k := range set.data {
		fn(k)
	}
}

func (set *Set[T]) Union(other Set[T]) *Set[T] {
	set.testAndInit()
	ans := NewSet(set.ToSlice()...)
	other.ForEach(func(item T) {
		ans.Add(item)
	})
	return ans
}

func (set *Set[T]) Size() int {
	return len(set.data)
}

func (set *Set[T]) Sub(other *Set[T]) *Set[T] {
	set.testAndInit()
	ans := NewSet(set.ToSlice()...)
	other.ForEach(func(item T) {
		ans.Remove(item)
	})
	return ans
}

func (set *Set[T]) Intersect(other *Set[T]) *Set[T] {
	set.testAndInit()
	ans := NewSet([]T{}...)
	other.ForEach(func(item T) {
		if set.Contains(item) {
			ans.Add(item)
		}
	})
	return ans
}

func NewSet[T cmp.Ordered](values ...T) *Set[T] {
	ans := Set[T]{data: make(map[T]bool)}
	for _, v := range values {
		ans.data[v] = true
	}
	return &ans
}
