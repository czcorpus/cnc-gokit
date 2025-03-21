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

package util

import "github.com/czcorpus/cnc-gokit/maths"

// Max provides a maximum value out of the ones provided
//
// Deprecated: use `maths.Max` instead
func Max[T maths.EssentialNumTypes](v1 ...T) T {
	return maths.Max[T](v1...)
}

// Min provides a minimum value out of the ones provided
//
// Deprecated: use `maths.Min` instead
func Min[T maths.EssentialNumTypes](v1 ...T) T {
	return maths.Min[T](v1...)
}

// Or returns first non-zero item out of provides ones.
// If no argument is non-zero, it returns the zero value.
// Note: this is backported from go 1.22
func Or[T comparable](vals ...T) T {
	var zero T
	for _, val := range vals {
		if val != zero {
			return val
		}
	}
	return zero
}

// Ternary returns either ifTrue or ifFalse based on condition.
// (similar to the ternary operator: condition ? ifTrue : ifFalse)
func Ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
