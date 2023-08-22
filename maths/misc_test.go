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

func TestRoundToN(t *testing.T) {

	assert.Equal(t, 3.279, RoundToN(3.2789, 3))
	assert.Equal(t, 0.1079, RoundToN(0.1079, 4))
	assert.Equal(t, 4.0, RoundToN(3.6, 0))

	assert.Equal(t, float32(3.279), RoundToN(float32(3.2789), 3))
	assert.Equal(t, float32(0.1079), RoundToN(float32(0.1079), 4))
	assert.Equal(t, float32(4.0), RoundToN(float32(3.6), 0))
}
