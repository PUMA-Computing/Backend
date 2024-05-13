package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
)

type MailgunService struct {
	domain      string
	apiKey      string
	senderEmail string
	mailgun     *mailgun.MailgunImpl
}

func NewMailgunService(domain, apiKey, senderEmail string) *MailgunService {
	return &MailgunService{
		domain:      domain,
		apiKey:      apiKey,
		senderEmail: senderEmail,
		mailgun:     mailgun.NewMailgun(domain, apiKey),
	}
}

func (ms *MailgunService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	body := "Click the link below to verify your email\n" +
		"http://localhost:8080/api/v1/auth/verify-email?token=" + token
	return ms.sendEmail(to, subject, body)
}

func (ms *MailgunService) sendEmail(toEmail, subject, body string) error {
	message := ms.mailgun.NewMessage(
		ms.senderEmail,
		subject,
		body,
		toEmail,
	)

	_, _, err := ms.mailgun.Send(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}
