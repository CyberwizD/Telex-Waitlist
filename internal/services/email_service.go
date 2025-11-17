package services

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/CyberwizD/Telex-Waitlist/internal/config"
	"github.com/CyberwizD/Telex-Waitlist/pkg/mailer"
)

// EmailService is responsible for sending outbound emails.
type EmailService interface {
	SendThankYou(ctx context.Context, recipientEmail, recipientName string) error
}

type smtpEmailService struct {
	cfg *config.Config
}

// NewEmailService constructs an SMTP-backed email sender.
func NewEmailService(cfg *config.Config) EmailService {
	return &smtpEmailService{cfg: cfg}
}

// SendThankYou sends a basic thank-you email after a waitlist submission.
func (s *smtpEmailService) SendThankYou(ctx context.Context, recipientEmail, recipientName string) error {
	if !s.cfg.EmailEnabled {
		return nil
	}

	if s.cfg.SMTPHost == "" || s.cfg.SMTPUsername == "" || s.cfg.SMTPPassword == "" || s.cfg.SMTPFrom == "" {
		return fmt.Errorf("smtp configuration incomplete")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	from := mail.Address{Name: s.cfg.AppName, Address: s.cfg.SMTPFrom}
	to := mail.Address{Name: recipientName, Address: recipientEmail}

	subject := fmt.Sprintf("Thanks for joining %s waitlist!", s.cfg.AppName)
	body := fmt.Sprintf("Hi %s,\n\nThanks for signing up for the %s waitlist. We'll keep you posted!\n\nCheers,\n%s Team\n", recipientName, s.cfg.AppName, s.cfg.AppName)

	return mailer.Send(ctx, mailer.SMTPConfig{
		Host:     s.cfg.SMTPHost,
		Port:     s.cfg.SMTPPort,
		Username: s.cfg.SMTPUsername,
		Password: s.cfg.SMTPPassword,
	}, mailer.Message{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	})
}
