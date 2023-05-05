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

import "sync"

type ConcurrentMap[T any] struct {
	sync.RWMutex
	data map[string]T
}

func (cm *ConcurrentMap[T]) Get(k string) T {
	cm.RLock()
	defer cm.RUnlock()
	return cm.data[k]
}

func (cm *ConcurrentMap[T]) GetWithTest(k string) (T, bool) {
	cm.RLock()
	defer cm.RUnlock()
	v, ok := cm.data[k]
	return v, ok
}

func (cm *ConcurrentMap[T]) HasKey(k string) bool {
	cm.RLock()
	defer cm.RUnlock()
	_, ok := cm.data[k]
	return ok
}

func (cm *ConcurrentMap[T]) Set(k string, v T) {
	cm.Lock()
	defer cm.Unlock()
	cm.data[k] = v
}

func (cm *ConcurrentMap[T]) ForEach(fn func(k string, v T)) {
	cm.RLock()
	defer cm.RUnlock()
	for k, v := range cm.data {
		fn(k, v)
	}
}

func (cm *ConcurrentMap[T]) Update(fn func(k string, v T) T) {
	cm.Lock()
	defer cm.Unlock()
	for k, v := range cm.data {
		cm.data[k] = fn(k, v)
	}
}

func (cm *ConcurrentMap[T]) Keys() []string {
	ans := make([]string, len(cm.data))
	i := 0
	for k, _ := range cm.data {
		ans[i] = k
		i++
	}
	return ans
}

func (cm *ConcurrentMap[T]) Values() []T {
	ans := make([]T, len(cm.data))
	i := 0
	for _, v := range cm.data {
		ans[i] = v
		i++
	}
	return ans
}

func NewConcurrentMap[T any]() *ConcurrentMap[T] {
	return &ConcurrentMap[T]{
		data: make(map[string]T),
	}
}

func NewConcurrentMapFrom[T any](data map[string]T) *ConcurrentMap[T] {
	return &ConcurrentMap[T]{
		data: data,
	}
}
