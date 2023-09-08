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

type OnlineMean struct {
	count int
	mean  float64
	stdev float64
	m2    float64
}

func (m OnlineMean) Add(incoming float64) OnlineMean {
	m.count++
	delta := incoming - m.mean
	m.mean = m.mean + delta/float64(m.count)
	delta2 := incoming - m.mean
	m.m2 += delta * delta2
	if m.count < 2 {
		m.stdev = 0

	} else {
		m.stdev = math.Sqrt(m.m2 / float64(m.count-1))
	}
	return m
}

func (m OnlineMean) Mean() float64 {
	return m.mean
}

func (m OnlineMean) Stdev() float64 {
	return m.stdev
}
