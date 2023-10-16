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
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/czcorpus/cnc-gokit/datetime"
)

// NotificationConf configures e-mail-based notification
type NotificationConf struct {
	Sender       string   `json:"sender"`
	Recipients   []string `json:"recipients"`
	SMTPServer   string   `json:"smtpServer"`
	SMTPUsername string   `json:"smtpUsername"`
	SMTPPassword string   `json:"smtpPassword"`
	// Signature defines multi-language signature for notification e-mails
	Signature map[string]string `json:"signature"`
}

// WithRecipients creates a new NotificationConf instance
// with recipients overwritten by the provided ones
func (nc NotificationConf) WithRecipients(r ...string) NotificationConf {
	return NotificationConf{
		Sender:       nc.Sender,
		Recipients:   r,
		SMTPServer:   nc.SMTPServer,
		SMTPUsername: nc.SMTPUsername,
		SMTPPassword: nc.SMTPPassword,
		Signature:    nc.Signature,
	}
}

// Notification represents a general notification e-mail
// subject and body. The contents of `Paragraphs`
// is always escaped by the notification function.
// (for raw HTML formatted notifications, use `FormattedNotification`)
type Notification struct {
	Subject    string
	Paragraphs []string
}

// FormattedNotification allows sending raw html message
// (i.e. nothing is HTMP escaped)
type FormattedNotification struct {
	Subject string
	Divs    []string
}

type EmailNotificationMessage interface {
	Notification | FormattedNotification
}

func mkParBody(not Notification) string {
	b := strings.Builder{}
	for _, p := range not.Paragraphs {
		b.WriteString("<p>" + html.EscapeString(p) + "</p>\r\n\r\n")
	}
	return b.String()
}

func mkDivBody(msg FormattedNotification) string {
	b := strings.Builder{}
	for _, div := range msg.Divs {
		b.WriteString("<div>" + div + "</div>\r\n\r\n")
	}
	return b.String()
}

// AsParagraph escapes `text` and wraps it into <p> and </p>
func AsParagraph(text string, style string) string {
	if style != "" {
		return fmt.Sprintf("<p style=\"%s\">%s</p>", style, html.EscapeString(text))
	}
	return "<p>" + html.EscapeString(text) + "</p>"
}

// SendNotification sends a general e-mail notification.
// Based on configuration, it is able to use SMTP servers
// requiring TLS and authentication (see Dial()).
func SendNotification[T EmailNotificationMessage](conf *NotificationConf, location *time.Location, msg T) error {
	client, err := DialServer(conf.SMTPServer, conf.SMTPUsername, conf.SMTPPassword)
	if err != nil {
		return err
	}
	defer client.Close()

	client.Mail(conf.Sender)
	for _, rcpt := range conf.Recipients {
		client.Rcpt(rcpt)
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	var aMsg any = msg
	var subj, msgBody string
	switch tMsg := aMsg.(type) {
	case FormattedNotification:
		subj = tMsg.Subject
		msgBody = mkDivBody(tMsg)
	case Notification:
		subj = tMsg.Subject
		msgBody = mkParBody(tMsg)
	}

	headers := make(map[string]string)
	headers["From"] = conf.Sender
	headers["To"] = strings.Join(conf.Recipients, ",")
	headers["Subject"] = subj
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	body := strings.Builder{}
	for k, v := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	body.WriteString("\r\n")
	body.WriteString(msgBody)
	body.WriteString(fmt.Sprintf("<p>Generated at %s</p>\r\n\r\n", datetime.GetCurrentDatetimeIn(location)))

	buf := bytes.NewBufferString(body.String())
	_, err = buf.WriteTo(wc)
	return err
}
