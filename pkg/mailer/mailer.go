package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"time"
)

// Message represents an outbound email.
type Message struct {
	From    mail.Address
	To      mail.Address
	Subject string
	Body    string
}

// SMTPConfig holds SMTP credentials and host/port.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Send dispatches an email via SMTP with STARTTLS when available.
func Send(ctx context.Context, cfg SMTPConfig, msg Message) error {
	if cfg.Host == "" {
		return fmt.Errorf("smtp host required")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: cfg.Host}
		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	if cfg.Username != "" || cfg.Password != "" {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(msg.From.Address); err != nil {
		return err
	}
	if err := client.Rcpt(msg.To.Address); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(serialize(msg))); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func serialize(msg Message) string {
	headers := map[string]string{
		"From":         msg.From.String(),
		"To":           msg.To.String(),
		"Subject":      msg.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"utf-8\"",
	}

	var raw string
	for k, v := range headers {
		raw += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	raw += "\r\n" + msg.Body
	return raw
}
