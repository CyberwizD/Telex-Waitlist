package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/CyberwizD/Telex-Waitlist/internal/models"
	"github.com/CyberwizD/Telex-Waitlist/internal/repository"
)

// ErrValidation is returned when user input fails validation.
var ErrValidation = errors.New("validation error")

// WaitlistService defines application logic for waitlist operations.
type WaitlistService interface {
	Submit(ctx context.Context, name, email string) (*models.WaitlistEntry, error)
	List(ctx context.Context, limit, offset int) ([]models.WaitlistEntry, int64, error)
}

type waitlistService struct {
	repo         repository.WaitlistRepository
	emailService EmailService
}

// NewWaitlistService constructs a WaitlistService.
func NewWaitlistService(repo repository.WaitlistRepository, emailSvc EmailService) WaitlistService {
	return &waitlistService{
		repo:         repo,
		emailService: emailSvc,
	}
}

func (s *waitlistService) Submit(ctx context.Context, name, email string) (*models.WaitlistEntry, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", ErrValidation)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("%w: invalid email", ErrValidation)
	}

	entry, err := s.repo.Create(ctx, name, email)
	if err != nil {
		return nil, err
	}

	// Send thank-you email asynchronously but track errors.
	if s.emailService != nil {
		go func() {
			_ = s.emailService.SendThankYou(context.Background(), email, name)
		}()
	}

	return entry, nil
}

func (s *waitlistService) List(ctx context.Context, limit, offset int) ([]models.WaitlistEntry, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	entries, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}
