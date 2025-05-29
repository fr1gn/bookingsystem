package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const mailerSendURL = "https://api.mailersend.com/v1/email"

var MailerSendAPIKey = "mlsn.52e41b77fc003e12d149917b08be5cc04d14bcebfdd78e3f1bf9a4d17cbf8704"

type EmailPayload struct {
	From    EmailAddress   `json:"from"`
	To      []EmailAddress `json:"to"`
	Subject string         `json:"subject"`
	Text    string         `json:"text"`
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func SendVerificationEmail(to string, code string) error {
	payload := EmailPayload{
		From: EmailAddress{
			Email: "noreply@gobooking.kz",
			Name:  "GoBooking",
		},
		To: []EmailAddress{
			{Email: to},
		},
		Subject: "Email Verification - GoBooking",
		Text:    fmt.Sprintf("Your verification code is: %s", code),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", mailerSendURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+MailerSendAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("MailerSend returned status %d", resp.StatusCode)
	}

	return nil
}
