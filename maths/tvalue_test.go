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

func TestTValueTwoTail(t *testing.T) {
	tv, err := TValueTwoTail(10, Significance_0_05)
	assert.NoError(t, err)
	assert.InDelta(t, 2.2281, tv, 0.001)

	tv, err = TValueTwoTail(10, Significance_0_01)
	assert.NoError(t, err)
	assert.InDelta(t, 3.1693, tv, 0.001)

	tv, err = TValueTwoTail(40, Significance_0_02)
	assert.NoError(t, err)
	assert.InDelta(t, 2.4233, tv, 0.001)
}

func TestTValueTwoTailErr(t *testing.T) {
	_, err := TValueTwoTail(0, Significance_0_01)
	assert.ErrorIs(t, ErrValueNotAvailable, err)
}

func TestTValueTwoTailLarge(t *testing.T) {
	tv, err := TValueTwoTail(400, Significance_0_05)
	assert.NoError(t, err)
	assert.InDelta(t, 1.984, tv, 0.001)
}

func TestTValueTwoTailLarge2(t *testing.T) {
	tv, err := TValueTwoTail(600, Significance_0_05)
	assert.NoError(t, err)
	assert.InDelta(t, 1.962, tv, 0.001)
}
