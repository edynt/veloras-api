package utils

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/edynnt/veloras-api/pkg/config"
)

type EmailAddress struct {
	Address string
	Name    string
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

var mailConfig *config.Config

func InitMail(cfg *config.Config) {
	mailConfig = cfg
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}

func send(to []string, from string, htmlTemplate string) error {
	contentEmail := Mail{
		From: EmailAddress{
			Address: from,
			Name:    "Go E-Commerce",
		},
		To:      to,
		Subject: "OTP Code",
		Body:    htmlTemplate,
	}

	messageMail := BuildMessage(contentEmail)
	// send smtp
	auth := smtp.PlainAuth("", mailConfig.SMTP.User, mailConfig.SMTP.Password, mailConfig.SMTP.Host)
	err := smtp.SendMail(mailConfig.SMTP.Host+":"+mailConfig.SMTP.Port, auth, from, to, []byte(messageMail))
	if err != nil {
		return err
	}
	return nil
}
