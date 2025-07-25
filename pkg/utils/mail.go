package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/edynnt/veloras-api/pkg/global"
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
	auth := smtp.PlainAuth("", global.Config.SMTP.User, global.Config.SMTP.Password, global.Config.SMTP.Host)
	err := smtp.SendMail(global.Config.SMTP.Host+":"+global.Config.SMTP.Port, auth, from, to, []byte(messageMail))
	if err != nil {
		return err
	}
	return nil
}

func SendTemplateEmailOtp(to []string, from string, nameTemplate string, dataTemplate map[string]interface{}) error {
	fmt.Println("to", to)
	fmt.Println("from", from)

	htmlBody, err := getMailTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}

	return send(to, from, htmlBody)
}

func getMailTemplate(nameTemplate string, dataTemplate map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates-email/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}

	return htmlTemplate.String(), nil
}
