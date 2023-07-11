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
	assert.Panics(t, func() {
		clist.ShiftUntil(func(item record) bool {
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
