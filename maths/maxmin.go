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

type EssentialNumTypes interface {
	int | int64 | float32 | float64
}

func asInt(v any) int {
	vx, ok := v.(int)
	if !ok {
		panic("invalid type for asInt")
	}
	return vx
}

func asInt64(v any) int64 {
	vx, ok := v.(int64)
	if !ok {
		panic("invalid type for asInt64")
	}
	return vx
}

func asFloat32(v any) float32 {
	vx, ok := v.(float32)
	if !ok {
		panic("invalid type for asFloat32")
	}
	return vx
}

func asFloat64(v any) float64 {
	vx, ok := v.(float64)
	if !ok {
		panic("invalid type for float64")
	}
	return vx
}

// Max finds maximum number for int, int64, float32, float64
func Max[T EssentialNumTypes](v1 ...T) T {
	v1x := any(v1[0])
	switch v1x.(type) {
	case int:
		var maxValIdx int
		for i, x := range v1 {
			curr := asInt(x)
			if curr > asInt(v1[maxValIdx]) {
				maxValIdx = i
			}
		}
		return v1[maxValIdx]
	case int64:
		var maxValIdx int
		for i, x := range v1 {
			curr := asInt64(x)
			if curr > asInt64(v1[maxValIdx]) {
				maxValIdx = i
			}
		}
		return v1[maxValIdx]
	case float32:
		var maxValIdx int
		for i, x := range v1 {
			curr := asFloat32(x)
			if curr > asFloat32(v1[maxValIdx]) {
				maxValIdx = i
			}
		}
		return v1[maxValIdx]
	case float64:
		var maxValIdx int
		for i, x := range v1 {
			curr := asFloat64(x)
			if curr > asFloat64(v1[maxValIdx]) {
				maxValIdx = i
			}
		}
		return v1[maxValIdx]
	}
	panic("invalid type for Max")
}

// Min finds manimum number for int, int64, float32, float64
func Min[T EssentialNumTypes](v1 ...T) T {
	v1x := any(v1[0])
	switch v1x.(type) {
	case int:
		var minValIdx int
		for i, x := range v1 {
			curr := asInt(x)
			if curr < asInt(v1[minValIdx]) {
				minValIdx = i
			}
		}
		return v1[minValIdx]
	case int64:
		var minValIdx int
		for i, x := range v1 {
			curr := asInt64(x)
			if curr < asInt64(v1[minValIdx]) {
				minValIdx = i
			}
		}
		return v1[minValIdx]
	case float32:
		var minValIdx int
		for i, x := range v1 {
			curr := asFloat32(x)
			if curr < asFloat32(v1[minValIdx]) {
				minValIdx = i
			}
		}
		return v1[minValIdx]
	case float64:
		var minValIdx int
		for i, x := range v1 {
			curr := asFloat64(x)
			if curr < asFloat64(v1[minValIdx]) {
				minValIdx = i
			}
		}
		return v1[minValIdx]
	}
	panic("invalid type for Min")

}
