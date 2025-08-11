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
	"sort"
)

type Identifier interface {
	ID() string
}

// HSet is a set implementation for ordered value types
type HSet[T Identifier] struct {
	data map[string]T
}

func (set *HSet[T]) testAndInit() {
	if set.data == nil {
		set.data = make(map[string]T)
	}
}

func (set *HSet[T]) Add(value T) {
	set.testAndInit()
	set.data[value.ID()] = value
}

func (set *HSet[T]) Remove(value T) {
	set.testAndInit()
	delete(set.data, value.ID())
}

func (set *HSet[T]) Contains(value T) bool {
	set.testAndInit()
	_, ok := set.data[value.ID()]
	return ok
}

func (set *HSet[T]) ToSlice() []T {
	set.testAndInit()
	ans := make([]T, 0, len(set.data))
	for _, v := range set.data {
		ans = append(ans, v)
	}
	return ans
}

func (set *HSet[T]) ToOrderedSlice() []T {
	ans := set.ToSlice()
	sort.Slice(
		ans,
		func(i, j int) bool {
			return ans[i].ID() < ans[j].ID()
		},
	)
	return ans
}

func (set *HSet[T]) ForEach(fn func(item T)) {
	for _, v := range set.data {
		fn(v)
	}
}

// Iterate supports "range" form iteration through
// all the values of the set. The order of items
// is not guaranteed to be stable.
func (set *HSet[T]) Iterate(yield func(item T) bool) {
	for _, v := range set.data {
		if !yield(v) {
			return
		}
	}
}

func (set *HSet[T]) Union(other HSet[T]) *HSet[T] {
	set.testAndInit()
	ans := NewHSet(set.ToSlice()...)
	other.ForEach(func(item T) {
		ans.Add(item)
	})
	return ans
}

func (set *HSet[T]) Size() int {
	return len(set.data)
}

func (set *HSet[T]) Sub(other *HSet[T]) *HSet[T] {
	set.testAndInit()
	ans := NewHSet(set.ToSlice()...)
	other.ForEach(func(item T) {
		ans.Remove(item)
	})
	return ans
}

func (set *HSet[T]) Intersect(other *HSet[T]) *HSet[T] {
	set.testAndInit()
	ans := NewHSet([]T{}...)
	other.ForEach(func(item T) {
		if set.Contains(item) {
			ans.Add(item)
		}
	})
	return ans
}

func NewHSet[T Identifier](values ...T) *HSet[T] {
	ans := HSet[T]{data: make(map[string]T)}
	for _, v := range values {
		ans.data[v.ID()] = v
	}
	return &ans
}
