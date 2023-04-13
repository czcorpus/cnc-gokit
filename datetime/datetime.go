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

import "time"

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
