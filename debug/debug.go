// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2025 Department of Linguistics,
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

package debug

import (
	"errors"
	"fmt"
	"reflect"
)

// GetAddress returns an address of a respective value in case
// it is a pointer. Otherwise, it returns 0.
func GetAddress(value any) string {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		return fmt.Sprintf("0x%x", v.Pointer())
	}
	return "N/A"
}

// GetTypeAndAddress prints the actual type of the value and its address (if applicable).
// This is useful for inspecting values hidden behind interfaces to see the concrete type
// and memory address.
func GetTypeAndAddress(value any) string {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	if value == nil {
		return "Type: <nil>, Address: N/A\n"
	}

	address := "N/A"
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			address = "<nil pointer>"
		} else {
			address = fmt.Sprintf("0x%x", v.Pointer())
		}
	} else if v.CanAddr() {
		address = fmt.Sprintf("0x%x", v.Addr().Pointer())
	}

	return fmt.Sprintf("Type: %v, Address: %s\n", t, address)
}

// PrintErrorChain prints a chain of errors starting from the provided error,
// unwrapping each error in the chain and printing them with indentation.
func PrintErrorChain(err error) {
	if err == nil {
		return
	}

	level := 0
	for current := err; current != nil; current = errors.Unwrap(current) {
		indent := ""
		for i := 0; i < level; i++ {
			indent += "  "
		}
		fmt.Printf("%s%v\n", indent, current)
		level++
	}
}
