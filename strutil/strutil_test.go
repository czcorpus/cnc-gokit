// Copyright 2024 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2024 Institute of the Czech National Corpus,
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

package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmartTruncateNoSpaceInput(t *testing.T) {
	s2 := SmartTruncate("0123456789", 5)
	assert.Equal(t, "01234\u2026", s2)
}

func TestSmartTruncateShorterResult(t *testing.T) {
	s2 := SmartTruncate("012 34 567 8 9", 9)
	assert.Equal(t, "012 34\u2026", s2)
}

func TestSmartTruncateExactResult(t *testing.T) {
	s2 := SmartTruncate("012 34 567 8 9", 11)
	assert.Equal(t, "012 34 567\u2026", s2)
}

func TestSmartTruncateEmpty(t *testing.T) {
	s2 := SmartTruncate("", 11)
	assert.Equal(t, "", s2)
}

func TestSmartTruncateTooLongLimit(t *testing.T) {
	s2 := SmartTruncate("012 34 567 8 9", 200)
	assert.Equal(t, "012 34 567 8 9", s2)
}

func TestSmartTruncateZeroLimit(t *testing.T) {
	s2 := SmartTruncate("012 34 567 8 9", 0)
	assert.Equal(t, "", s2)
}

func TestSmartTruncateNegativeLimit(t *testing.T) {
	assert.Panics(t, func() {
		SmartTruncate("012 34 567 8 9", -5)
	})
}
