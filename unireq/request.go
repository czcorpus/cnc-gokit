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
	"net"
	"net/http"
	"strconv"

	"github.com/czcorpus/cnc-gokit/collections"
	"github.com/czcorpus/cnc-gokit/uniresp"
	"github.com/gin-gonic/gin"
)

// CheckSuperfluousURLArgs allows checking for presence of supported only
// arguments in URL. It returns first error it encounters.
func CheckSuperfluousURLArgs(req *http.Request, allowedArgs []string) error {
	for name := range req.URL.Query() {
		if !collections.SliceContains(allowedArgs, name) {
			return fmt.Errorf("unsupported URL argument %s", name)
		}
	}
	return nil
}

// ClientIP tries to capture actual remote client IP address
// even if the client communicates with proxy servers.
// Please note that the "forwarded" HTTP header is not supported.
//
// In case nothing is found, nil is returned.
func ClientIP(req *http.Request) net.IP {
	src := req.Header.Get("x-forwarded-for")
	if src != "" {
		return net.ParseIP(src)
	}
	src = req.Header.Get("x-client-ip")
	if src != "" {
		return net.ParseIP(src)
	}
	src = req.Header.Get("x-real-ip")
	if src != "" {
		return net.ParseIP(src)
	}
	return net.ParseIP(req.RemoteAddr)
}

// GetURLIntArgOrFail reads a string-encoded integer argument from URL query.
// If not set, then `dflt` is returned (i.e. value not present is considered a non-error).
// The second returned value is an "OK" flag.
// In case of an error, the function writes a HTTP response and returns
// false as a second argument.
func GetURLIntArgOrFail(ctx *gin.Context, name string, dflt int) (int, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		return dflt, true
	}
	tmp := ctx.Request.URL.Query().Get(name)
	value, err := strconv.Atoi(tmp)
	if err != nil {
		uniresp.WriteJSONErrorResponse(
			ctx.Writer,
			uniresp.NewActionErrorFrom(err),
			http.StatusUnprocessableEntity,
		)
		return 0, false
	}
	return value, true
}

// GetURLBoolArgOrFail reads a string-encoded bool argument (= '1', '0') from URL query.
// If not set, then `dflt` is returned (i.e. value not present is considered a non-error).
// The second returned value is an "OK" flag.
// In case of an error, the function writes a HTTP response and returns
// false as a second argument.
func GetURLBoolArgOrFail(ctx *gin.Context, name string, dflt bool) (bool, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		return dflt, true
	}
	tmp := ctx.Request.URL.Query().Get(name)
	if tmp == "0" {
		return false, true

	} else if tmp == "1" {
		return true, true

	} else {
		err := fmt.Errorf("invalid URL bool value: %s", tmp)
		uniresp.WriteJSONErrorResponse(
			ctx.Writer,
			uniresp.NewActionErrorFrom(err),
			http.StatusUnprocessableEntity,
		)
		return false, false
	}
}
