package email

import (
	"fmt"
	"net/smtp"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	smtpUser = "ggtaiga06@gmail.com"
	smtpPass = "xgarnsesmpvgktzd"
)

func SendVerificationEmail(to string, code string) error {
	from := smtpUser
	pass := smtpPass

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Email Verification - GoBooking\r\n\r\n"+
		"Your verification code is: %s\r\n", from, to, code)

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
