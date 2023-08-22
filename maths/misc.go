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

import "math"

type floats interface {
	float32 | float64
}

func RoundToN[T floats](value T, places int) T {
	multiplier := math.Pow(10, float64(places))
	var v any = value
	switch tv := v.(type) {
	case float64:
		return T(math.Round(tv*multiplier) / multiplier)
	case float32:
		return T(float32(math.Round(float64(tv)*multiplier) / multiplier))
	default:
		panic("usupported value for RoundToN")
	}
}
