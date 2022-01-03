package service

import (
	"HundredToFive/internal/config"
	"HundredToFive/pkg/email"
	"fmt"
)

const (
	verificationLinkTmpl = "http://%s/Reset?token=%s"
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig
}
type verificationEmailInput struct {
	VerificationLink string
}

func NewEmailService(sender email.Sender, config config.EmailConfig) *EmailService {
	return &EmailService{sender: sender, config: config}
}

func (s *EmailService) SendUserVerificationEmail(input VerificationEmailInput) error {

	subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)

	templateInput := verificationEmailInput{s.createVerificationSecretCode(input.Domain, input.Token)}

	sendInput := email.SendEmailInput{
		To:      input.Email,
		Subject: subject,
	}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return fmt.Errorf("service.SendUserVerificationEmail: %w", err)
	}

	return s.sender.Send(sendInput)

}

func (s *EmailService) createVerificationSecretCode(domain, token string) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, token)
}
