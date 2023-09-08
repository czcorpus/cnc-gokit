// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
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

func TestOnlineMean(t *testing.T) {
	var om OnlineMean
	om2 := om.Add(1.0).Add(2.0).Add(3.0).Add(4.0).Add(5.0)
	assert.Equal(t, 3.0, om2.Mean())
	assert.InDelta(t, 1.58114, om2.Stdev(), 0.00001)
}

func TestOnlineMeanNoVariance(t *testing.T) {
	var om OnlineMean
	om2 := om.Add(3.0).Add(3.0).Add(3.0).Add(3.0).Add(3.0)
	assert.Equal(t, 3.0, om2.Mean())
	assert.InDelta(t, 0.0, om2.Stdev(), 0.00001)
}

func TestOnlineMeanOneValue(t *testing.T) {
	var om OnlineMean
	om2 := om.Add(7.3)
	assert.Equal(t, 7.3, om2.Mean())
	assert.InDelta(t, 0.0, om2.Stdev(), 0.00001)
}

func TestOnlineMeanNo(t *testing.T) {
	var om OnlineMean
	assert.Equal(t, 0.0, om.Mean())
	assert.InDelta(t, 0.0, om.Stdev(), 0.00001)
}
