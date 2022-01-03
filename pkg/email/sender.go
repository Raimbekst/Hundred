package email

import (
	"HundredToFive/pkg/logger"
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

type Sender interface {
	Send(input SendEmailInput) error
}

func (e *SendEmailInput) GenerateBodyFromHTML(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		logger.Errorf("failed to parse file %s:%s", templateFileName, err.Error())
		return fmt.Errorf("email.GenerateBodyFromHTML : %w", err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return fmt.Errorf("email.GenerateBodyFromHTML : %w", err)
	}
	e.Body = buf.String()

	return nil
}

func (e *SendEmailInput) Validate() error {
	if e.To == "" {
		return fmt.Errorf("email.GenerateBodyFromHTML: %w", errors.New("empty to"))
	}
	if e.Subject == "" || e.Body == "" {
		return fmt.Errorf("email.GenerateBodyFromHTML: %w", errors.New("empty subject/body"))
	}
	return nil
}
