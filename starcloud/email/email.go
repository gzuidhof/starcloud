/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package email

import (
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type Emailer interface {
	SendEmail(from string, recipients []string, subject string, html []byte, text []byte) error
}

type SMTPEmailClient struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
}

type MockEmailClient struct {
	SentEmailTexts []string
}

func (c *MockEmailClient) SendEmail(from string, recipients []string, subject string, html []byte, text []byte) error {
	c.SentEmailTexts = append(c.SentEmailTexts, string(text))
	return nil
}

func (c *MockEmailClient) GetLastSentEmail() string {
	return c.SentEmailTexts[len(c.SentEmailTexts)-1]
}

func NewSMTPClient() SMTPEmailClient {
	return SMTPEmailClient{
		smtpHost:     viper.GetString("email.smtp_host"),
		smtpPort:     viper.GetString("email.smtp_port"),
		smtpUsername: viper.GetString("email.smtp_username"),
		smtpPassword: viper.GetString("email.smtp_password"),
	}
}

func (c *SMTPEmailClient) SendEmail(from string, recipients []string, subject string, html []byte, text []byte) error {
	e := &email.Email{
		To:      recipients,
		From:    from,
		Subject: subject,
		Text:    text,
		HTML:    html,
		Headers: textproto.MIMEHeader{},
	}
	err := e.Send(c.smtpHost+":"+c.smtpPort, smtp.PlainAuth("", c.smtpUsername, c.smtpPassword, c.smtpHost))

	if err != nil {
		return fmt.Errorf("SendEmail failed: %v", err)
	}
	return nil
}
