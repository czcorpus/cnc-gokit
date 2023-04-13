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

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationToHMS(t *testing.T) {
	d1 := time.Duration(12353 * time.Second) // = 3 * 3600 + 25 * 60 + 53
	ans := DurationToHMS(d1)
	assert.Equal(t, "03:25:53", ans)
}

func TestDurationToHMSZero(t *testing.T) {
	d1 := time.Duration(0)
	ans := DurationToHMS(d1)
	assert.Equal(t, "00:00:00", ans)
}

func TestDurationToHMSNegative(t *testing.T) {
	d1 := time.Duration(-8259 * time.Second)
	ans := DurationToHMS(d1)
	assert.Equal(t, "-02:17:39", ans)
}
