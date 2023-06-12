package postoffice

import (
	"YenExpress/config"
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type PostMan[T any] struct {
	To               string
	Subject          string
	MailTemplatePath string
	MailBodyVal      T
	HTMLBody         bytes.Buffer
}

func (m *PostMan[T]) LoadTemplate() error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.ParseFiles(wd + m.MailTemplatePath)
	if err != nil {
		return err
	}
	return t.Execute(&m.HTMLBody, m.MailBodyVal)
}

func (m *PostMan[T]) SendMail() (*MailResponse, error) {
	err := m.LoadTemplate()
	if err != nil {
		log.Println("error from template loader ==> ", err)
		return &MailResponse{}, err
	}
	client := sendgrid.NewSendClient(config.SendGridAPIKey)
	message := mail.NewSingleEmail(mail.NewEmail("YenExpress", config.EmailAccountAddress),
		m.Subject, mail.NewEmail(strings.Split(m.To, "@")[0], m.To),
		"", m.HTMLBody.String())
	response, err := client.Send(message)
	if err != nil {
		log.Println("error within delivery method ==> ", err)
		return &MailResponse{}, err
	} else {
		return &MailResponse{
			StatusCode: response.StatusCode,
			Body:       response.Body,
		}, nil
	}
}

type OneTimePassword struct {
	OTP       string
	Validity  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (otp *OneTimePassword) Marshal() ([]byte, error) {
	return json.Marshal(otp)
}

func (otp *OneTimePassword) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &otp); err != nil {
		return err
	}
	return nil
}

type MailResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}
