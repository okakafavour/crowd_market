package services

import (
	"fmt"
	"net/smtp"
	"CROWD_MARKET/config"
)

func SendVerificationEmail(toEmail, verificationCode string) error {
	from := config.EmailFrom
	password := config.EmailPassword

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	link := fmt.Sprintf("http://localhost:8080/verify?code=%s", verificationCode)

	message := []byte(fmt.Sprintf("Subject: Verify your email\n\nClick here to verify your account: %s", link))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	return err
}
