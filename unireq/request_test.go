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

package unireq

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSuperfluousURLArgs(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("foo", "10")
	args.Add("bar", "hit")
	req.URL.RawQuery = args.Encode()
	ans := CheckSuperfluousURLArgs(req, []string{"foo", "bar"})
	assert.NoError(t, ans)
}

func TestCheckSuperfluousURLArgs1Excess(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("foo", "10")
	args.Add("bar", "hit")
	args.Add("alien", "xxx")
	req.URL.RawQuery = args.Encode()
	ans := CheckSuperfluousURLArgs(req, []string{"foo", "bar"})
	assert.Error(t, ans)
}

func TestClientIPForwFor(t *testing.T) {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Set("X-Forwarded-For", "192.168.1.17")
	req.Header.Set("X-Client-Ip", "192.168.1.20")
	ip := ClientIP(req)
	assert.Equal(t, net.ParseIP("192.168.1.17"), ip)
}

func TestClientIPClientIP(t *testing.T) {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Set("X-Client-Ip", "192.168.1.17")
	req.Header.Set("X-Real-Ip", "192.168.1.20")
	ip := ClientIP(req)
	assert.Equal(t, net.ParseIP("192.168.1.17"), ip)
}

func TestClientIPRealIp(t *testing.T) {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Set("X-Real-Ip", "192.168.1.17")
	req.RemoteAddr = "192.168.1.20"
	ip := ClientIP(req)
	assert.Equal(t, net.ParseIP("192.168.1.17"), ip)
}

func TestClientIPRemoteAddr(t *testing.T) {
	req := &http.Request{Header: make(http.Header)}
	req.RemoteAddr = "192.168.1.17"
	ip := ClientIP(req)
	assert.Equal(t, net.ParseIP("192.168.1.17"), ip)
}
