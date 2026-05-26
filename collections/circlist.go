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
	"bytes"
	"encoding/gob"
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

// calcIdx converts a logical index (0 = oldest item) to the physical slice index.
func (clist *CircularList[T]) calcIdx(idx int) int {
	return (clist.nextIdx + clist.numUnused + idx) % len(clist.items)
}

// Append adds a new item to the end of the list. In case the
// free capacity is depleted, then the oldest item is replaced by
// this new one.
func (clist *CircularList[T]) Append(v T) {
	clist.AppendAndGetInternalIdx(v)
}

// AppendAndGetInternalIdx adds a new item to the end of the list. In case the
// free capacity is depleted, then the oldest item is replaced by
// this new one.
// The method returns item's internal index for possible additional manipulation.
func (clist *CircularList[T]) AppendAndGetInternalIdx(v T) int {
	usedIdx := clist.nextIdx
	clist.items[usedIdx] = v
	clist.nextIdx = (usedIdx + 1) % len(clist.items)
	if clist.numUnused > 0 {
		clist.numUnused--
	}
	return usedIdx
}

// Prepend inserts v at the logical beginning (oldest position) of the list.
// If the list has unused capacity, existing items are shifted right and v occupies
// index 0. If the list is full, the newest item is silently overwritten, since there
// is no free slot available at the beginning.
func (clist *CircularList[T]) Prepend(v T) {
	if clist.numUnused == 0 {
		idx := (len(clist.items) + clist.nextIdx - 1) % len(clist.items)
		clist.items[idx] = v

	} else {
		for i := clist.nextIdx; i > 0; i-- {
			clist.items[i] = clist.items[i-1]
		}
		clist.items[0] = v
		clist.numUnused--
		clist.nextIdx++
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
	if idx < 0 {
		idx = len(clist.items) + idx
	}
	return clist.items[idx]
}

// ShiftUntil
//
// Deprecated: use DeleteWhile
func (clist *CircularList[T]) ShiftUntil(fn func(item T) bool) {
	clist.DeleteWhile(fn)
}

// DeleteWhile removes items starting from the oldest one, continuing while
// fn returns true, and stopping at the first item for which fn returns false.
// If the list becomes empty before fn returns false, the function returns
// without error. This is useful for discarding expired records, e.g. log
// entries older than a given timestamp.
func (clist *CircularList[T]) DeleteWhile(fn func(item T) bool) {
	if clist.Len() == 0 {
		return
	}
	pred := fn(clist.Head())
	for pred {
		clist.numUnused++
		if clist.Len() == 0 {
			return
		}
		pred = fn(clist.Head())
	}
}

// Get returns an item based on its order from the oldest (0),
// to newest (Len() - 1).
func (clist *CircularList[T]) Get(idx int) T {
	if idx >= clist.Len() {
		panic(fmt.Sprintf("index out of range [%d] with length %d", idx, clist.Len()))
	}
	return clist.items[clist.calcIdx(idx)]
}

// GetByInternalIdx returns an item by its internal index, as obtained via
// AppendAndGetInternalIdx.
func (clist *CircularList[T]) GetByInternalIdx(idx int) T {
	return clist.items[idx]
}

// Len returns size of the list
func (clist *CircularList[T]) Len() int {
	return len(clist.items) - clist.numUnused
}

// ForEach is just an alias for Iterate
//
// Deprecated: use Iterate instead
func (clist *CircularList[T]) ForEach(fn func(i int, item T) bool) {
	clist.Iterate(fn)
}

// Iterate goes through all items from oldest to newest, calling yield for each.
// The first argument passed to yield is the item's internal (physical) array index,
// not a logical sequence number. Use that index with GetByInternalIdx if needed.
// Returning false from yield stops the iteration early.
//
// Deprecated: use either IterateInternal or IterateLogical based on your
// needs.
func (clist *CircularList[T]) Iterate(yield func(i int, item T) bool) {
	clist.IterateInternal(yield)
}

// IterateInternal goes through all items from oldest to newest, calling yield for each.
// The first argument passed to yield is the item's internal (physical) array index,
// not a logical sequence number. Use that index with GetByInternalIdx if needed.
// Returning false from yield stops the iteration early.
func (clist *CircularList[T]) IterateInternal(yield func(i int, item T) bool) {
	for i := 0; i < clist.Len(); i++ {
		ii := (clist.nextIdx + clist.numUnused + i) % len(clist.items)
		cnt := yield(ii, clist.items[ii])
		if !cnt {
			break
		}
	}
}

// IterateLogical goes through all items from oldest to newest, calling yield for each.
// The first argument passed to yield is the logical index (0 = oldest, Len()-1 = newest).
// Returning false from yield stops the iteration early.
func (clist *CircularList[T]) IterateLogical(yield func(i int, item T) bool) {
	for i := 0; i < clist.Len(); i++ {
		ii := (clist.nextIdx + clist.numUnused + i) % len(clist.items)
		cnt := yield(i, clist.items[ii])
		if !cnt {
			break
		}
	}
}

// IterateOverInternalRange returns an iterator over a range of internal indices
// [i1, i2] inclusive, where indices are obtained via AppendAndGetInternalIdx.
// Wrap-around ranges are supported: e.g. i1=13, i2=7 on a capacity-16 list
// iterates through 13, 14, 15, 0, 1, ..., 7.
// Note: if the list is not full, unused slots within the range will yield
// zero values of type T.
func (clist *CircularList[T]) IterateOverInternalRange(i1, i2 int) func(yield func(i int, item T) bool) {
	return func(fn func(i int, item T) bool) {
		var numIter int
		if i2 >= i1 {
			numIter = i2 - i1 + 1

		} else {
			numIter = len(clist.items) - i1 + i2 + 1
		}
		for i := i1; i < i1+numIter; i++ {
			ii := i % len(clist.items)
			cnt := fn(ii, clist.items[ii])
			if !cnt {
				break
			}
		}
	}
}

func (clist *CircularList[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(clist.items); err != nil {
		return []byte{}, fmt.Errorf("failed to GOB encode CircularList: %w", err)
	}
	if err := encoder.Encode(clist.nextIdx); err != nil {
		return []byte{}, fmt.Errorf("failed to GOB encode CircularList: %w", err)
	}
	if err := encoder.Encode(clist.numUnused); err != nil {
		return []byte{}, fmt.Errorf("failed to GOB encode CircularList: %w", err)
	}
	return buf.Bytes(), nil
}

func (clist *CircularList[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&clist.items); err != nil {
		return fmt.Errorf("failed to GOB decode CircularList: %w", err)
	}
	if err := decoder.Decode(&clist.nextIdx); err != nil {
		return fmt.Errorf("failed to GOB decode CircularList: %w", err)
	}
	if err := decoder.Decode(&clist.numUnused); err != nil {
		return fmt.Errorf("failed to GOB decode CircularList: %w", err)
	}
	return nil
}

// NewCircularList is the recommended factory function for CircularList.
// The capacity parameter sets the maximum number of items the list can hold.
// All capacity slots are allocated upfront as a fixed-size array, so the full
// memory is reserved immediately regardless of how many items have been appended.
func NewCircularList[T any](capacity int) *CircularList[T] {
	return &CircularList[T]{
		items:     make([]T, capacity),
		numUnused: capacity,
	}
}
