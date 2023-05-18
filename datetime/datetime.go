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
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	durationPattern = regexp.MustCompile(`(\d+)([smhdwy])`)
)

// GetCurrentDatetime returns current UTC date and time in ISO 8601
// date format without timezone info (e.g. 2023-04-13T08:19:03)
func GetCurrentDatetime() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

// GetCurrentDatetimeIn returns current local date and time in ISO 8601
// date format without timezone info (e.g. 2023-04-13T08:19:03)
func GetCurrentDatetimeIn(loc *time.Location) string {
	return time.Now().In(loc).Format("2006-01-02T15:04:05")
}

// FormatDatetime creates an ISO 8601 formatted string
// based on provided time.
func FormatDatetime(dt time.Time) string {
	return dt.Format("2006-01-02T15:04:05")
}

func DurationToHMS(dur time.Duration) string {
	numSec := int(dur.Seconds())
	var sPrefix string
	if numSec < 0 {
		sPrefix = "-"
		numSec = -numSec
	}
	hours := numSec / 3600
	rest := numSec % 3600
	mins := rest / 60
	rest = rest % 60
	return fmt.Sprintf("%s%02d:%02d:%02d", sPrefix, hours, mins, rest)
}

// ParseDuration decodes string-encoded durations
// like 10s, 17h, 3y, or even '5d 7m
func ParseDuration(v string) (time.Duration, error) {
	srch := durationPattern.FindAllStringSubmatch(v, -1)
	var total time.Duration
	var check int
	for _, v := range srch {
		check += len(v[0])
	}
	if check != len(v) {
		return total, fmt.Errorf("failed to parse duration expression %s", v)
	}
	for _, v := range srch {
		tv, err := strconv.Atoi(v[1])
		if err != nil {
			return total, fmt.Errorf("failed to parse duration (value: %s)", v)
		}
		switch v[2] {
		case "s":
			total += time.Duration(tv) * time.Second
		case "m":
			total += time.Duration(tv) * time.Minute
		case "h":
			total += time.Duration(tv) * time.Hour
		case "d":
			total += time.Duration(tv) * 24 * time.Hour
		case "w":
			total += time.Duration(tv) * 7 * 24 * time.Hour
		case "y":
			total += time.Duration(tv) * 365 * 24 * time.Hour
		default:
			return total, fmt.Errorf("unknown code %s", v[2])
		}
	}
	return total, nil
}
