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

type Multidict[T any] struct {
	data map[string][]T
}

func (md *Multidict[T]) Add(k string, v T) {
	if md.data == nil {
		md.data = make(map[string][]T)
	}
	_, ok := md.data[k]
	if !ok {
		md.data[k] = make([]T, 0, 10)
	}
	md.data[k] = append(md.data[k], v)
}

func (md *Multidict[T]) Get(k string) []T {
	if md.data == nil {
		md.data = make(map[string][]T)
	}
	return md.data[k]
}

func (md *Multidict[T]) Iterate(yield func(k string, v []T) bool) {
	if md.data == nil {
		md.data = make(map[string][]T)
	}
	for k, v := range md.data {
		if !yield(k, v) {
			return
		}
	}
}

func (md *Multidict[T]) IterateFlat(yield func(k string, v T) bool) {
	if md.data == nil {
		md.data = make(map[string][]T)
	}
	for k, v := range md.data {
		for _, v2 := range v {
			if !yield(k, v2) {
				return
			}
		}
	}
}

func NewMultidict[T any]() *Multidict[T] {
	return &Multidict[T]{
		data: make(map[string][]T),
	}
}
