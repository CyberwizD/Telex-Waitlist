package models

import (
	"time"
)

// WaitlistEntry represents a single waitlist record.
type WaitlistEntry struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
