/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package email

import (
	"fmt"

	"github.com/gzuidhof/starcloud/starcloud/constants"
	"github.com/matcornic/hermes/v2"
	"github.com/spf13/viper"
)

const defaultFrom = constants.EmailDefaultFrom
const defaultButtonColor = "#3c64ff"
const defaultSignature = "Best"

type HermesEmailer struct {
	emailer Emailer
	hermes  hermes.Hermes

	systemName string
}

func NewHermesEmailer(client Emailer) HermesEmailer {

	return HermesEmailer{
		systemName: viper.GetString("system.name"),
		emailer:    client,
		hermes: hermes.Hermes{
			// Optional Theme
			// Theme: new(Default)
			Product: hermes.Product{
				// Appears in header & footer of e-mails
				Name: viper.GetString("system.name"),
				Link: viper.GetString("system.url"),
				// Optional product logo
				// Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
				Copyright: viper.GetString("system.name") + " - " + viper.GetString("system.url"),
			},
		},
	}
}

func (h *HermesEmailer) SendHermesEmail(from string, recipients []string, subject string, email hermes.Email) error {

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.hermes.GenerateHTML(email)
	if err != nil {
		return fmt.Errorf("failed to generate html for email with subject %s: %v", subject, err)
	}

	// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
	emailText, err := h.hermes.GeneratePlainText(email)
	if err != nil {
		return fmt.Errorf("failed to generate plaintext for email with subject %s: %v", subject, err)
	}

	return h.emailer.SendEmail(from, recipients, subject, []byte(emailBody), []byte(emailText))
}

func (h *HermesEmailer) SendWelcomeConfirmEmailEmail(userEmail string, username string, displayName string, confirmEmailUrl string) error {
	subject := fmt.Sprintf("Welcome to %s, please verify your e-mail address", h.systemName)

	e := hermes.Email{
		Body: hermes.Body{
			Name: displayName,
			Intros: []string{
				fmt.Sprintf("Welcome to %s, your details:", h.systemName),
			},
			Dictionary: []hermes.Entry{
				{Key: "Display Name", Value: displayName},
				{Key: "Username (use this to log in)", Value: username},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started please click here:",
					Button: hermes.Button{
						Text:  "Verify e-mail address",
						Link:  confirmEmailUrl,
						Color: defaultButtonColor,
					},
				},
			},
			Signature: defaultSignature,
		},
	}

	return h.SendHermesEmail(defaultFrom, []string{userEmail}, subject, e)
}

func (h *HermesEmailer) SendResetPasswordEmail(userEmail string, username string, displayName string, resetPasswordUrl string) error {
	subject := fmt.Sprintf("%s password reset", h.systemName)

	e := hermes.Email{
		Body: hermes.Body{
			Name: displayName,
			Intros: []string{
				"Someone - hopefully you - requested to reset your password. If it wasn't you, don't worry about it and just ignore this e-mail.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Username (use this to log in)", Value: username},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To reset your password click here",
					Button: hermes.Button{
						Text:  "Reset My Password",
						Link:  resetPasswordUrl,
						Color: defaultButtonColor,
					},
				},
			},
			Signature: defaultSignature,
		},
	}

	return h.SendHermesEmail(defaultFrom, []string{userEmail}, subject, e)
}

// func (h *HermesEmailer) SendConfirmEmailEmail(userEmail string, username string, displayName string, confirmEmailUrl string) error {
// 	subject := fmt.Sprintf("Welcome to %s, please verify your e-mail address", h.systemName)

// 	e := hermes.Email{
// 		Body: hermes.Body{
// 			Name: displayName,
// 			Actions: []hermes.Action{
// 				{
// 					Instructions: "To confirm your e-mail address please click here:",
// 					Button: hermes.Button{
// 						Text:  "Verify e-mail address",
// 						Link:  confirmEmailUrl,
// 						Color: defaultButtonColor,
// 					},
// 				},
// 			},
// 			Signature: defaultSignature,
// 		},
// 	}

// 	return h.SendHermesEmail(defaultFrom, []string{userEmail}, subject, e)
// }
