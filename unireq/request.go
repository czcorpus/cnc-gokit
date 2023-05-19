// Copyright 2022 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2022 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2022 Institute of the Czech National Corpus,
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

package unireq

import (
	"fmt"
	"net/http"

	"github.com/czcorpus/cnc-gokit/collections"
)

// CheckSuperfluousURLArgs allows checking for presence of supported only
// arguments in URL. It returns first error it encounters.
func CheckSuperfluousURLArgs(req http.Request, allowedArgs []string) error {
	for name := range req.URL.Query() {
		if !collections.SliceContains(allowedArgs, name) {
			return fmt.Errorf("unsupported URL argument %s", name)
		}
	}
	return nil
}
