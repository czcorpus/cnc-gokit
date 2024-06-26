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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myInt int

func (v myInt) Compare(other Comparable) int {
	tOther := other.(myInt)
	if v > tOther {
		return 1

	} else if v == other {
		return 0
	}
	return -1
}

func TestBinTreeAdd(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(10, 20, 8, 15, 4, 21, 20)

	assert.Equal(t, 7, bt.Len())

	assert.Equal(t, myInt(4), bt.Get(0))
	assert.Equal(t, myInt(8), bt.Get(1))
	assert.Equal(t, myInt(10), bt.Get(2))
	assert.Equal(t, myInt(15), bt.Get(3))
	assert.Equal(t, myInt(20), bt.Get(4))
	assert.Equal(t, myInt(20), bt.Get(5))
	assert.Equal(t, myInt(21), bt.Get(6))
}

func TestBinTreeRemove(t *testing.T) {
	var bt BinTree[myInt]
	// 4, 8, 10, 15, 20, 20, 21
	bt.Add(10, 20, 8, 15, 4, 21, 20)
	x := bt.Remove(2)
	assert.Equal(t, myInt(10), x)

	assert.Equal(t, myInt(4), bt.Get(0))
	assert.Equal(t, myInt(8), bt.Get(1))
	assert.Equal(t, myInt(15), bt.Get(2))
	assert.Equal(t, myInt(20), bt.Get(3))
	assert.Equal(t, myInt(20), bt.Get(4))
	assert.Equal(t, myInt(21), bt.Get(5))
}

func TestBinTreeRemoveFromRGTLinkedListLike(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(4, 8, 10, 15, 20, 20, 21)
	x := bt.Remove(2)
	assert.Equal(t, myInt(10), x)

	assert.Equal(t, myInt(4), bt.Get(0))
	assert.Equal(t, myInt(8), bt.Get(1))
	assert.Equal(t, myInt(15), bt.Get(2))
	assert.Equal(t, myInt(20), bt.Get(3))
	assert.Equal(t, myInt(20), bt.Get(4))
	assert.Equal(t, myInt(21), bt.Get(5))
}

func TestBinTreeRemoveFromLFTLinkedListLike(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(21, 20, 15, 10, 8, 4)
	x := bt.Remove(3)
	assert.Equal(t, myInt(15), x)

	assert.Equal(t, myInt(4), bt.Get(0))
	assert.Equal(t, myInt(8), bt.Get(1))
	assert.Equal(t, myInt(10), bt.Get(2))
	assert.Equal(t, myInt(20), bt.Get(3))
	assert.Equal(t, myInt(21), bt.Get(4))
}

func TestBinTreeRemoveFromEmpty(t *testing.T) {
	var bt BinTree[myInt]
	assert.Panics(t, func() {
		bt.Remove(0)
	})
}

func TestBinTreeRemoveLeaf(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(4, 8, 10, 15, 20, 20, 21)
	x := bt.Remove(6)
	assert.Equal(t, myInt(21), x)

	assert.Equal(t, myInt(4), bt.Get(0))
	assert.Equal(t, myInt(8), bt.Get(1))
	assert.Equal(t, myInt(10), bt.Get(2))
	assert.Equal(t, myInt(15), bt.Get(3))
	assert.Equal(t, myInt(20), bt.Get(4))
	assert.Equal(t, myInt(20), bt.Get(5))
}

func TestBinTreeGetEmpty(t *testing.T) {
	var bt BinTree[myInt]
	assert.Panics(t, func() {
		bt.Get(0)
	})
}

func TestBinTreeToSlice(t *testing.T) {
	var bt BinTree[myInt]
	// 4, 8, 10, 15, 20, 20, 21
	bt.Add(10, 20, 8, 15, 4, 21, 20)
	slc := bt.ToSlice()
	assert.Equal(t, []myInt{4, 8, 10, 15, 20, 20, 21}, slc)
}

func TestBinTreeToSliceLinkedListLFT(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(21, 20, 15, 10, 8, 4)
	slc := bt.ToSlice()
	assert.Equal(t, []myInt{4, 8, 10, 15, 20, 21}, slc)
}

func TestBinTreeToSliceLinkedListRGT(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(4, 8, 10, 15, 20, 20, 21)
	slc := bt.ToSlice()
	assert.Equal(t, []myInt{4, 8, 10, 15, 20, 20, 21}, slc)
}

func TestBinTreeToSliceEmpty(t *testing.T) {
	var bt BinTree[myInt]
	slc := bt.ToSlice()
	assert.Equal(t, []myInt{}, slc)
}

func TestBinTreeGetOverflow(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(10, 20, 8, 15)
	assert.Panics(t, func() {
		bt.Get(20)
	})
}

func TestBinTreeGetNegativeIndex(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(10, 20, 8, 15)
	// 8, 10, 15, 20
	assert.Equal(t, myInt(20), bt.Get(-1))
	assert.Equal(t, myInt(15), bt.Get(-2))
	assert.Equal(t, myInt(10), bt.Get(-3))
	assert.Equal(t, myInt(8), bt.Get(-4))
}

func TestBinTreeGetNegativeIndexOverflow(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(10, 20, 8, 15)
	// 8, 10, 15, 20
	assert.Panics(t, func() {
		bt.Get(-5)
	})
}

func TestBinTreeForEach(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(237, 30, 10, 430, 1, 147, 128, 12)
	// 1, 10, 12, 30, 128, 147, 237, 430
	ans := strings.Builder{}
	bt.ForEach(func(i int, v myInt) bool {
		if i > 6 {
			return false
		}
		ans.WriteString(fmt.Sprintf("%d:%d|", i, v))
		return true
	})
	assert.Equal(t, "0:1|1:10|2:12|3:30|4:128|5:147|6:237|", ans.String())
}

func TestBinTreeForEachOnEmpty(t *testing.T) {
	var bt BinTree[myInt]
	var tst bool
	assert.NotPanics(t, func() {
		bt.ForEach(func(i int, v myInt) bool {
			tst = true
			return true
		})
	})
	assert.False(t, tst)
}

func TestUniqueSetting(t *testing.T) {
	var bt BinTree[myInt]
	bt.Add(2, 1, 3)
	assert.Equal(t, []myInt{1, 2, 3}, bt.ToSlice())
	bt.Add(1)
	bt.UniqValues = true
	bt.Add(2)
	bt.Add(3)
	assert.Equal(t, []myInt{1, 1, 2, 3}, bt.ToSlice())
}
