package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(
	templateFile, username, email string, data any, isSandbox bool,
) error {

	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	temp, err := template.ParseFS(FS, fmt.Sprintf("templates/%s", UserWelcomeTemplate))
	if err != nil {
		return nil
	}

	subject := new(bytes.Buffer)
	err = temp.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = temp.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := range maxRetries {
		_, err := m.client.Send(message)
		if err != nil {

			slog.Error("Failed to send email", "email", email, "attempt", i+1, "maxRetries", maxRetries)
			slog.Error("Error", "error", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)

}
