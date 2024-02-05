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

import "unicode/utf8"

// SmartTruncate truncates an input string so it does
// not exceed maxSize while respecting the "space"
// character
func SmartTruncate(inStr string, maxSize int) string {
	if maxSize < 0 {
		panic("negative maxSize")
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
	if prevSpace > 0 {
		return string([]rune(inStr)[:prevSpace])
	}
	return string([]rune(inStr)[:maxSize])
}
