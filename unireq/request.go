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
	"strings"

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

// IsAIBot provides basic detection of mainstream AI bots.
//
// The function is intended for informational purposes only
// (e.g. for API usage stats) and should not be used to decide
// who to ban, etc. as the detection is not very sophisticated.
func IsAIBot(req *http.Request) bool {
	return strings.Contains(req.UserAgent(), "GPT") ||
		strings.Contains(req.UserAgent(), "Claude") ||
		strings.Contains(req.UserAgent(), "Perplexity") ||
		strings.Contains(req.UserAgent(), "Gemini-Deep-Research") // most stuff is done by GoogleBot though
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

// RequireURLIntArgOrFail is a variant of GetURLIntArgOrFail without a default value
// which triggers an error in case the value is not found.
func RequireURLIntArgOrFail(ctx *gin.Context, name string) (int, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		uniresp.RespondWithErrorJSON(
			ctx,
			fmt.Errorf("missing argument %s", name),
			http.StatusBadRequest,
		)
		return 0, false
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

// GetURLFloatArgOrFail reads a string-encoded float argument from URL query.
// If not set, then `dflt` is returned (i.e. value not present is considered a non-error).
// The second returned value is an "OK" flag.
// In case of an error, the function writes a HTTP response and returns
// false as a second argument.
func GetURLFloatArgOrFail(ctx *gin.Context, name string, dflt float64) (float64, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		return dflt, true
	}
	tmp := ctx.Request.URL.Query().Get(name)
	value, err := strconv.ParseFloat(tmp, 32)
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

// RequireURLFloatArgOrFail is a variant of GetURLFloatArgOrFail without a default value
// which triggers an error in case the value is not found.
func RequireURLFloatArgOrFail(ctx *gin.Context, name string) (float64, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		uniresp.RespondWithErrorJSON(
			ctx,
			fmt.Errorf("missing argument %s", name),
			http.StatusBadRequest,
		)
		return 0, false
	}
	tmp := ctx.Request.URL.Query().Get(name)
	value, err := strconv.ParseFloat(tmp, 32)
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

// RequireURLBoolArgOrFail is a variant of GetURLBoolArgOrFail without a default value
// which triggers an error in case the value is not found.
func RequireURLBoolArgOrFail(ctx *gin.Context, name string, dflt bool) (bool, bool) {
	if !ctx.Request.URL.Query().Has(name) {
		uniresp.RespondWithErrorJSON(
			ctx,
			fmt.Errorf("missing argument %s", name),
			http.StatusBadRequest,
		)
		return false, false
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
