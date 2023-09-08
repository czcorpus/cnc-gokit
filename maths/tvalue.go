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

import (
	"math"
)

type SignificanceLevel string

func findTValue(df int, ci SignificanceLevel) (float64, error) {
	if df == 0 {
		return 0, ErrValueNotAvailable
	}
	col := idxMap[ci]
	if df <= 100 {
		return tTable[df][col], nil
	}
	if df <= 500 {
		return tTable[100][col], nil
	}
	return tTable[1000][col], nil
}

// TDistribConfInterval calculates a confidence interval
// for a sample mean and standard deviation in case
// population std. deviation is unknown and the values
// are "roughly normal".
// Please note that the function calculates respective
// t-values using a simple lookup table and is reliable
// up to 100 degrees of freedom. Higher values will be
// likely approximated by 1000 df which may or may not
// serve well.
// The provided confidence level is always applied
// in "two tails" mode.
func TDistribConfInterval(mean, stdev float64, sampleSize int, conf SignificanceLevel) (float64, float64, error) {
	tVal, err := findTValue(sampleSize-1, conf)
	if err != nil {
		return 0, 0, err
	}
	lft := mean + tVal*stdev/math.Sqrt(float64(sampleSize))
	rgt := mean + tVal*stdev/math.Sqrt(float64(sampleSize))
	return lft, rgt, nil
}

// TValueTwoTail gets t-value with two-tailed confidence level
func TValueTwoTail(df int, conf SignificanceLevel) (float64, error) {
	return findTValue(df, conf)
}
