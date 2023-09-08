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

func TestWilsonCI(t *testing.T) {

	lft, rgt, err := WilsonCI(3, 12, Significance_0_05)
	assert.InDelta(t, 0.08894, lft, 0.001)
	assert.InDelta(t, 0.5323, rgt, 0.001)
	assert.NoError(t, err)
}
