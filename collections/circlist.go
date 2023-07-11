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
)

// CircularList is a structure allowing infinite appending of new items
// while rewriting the oldest (in terms of order of respective Append()
// operations) records if needed. It also allows removing oldest records
// based on a condition - here we expect that the values contain some value
// representing their order in which they have been added - typically
// it is a time information.
type CircularList[T any] struct {
	items     []T
	nextIdx   int
	numUnused int
}

func (clist *CircularList[T]) calcIdx(idx int) int {
	return (clist.nextIdx + clist.numUnused + idx) % len(clist.items)
}

// Append adds a new item to the end of the list. In case the
// free capacity is depleted, then the oldest item is replaced by
// this new one.
func (clist *CircularList[T]) Append(v T) {
	clist.items[clist.nextIdx] = v
	clist.nextIdx = (clist.nextIdx + 1) % len(clist.items)
	if clist.numUnused > 0 {
		clist.numUnused--
	}
}

// Head returns the oldest item of the list. In case the list
// is empty, panic() is caused.
func (clist *CircularList[T]) Head() T {
	if clist.Len() == 0 {
		panic("calling Head() on an empty CircularList")
	}
	idx := (clist.nextIdx + clist.numUnused) % len(clist.items)
	return clist.items[idx]
}

// Last returns the most recent item of the list. In case the list
// is empty, panic() is caused.
func (clist *CircularList[T]) Last() T {
	if clist.Len() == 0 {
		panic("calling Last() on an empty CircularList")
	}
	idx := (clist.nextIdx - 1) % len(clist.items)
	fmt.Println("idx = ", idx)
	if idx < 0 {
		idx = len(clist.items) + idx
	}
	return clist.items[idx]
}

// ShiftUntil removes old items starting from the oldest one
// and moving towards newer ones until 'fn' returns true.
// This can be used to e.g. clean old log records.
func (clist *CircularList[T]) ShiftUntil(fn func(item T) bool) {
	pred := fn(clist.Head())
	for pred {
		clist.numUnused++
		pred = fn(clist.Head())
	}
}

// Get returns an item with a specific index.
// Please note that for CircularList, this
// method is not very usable as the index
// does not represent anything specific for
// the outside world.
func (clist *CircularList[T]) Get(idx int) T {
	if idx >= clist.Len() {
		panic(fmt.Sprintf("index out of range [%d] with length %d", idx, clist.Len()))
	}
	return clist.items[clist.calcIdx(idx)]
}

// Len returns size of the list
func (clist *CircularList[T]) Len() int {
	return len(clist.items) - clist.numUnused
}

// ForEach runs a function fn for all the items
// starting from the oldest one. The iteration
// continues until fn returns true.
func (clist *CircularList[T]) ForEach(fn func(i int, item T) bool) {
	for i := 0; i < clist.Len(); i++ {
		ii := (clist.nextIdx + clist.numUnused + i) % len(clist.items)
		cnt := fn(ii, clist.items[ii])
		if !cnt {
			break
		}
	}
}

// NewCircularList is a recommended factory function
// for CircularList.
// The parameter `capacity` defines max. number
// of items the instance will be able to store.
func NewCircularList[T any](capacity int) *CircularList[T] {
	return &CircularList[T]{
		items:     make([]T, capacity),
		numUnused: capacity,
	}
}
