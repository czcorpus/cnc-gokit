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
	"fmt"
	"strconv"
	"strings"
)

type anyNumber interface {
	int | int64 | float32 | float64
}

func numToString[T anyNumber](n T) string {
	switch nx := any(n).(type) {
	case int:
		return strconv.Itoa(nx)
	case int64:
		return strconv.FormatInt(nx, 10)
	case float32, float64:
		return fmt.Sprintf("%01.2f", nx)
	}
	return ""
}

// JoinNumbersAsString joins numbers to a string where items
// are separated with comma and space (e.g. '1, 173, 407').
// Float numbers are formatted with 2 decimal places.
func JoinNumbersAsString[T anyNumber](nums []T) string {
	var b strings.Builder
	for i, n := range nums {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(numToString(n))
	}
	return b.String()
}
