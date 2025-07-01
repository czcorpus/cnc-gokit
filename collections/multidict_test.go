// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2025 Institute of the Czech National Corpus,
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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiDictAutoInitializes(t *testing.T) {
	var md Multidict[string]
	assert.NotPanics(t, func() {
		md.Add("foo", "foo-value")
	})
}

func TestMultidictSetGet(t *testing.T) {
	md := NewMultidict[string]()
	md.Add("foo", "foo-value")
	v := md.Get("foo")
	assert.Equal(t, []string{"foo-value"}, v)
}

type mdItem struct {
	K string
	V []string
}

func TestMultidictIterate(t *testing.T) {
	md := NewMultidict[string]()
	md.Add("foo", "foo-v1")
	md.Add("foo", "foo-v2")
	md.Add("bar", "bar-v1")

	itemTest := make([]mdItem, 0, 2)
	for k, v := range md.Iterate {
		itemTest = append(itemTest, mdItem{k, v})
	}
	sort.Slice(itemTest, func(i, j int) bool {
		return strings.Compare(itemTest[i].K, itemTest[j].K) > 0
	})
	assert.Equal(t, []mdItem{{"foo", []string{"foo-v1", "foo-v2"}}, {"bar", []string{"bar-v1"}}}, itemTest)
}

type mdFlatItem struct {
	K string
	V string
}

func TestMultidictIterateFlat(t *testing.T) {
	md := NewMultidict[string]()
	md.Add("foo", "foo-v1")
	md.Add("foo", "foo-v2")
	md.Add("bar", "bar-v1")

	itemTest := make([]mdFlatItem, 0, 2)
	for k, v := range md.IterateFlat {
		itemTest = append(itemTest, mdFlatItem{k, v})
	}
	sort.Slice(itemTest, func(i, j int) bool {
		return strings.Compare(itemTest[i].K, itemTest[j].K) > 0
	})
	assert.Equal(t, []mdFlatItem{{"foo", "foo-v1"}, {"foo", "foo-v2"}, {"bar", "bar-v1"}}, itemTest)
}
