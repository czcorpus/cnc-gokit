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
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
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

func TestGetURLBoolArgOrFail(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("foo", "1")
	args.Add("bar", "hit")
	req.URL.RawQuery = args.Encode()
	ctx := new(gin.Context)
	ctx.Request = req
	v, ok := GetURLBoolArgOrFail(ctx, "foo", false)
	assert.True(t, v)
	assert.True(t, ok)
}

func TestGetURLBoolArgOrFailDefault(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("bar", "hit")
	req.URL.RawQuery = args.Encode()
	ctx := new(gin.Context)
	ctx.Request = req
	v, ok := GetURLBoolArgOrFail(ctx, "foo", true)
	assert.True(t, v)
	assert.True(t, ok)
}

func TestGetURLBoolArgOrFailInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{URL: &url.URL{}}
	args := ctx.Request.URL.Query()
	args.Add("foo", "30")
	args.Add("bar", "hit")
	ctx.Request.URL.RawQuery = args.Encode()
	v, ok := GetURLBoolArgOrFail(ctx, "foo", false)
	assert.False(t, v)
	assert.False(t, ok)
}

func TestGetURLIntArgOrFail(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("foo", "37")
	args.Add("bar", "hit")
	req.URL.RawQuery = args.Encode()
	ctx := new(gin.Context)
	ctx.Request = req
	v, ok := GetURLIntArgOrFail(ctx, "foo", 0)
	assert.Equal(t, 37, v)
	assert.True(t, ok)
}

func TestGetURIntArgOrFailDefault(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	args := req.URL.Query()
	args.Add("bar", "hit")
	req.URL.RawQuery = args.Encode()
	ctx := new(gin.Context)
	ctx.Request = req
	v, ok := GetURLIntArgOrFail(ctx, "foo", 137)
	assert.Equal(t, 137, v)
	assert.True(t, ok)
}

func TestGetURLIntArgOrFailInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{URL: &url.URL{}}
	args := ctx.Request.URL.Query()
	args.Add("foo", "a30")
	args.Add("bar", "hit")
	ctx.Request.URL.RawQuery = args.Encode()
	v, ok := GetURLIntArgOrFail(ctx, "foo", 0)
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}
