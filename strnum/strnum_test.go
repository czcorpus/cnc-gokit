// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
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

package strnum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinNumbersAsStringEmpty(t *testing.T) {
	ans := JoinNumbersAsString([]float32{})
	assert.Equal(t, "", ans)
}

func TestJoinNumbersAsStringSingle(t *testing.T) {
	ans := JoinNumbersAsString([]float32{0.747})
	assert.Equal(t, "0.75", ans)
}

func TestJoinNumbersAsStringFloats(t *testing.T) {
	ans := JoinNumbersAsString([]float32{3.1416, -2.9917, 189.1, 297.875, 0.738})
	assert.Equal(t, "3.14, -2.99, 189.10, 297.88, 0.74", ans)
}

func TestJoinNumbersAsStringInts(t *testing.T) {
	ans := JoinNumbersAsString([]int{3, -2, 189, 297, 0})
	assert.Equal(t, "3, -2, 189, 297, 0", ans)
}

func TestJoinNumbersAsStringInt64s(t *testing.T) {
	ans := JoinNumbersAsString([]int64{3, 2, 189, 297, 0})
	assert.Equal(t, "3, 2, 189, 297, 0", ans)
}
