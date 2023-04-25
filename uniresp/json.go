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

package uniresp

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// ActionError represents a basic user action error (e.g. a wrong parameter,
// non-existing record etc.)
type ActionError struct {
	error
}

// MarshalJSON serializes the error to JSON
func (me ActionError) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Error())
}

func NewActionErrorFrom(err error) ActionError {
	return ActionError{err}
}

// NewActionError creates an Action error from provided message using
// a newly defined general error as the original error
func NewActionError(msg string, args ...any) ActionError {
	return ActionError{fmt.Errorf(msg, args...)}
}

// ErrorResponse describes a wrapping object for all error HTTP responses
type ErrorResponse struct {
	Error   *ActionError `json:"error"`
	Details []string     `json:"details"`
	Code    int          `json:"code"`
}

// WriteJSONResponse writes 'value' to an HTTP response encoded as JSON
func WriteJSONResponse(w http.ResponseWriter, value any) {
	jsonAns, err := json.Marshal(value)
	if err != nil {
		log.Err(err).Msg("failed to encode a result to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonAns)
}

// WriteJSONResponseWithStatus writes 'value' to an HTTP response encoded as JSON
func WriteJSONResponseWithStatus(w http.ResponseWriter, status int, value any) {
	jsonAns, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(jsonAns)
}

func testEtagValues(headerValue, testValue string) bool {
	if headerValue == "" {
		return false
	}
	for _, item := range strings.Split(headerValue, ", ") {
		if strings.HasPrefix(item, "\"") && strings.HasSuffix(item, "\"") {
			val := item[1 : len(item)-1]
			if val == testValue {
				return true
			}

		} else {
			log.Warn().Msgf("Invalid ETag value: %s", item)
		}
	}
	return false
}

// WriteCacheableJSONResponse writes 'value' to an HTTP response encoded as JSON
// but before doing that it calculates a checksum of the JSON and in case it is
// equal to provided 'If-Match' header, 304 is returned. Otherwise a value with
// ETag header is returned.
func WriteCacheableJSONResponse(w http.ResponseWriter, req *http.Request, value any) {
	jsonAns, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 3600*24*30))
		crc := crc32.ChecksumIEEE(jsonAns)
		newEtag := fmt.Sprintf("chksm-%d", crc)
		reqEtagString := req.Header.Get("If-Match")
		if testEtagValues(reqEtagString, newEtag) {
			http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)

		} else {
			w.Header().Set("Etag", newEtag)
			w.Write(jsonAns)
		}
	}
}

// WriteJSONErrorResponse writes 'aerr' to an HTTP error response as JSON
func WriteJSONErrorResponse(w http.ResponseWriter, aerr ActionError, status int, details ...string) {
	ans := &ErrorResponse{
		Code:    status,
		Error:   &aerr,
		Details: details,
	}
	jsonAns, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	w.Write(jsonAns)
}

// WriteCustonmJSONErrorResponse writes any JSON serializable object as an HTTP error response.
// In case the value cannot be serialized into JSON, the function will write error
// 500 (Internal Server Error).
func WriteCustomJSONErrorResponse(w http.ResponseWriter, value any, status int, details ...string) {
	jsonAns, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	w.Write(jsonAns)
}
