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
	"strings"
	"unicode/utf8"
)

// SmartTruncate truncates an input string so it does
// not exceed maxSize while respecting the "space"
// character.
// In case the resulting string is less than 20% of
// the original, an alternative "non-smart" cut
// is performed.
func SmartTruncate(inStr string, maxSize int) string {
	if maxSize < 0 {
		panic("negative maxSize")
	}
	if maxSize == 0 {
		return ""
	}
	sLength := utf8.RuneCountInString(inStr)
	if sLength < maxSize {
		return inStr
	}
	var prevSpace int
	for i, s := range inStr {
		if i >= maxSize {
			break
		}
		if s == ' ' {
			prevSpace = i
		}
	}
	ans := []rune(inStr)[:prevSpace]
	if prevSpace > 0 && float64(len(ans)) > 0.2*float64(len([]rune(inStr))) {
		return string([]rune(inStr)[:prevSpace]) + "\u2026"
	}
	return string([]rune(inStr)[:maxSize]) + "\u2026"
}

// JoinAny joins any items to a string based on mapFn which converts individual
// items to strings.
func JoinAny[T any](elems []T, mapFn func(v T) string, sep string) string {
	var ans strings.Builder
	for i := 0; i < len(elems); i++ {
		if i > 0 {
			ans.WriteString(sep)
		}
		ans.WriteString(mapFn(elems[i]))
	}
	return ans.String()
}
