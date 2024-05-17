package services

import (
	"Backend/configs"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"log"
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

	// Load BaseURL from configs
	baseURL := configs.LoadConfig().BaseURL

	// Change the hardcoded URL to use the baseURL from configs
	body := "Click the link below to verify your email\n" +
		baseURL + "/auth/verify-email?token=" + token

	log.Println("Sending email to: ", to)
	log.Println("Subject: ", subject)
	log.Println("Body: ", body)

	if err := ms.sendEmail(to, subject, body); err != nil {
		return err
	}

	log.Println()

	return nil
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
