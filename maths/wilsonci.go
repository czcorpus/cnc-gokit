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

// WilsonCI calculates Wilson confidence interval for a random
// variable with binomial distribution. The input arguments
// are represented as: `succ` successful trials out of `sampleSize`
func WilsonCI(succ float64, sampleSize int, signif SignificanceLevel) (float64, float64, error) {
	z, ok := zTable[signif]
	if !ok {
		return 0, 0, ErrUnsupportedSignifLevel
	}
	p := succ / float64(sampleSize)
	sq := z * math.Sqrt(p*(1-p)/float64(sampleSize)+math.Pow(z, 2)/(4*math.Pow(float64(sampleSize), 2)))
	denom := 1 + math.Pow(z, 2)/float64(sampleSize)
	a := p + math.Pow(z, 2)/(2*float64(sampleSize))
	return (a - sq) / denom, (a + sq) / denom, nil
}
