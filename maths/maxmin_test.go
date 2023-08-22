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

func TestMinInt(t *testing.T) {
	assert.Equal(t, -7, Min(10, -3, 1, 2, 11, -6, -7, 1, 0))
	assert.Equal(t, 7, Min(7))
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 11, Max(10, -3, 1, 2, 11, -6, -7, 1, 0))
	assert.Equal(t, 7, Max(7))
}

func TestMinInt64(t *testing.T) {
	assert.Equal(t, int64(-7), Min(int64(10), -3, 1, 2, 11, -6, -7, 1, 0))
	assert.Equal(t, int64(7), Min(int64(7)))
}

func TestMaxInt64(t *testing.T) {
	assert.Equal(t, int64(11), Max(int64(10), -3, 1, 2, 11, -6, -7, 1, 0))
	assert.Equal(t, int64(7), Max(int64(7)))
}

func TestMinFloat32(t *testing.T) {
	assert.Equal(t, float32(-7.37), Min(float32(10.7), -3.3, 1.11324, 2.554, 11.73, -6.74, -7.37, 1.0, 0.7))
	assert.Equal(t, float32(7.9), Min(float32(7.9)))
}

func TestMaxFloat32(t *testing.T) {
	assert.Equal(t, float32(11.73), Max(float32(10.7), -3.3, 1.11324, 2.554, 11.73, -6.74, -7.37, 1.0, 0.7))
	assert.Equal(t, float32(7), Max(float32(7)))
}

func TestMaxFloat64(t *testing.T) {
	assert.Equal(t, 11.73, Max(10.7, -3.3, 1.11324, 2.554, 11.73, -6.74, -7.37, 1.0, 0.7))
	assert.Equal(t, 7.9, Max(7.9))
}

func TestMinFloat64(t *testing.T) {
	assert.Equal(t, -7.37, Min(10.7, -3.3, 1.11324, 2.554, 11.73, -6.74, -7.37, 1.0, 0.7))
	assert.Equal(t, 7, Min(7))
}
