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

package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// DialServer dials a SMTP server using provided arguments. The 'server' argument
// should also contain a port number (e.g. 'localhost:25'). Anything other than
// port 25 is considered to be a connection with authentication over TLS.
func DialServer(server, username, password string) (*smtp.Client, error) {
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP server info: %w", err)
	}
	if port == "25" {
		ans, err := smtp.Dial(server)
		if err != nil {
			return nil, fmt.Errorf("failed to dial: %w", err)
		}
		return ans, err
	}
	auth := smtp.PlainAuth("", username, password, host)
	client, err := smtp.Dial(server)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to StartTLS: %w", err)
	}
	err = client.Auth(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	return client, nil
}
