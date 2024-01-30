// Copyright 2024 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2024 Institute of the Czech National Corpus,
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

package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	TransportMaxIdleConns        = 100
	TransportMaxConnsPerHost     = 100
	TransportMaxIdleConnsPerHost = 80
)

type httpClientConf struct {
	followRedirects bool
	idleConnTimeout time.Duration
	timeout         time.Duration
	tslSkipVerify   bool
}

func WithFollowRedirects() func(args *httpClientConf) {
	return func(args *httpClientConf) {
		args.followRedirects = true
	}
}

// WithTimeout specifies request timeout (incl. connecting, redirects,
// body reading).
func WithTimeout(value time.Duration) func(args *httpClientConf) {
	return func(args *httpClientConf) {
		args.timeout = value
	}
}

// WithIdleConnTimeout specifies the maximum amount of time an idle
// (keep-alive) connection will remain idle before closing
// itself. To set no limit, use zero.
func WithIdleConnTimeout(value time.Duration) func(args *httpClientConf) {
	return func(args *httpClientConf) {
		args.idleConnTimeout = value
	}
}

// WithInsecureSkipVerify disables certificate verification.
// Please note that this is intended for internal network
// usage due to possible security implications.
func WithInsecureSkipVerify() func(args *httpClientConf) {
	return func(args *httpClientConf) {
		args.tslSkipVerify = true
	}
}

func New(options ...func(args *httpClientConf)) *http.Client {
	var conf httpClientConf
	for _, opt := range options {
		opt(&conf)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = TransportMaxIdleConns
	transport.MaxConnsPerHost = TransportMaxConnsPerHost
	transport.MaxIdleConnsPerHost = TransportMaxIdleConnsPerHost
	transport.IdleConnTimeout = conf.idleConnTimeout
	if conf.tslSkipVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := http.Client{
		Timeout:   conf.timeout,
		Transport: transport,
	}
	if conf.followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &client

}
