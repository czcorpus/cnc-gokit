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

type ConcurrentMap[K comparable, T any] struct {
	sync.RWMutex
	data map[K]T
}

func (cm *ConcurrentMap[K, T]) Get(k K) T {
	cm.RLock()
	defer cm.RUnlock()
	return cm.data[k]
}

func (cm *ConcurrentMap[K, T]) GetWithTest(k K) (T, bool) {
	cm.RLock()
	defer cm.RUnlock()
	v, ok := cm.data[k]
	return v, ok
}

func (cm *ConcurrentMap[K, T]) HasKey(k K) bool {
	cm.RLock()
	defer cm.RUnlock()
	_, ok := cm.data[k]
	return ok
}

func (cm *ConcurrentMap[K, T]) Set(k K, v T) {
	cm.Lock()
	defer cm.Unlock()
	cm.data[k] = v
}

func (cm *ConcurrentMap[K, T]) Delete(k K) {
	delete(cm.data, k)
}

func (cm *ConcurrentMap[K, T]) ForEach(fn func(k K, v T)) {
	cm.RLock()
	defer cm.RUnlock()
	for k, v := range cm.data {
		fn(k, v)
	}
}

func (cm *ConcurrentMap[K, T]) Update(fn func(k K, v T) T) {
	cm.Lock()
	defer cm.Unlock()
	for k, v := range cm.data {
		cm.data[k] = fn(k, v)
	}
}

func (cm *ConcurrentMap[K, T]) Keys() []K {
	ans := make([]K, len(cm.data))
	i := 0
	for k, _ := range cm.data {
		ans[i] = k
		i++
	}
	return ans
}

func (cm *ConcurrentMap[K, T]) Values() []T {
	ans := make([]T, len(cm.data))
	i := 0
	for _, v := range cm.data {
		ans[i] = v
		i++
	}
	return ans
}

// AsMap creates a shallow copy of a map wrapped
// by this ConcurrentMap
func (cm *ConcurrentMap[K, T]) AsMap() map[K]T {
	cm.RLock()
	defer cm.RUnlock()
	ans := make(map[K]T)
	for k, v := range cm.data {
		ans[k] = v
	}
	return ans
}

// Len returns number of key-value pairs stored in the map
func (cm *ConcurrentMap[K, T]) Len() int {
	return len(cm.data)
}

func NewConcurrentMap[K comparable, T any]() *ConcurrentMap[K, T] {
	return &ConcurrentMap[K, T]{
		data: make(map[K]T),
	}
}

func NewConcurrentMapFrom[K comparable, T any](data map[K]T) *ConcurrentMap[K, T] {
	return &ConcurrentMap[K, T]{
		data: data,
	}
}
