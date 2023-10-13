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

package maths

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type item int

func (item item) Freq() int {
	return int(item)
}

type items []item

func (items items) Len() int {
	return len(items)
}

func (items items) Get(idx int) FreqInfo {
	return items[idx]
}

func TestStatsGetQuartiles(t *testing.T) {
	//                     *       *           *
	//             0   1   2   3   4   5   6   7   8   9
	data := items{10, 15, 20, 25, 30, 35, 40, 45, 50, 55}
	q, err := GetQuartiles[FreqInfo](data)
	assert.NoError(t, err)
	assert.Equal(t, 2, q.Q1Idx)
	assert.Equal(t, 20, q.Q1)
	assert.Equal(t, 4, q.Q2Idx)
	assert.Equal(t, 30, q.Q2)
	assert.Equal(t, 7, q.Q3Idx)
	assert.Equal(t, 45, q.Q3)
	assert.Equal(t, 25, q.IQR())
}

func TestStatsGetQuartilesOnEmpty(t *testing.T) {
	data := items{}
	_, err := GetQuartiles[FreqInfo](data)
	assert.ErrorIs(t, ErrTooSmallDataset, err)
}

func TestStatsGetQuartilesOnSix(t *testing.T) {
	data := items{10, 15, 20, 25, 30, 35}
	_, err := GetQuartiles[FreqInfo](data)
	assert.ErrorIs(t, ErrTooSmallDataset, err)
}
