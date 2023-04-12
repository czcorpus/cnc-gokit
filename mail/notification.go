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
	"bytes"
	"cnc-gokit/datetime"
	"fmt"
	"net/smtp"
	"strings"
)

type NotificationConf struct {
	Sender       string   `json:"sender"`
	Receivers    []string `json:"receivers"`
	SMTPServer   string   `json:"smtpServer"`
	SMTPUsername string   `json:"smtpUsername"`
	SMTPPassword string   `json:"smtpPassword"`
	// Signature defines multi-language signature for notification e-mails
	Signature map[string]string `json:"signature"`
}

// SendNotification sends a general e-mail notification based on
// a respective monitoring configuration. The 'alarmToken' argument
// can be nil - in such case the 'turn of the alarm' text won't be
// part of the message.
func SendNotification(conf *NotificationConf, subject string, msgParagraphs ...string) error {
	client, err := smtp.Dial(conf.SMTPServer)
	if err != nil {
		return err
	}
	defer client.Close()

	client.Mail(conf.Sender)
	for _, rcpt := range conf.Receivers {
		client.Rcpt(rcpt)
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	headers := make(map[string]string)
	headers["From"] = conf.Sender
	headers["To"] = strings.Join(conf.Receivers, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	body := ""
	for k, v := range headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	for _, par := range msgParagraphs {
		body += "<p>" + par + "</p>\r\n\r\n"
	}
	body += fmt.Sprintf("<p>Generated at %s</p>\r\n\r\n", datetime.GetCurrentDatetime())

	buf := bytes.NewBufferString(body)
	_, err = buf.WriteTo(wc)
	return err
}
