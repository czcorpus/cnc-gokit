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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGet(t *testing.T) {
	clist := NewCircularList[string](5)
	clist.Append("A")
	assert.Equal(t, "A", clist.Get(0))
	assert.Equal(t, 1, clist.Len())
}

func TestGetPanicsOnIndexOutOfRange(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	assert.Panics(t, func() {
		clist.Get(1)
	})
}

func TestAddMoreThanCapacity(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	clist.Append("B")
	clist.Append("C")
	clist.Append("D") // now se start to rewrite A
	clist.Append("E") // here we rewrite B
	assert.Equal(t, "C", clist.Get(0))
	assert.Equal(t, "D", clist.Get(1))
	assert.Equal(t, "E", clist.Get(2))
	assert.Equal(t, 3, clist.Len())
}

func TestPrependOnNonFull(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("B")
	clist.Append("C")
	clist.Prepend("A")
	assert.Equal(t, "B", clist.items[1])
	assert.Equal(t, "C", clist.items[2])
	assert.Equal(t, "A", clist.items[0])
}

func TestPrependOnFull(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	clist.Append("B")
	clist.Append("C")
	clist.Prepend("D")
	assert.Equal(t, "A", clist.items[0]) // TODO the following 3 items are too low-level
	assert.Equal(t, "B", clist.items[1])
	assert.Equal(t, "D", clist.items[2])
	assert.Equal(t, "D", clist.Last())
	assert.Equal(t, "A", clist.Head())
}

func TestPrependOnEmpty(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Prepend("X")
	assert.Equal(t, "X", clist.Head())
	assert.Equal(t, "X", clist.Last())
}

func TestAppendWithPrepend(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	clist.Append("B")
	clist.Append("C")
	fmt.Println("next ", clist.nextIdx)
	clist.Prepend("D")
	fmt.Println("next ", clist.nextIdx)
	clist.Append("E")
	fmt.Println("next ", clist.nextIdx)
	assert.Equal(t, "B", clist.Head())
	assert.Equal(t, "E", clist.Last())
}

func TestHead(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	assert.Equal(t, "A", clist.Head())
	clist.Append("B")
	assert.Equal(t, "A", clist.Head())
	clist.Append("C")
	assert.Equal(t, "A", clist.Head())
	clist.Append("D")
	assert.Equal(t, "B", clist.Head())
	clist.Append("E") // here we rewrite B
	assert.Equal(t, "C", clist.Head())
}

func TestHeadPanicsOnEmpty(t *testing.T) {
	clist := NewCircularList[int](3)
	assert.Panics(t, func() {
		clist.Head()
	})
}

func TestLast(t *testing.T) {
	clist := NewCircularList[string](3)
	clist.Append("A")
	assert.Equal(t, "A", clist.Last())
	clist.Append("B")
	assert.Equal(t, "B", clist.Last())
	clist.Append("C")
	assert.Equal(t, "C", clist.Last())
	clist.Append("D")
	assert.Equal(t, "D", clist.Last())
	clist.Append("E")
	assert.Equal(t, "E", clist.Last())
}

func TestLastPanicsOnEmpty(t *testing.T) {
	clist := NewCircularList[int](3)
	assert.Panics(t, func() {
		clist.Last()
	})
}

type record struct {
	name    string
	created int
}

func TestShiftUntil(t *testing.T) {
	clist := NewCircularList[record](5)
	clist.Append(record{"action1", 10})
	clist.Append(record{"action2", 14})
	clist.Append(record{"action3", 16})
	clist.Append(record{"action4", 19})
	clist.Append(record{"action5", 23})
	clist.Append(record{"action6", 40})
	clist.Append(record{"action7", 49})

	clist.ShiftUntil(func(item record) bool {
		return item.created < 23
	})
	assert.Equal(t, "action5", clist.Get(0).name)
	assert.Equal(t, "action6", clist.Get(1).name)
	assert.Equal(t, "action7", clist.Get(2).name)
	assert.Equal(t, 3, clist.Len())
}

func TestShiftUntilOnEmpty(t *testing.T) {
	clist := NewCircularList[record](5)
	assert.NotPanics(t, func() {
		clist.ShiftUntil(func(item record) bool {
			return true
		})
	})
}

func TestShiftUntilTooGreedy(t *testing.T) {
	clist := NewCircularList[string](4)
	clist.Append("action1")
	clist.Append("action2")
	clist.Append("action3")
	clist.Append("action4")
	assert.NotPanics(t, func() {
		clist.ShiftUntil(func(item string) bool {
			return true
		})
	})
}

func TestForEach(t *testing.T) {
	clist := NewCircularList[string](4)
	clist.Append("action1")
	clist.Append("action2")
	clist.Append("action3")
	clist.Append("action4")
	clist.Append("action5")
	tmp := make([]string, 0, 5)
	clist.ForEach(func(i int, v string) bool {
		tmp = append(tmp, v)
		return true
	})
	assert.Equal(t, []string{"action2", "action3", "action4", "action5"}, tmp)
}

func TestForEachOnEmpty(t *testing.T) {
	var cnt int
	clist := NewCircularList[string](4)
	clist.ForEach(func(i int, v string) bool {
		cnt++
		return true
	})
	assert.Equal(t, 0, cnt)
}

func TestGOBEncodeDecode(t *testing.T) {
	clist := NewCircularList[string](4)
	clist.Append("a")
	clist.Append("b")
	clist.Append("c")
	clist.Append("d")
	origNextIdx := clist.nextIdx
	origNumUnused := clist.numUnused
	origItems := clist.items
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(clist)
	assert.NoError(t, err)

	var clist2 CircularList[string]
	buf2 := bytes.NewBuffer(buf.Bytes())
	decoder := gob.NewDecoder(buf2)
	err = decoder.Decode(&clist2)
	assert.NoError(t, err)
	assert.Equal(t, origNumUnused, clist2.numUnused)
	assert.Equal(t, origNextIdx, clist2.nextIdx)
	assert.Equal(t, origItems, clist2.items)
	// following two lines are not that necessary
	assert.Equal(t, "a", clist2.Head())
	assert.Equal(t, "d", clist2.Last())
}
