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
	"slices"
)

func MapUpdate[K cmp.Ordered, V any](curr map[K]V, incom map[K]V) {
	for k, v := range incom {
		curr[k] = v
	}
}

type MapEntry[K cmp.Ordered, V any] struct {
	K K
	V V
}

func mapToEntries[K cmp.Ordered, V any](data map[K]V, sortBy func(a, b MapEntry[K, V]) int) []MapEntry[K, V] {
	ans := make([]MapEntry[K, V], len(data))
	i := 0
	for k, v := range data {
		ans[i] = MapEntry[K, V]{K: k, V: v}
		i++
	}
	if sortBy != nil {
		slices.SortFunc(ans, sortBy)
	}
	return ans
}

// MapToEntries transforms any map with keys cmp.Ordered to a slice of entries
func MapToEntries[K cmp.Ordered, V any](data map[K]V) []MapEntry[K, V] {
	return mapToEntries(data, nil)
}

// MapToEntries transforms any map with keys cmp.Ordered to a slice of entries,
// sorted by provided function (the rules for the function are the same as in slices.SortFunc)
func MapToEntriesSorted[K cmp.Ordered, V any](data map[K]V, sortBy func(a, b MapEntry[K, V]) int) []MapEntry[K, V] {
	return mapToEntries(data, sortBy)
}
