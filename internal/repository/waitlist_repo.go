package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CyberwizD/Telex-Waitlist/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WaitlistRepository defines persistence operations for waitlist entries.
type WaitlistRepository interface {
	Create(ctx context.Context, name, email string) (*models.WaitlistEntry, error)
	List(ctx context.Context, limit, offset int) ([]models.WaitlistEntry, error)
	Count(ctx context.Context) (int64, error)
}

type waitlistRepo struct {
	db *pgxpool.Pool
}

// NewWaitlistRepository returns a Postgres-backed implementation.
func NewWaitlistRepository(db *pgxpool.Pool) WaitlistRepository {
	return &waitlistRepo{db: db}
}

func (r *waitlistRepo) Create(ctx context.Context, name, email string) (*models.WaitlistEntry, error) {
	const query = `
		INSERT INTO waitlist (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at;
	`

	row := r.db.QueryRow(ctx, query, name, email)
	entry := models.WaitlistEntry{}
	if err := row.Scan(&entry.ID, &entry.Name, &entry.Email, &entry.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("waitlist insert returned no rows")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *waitlistRepo) List(ctx context.Context, limit, offset int) ([]models.WaitlistEntry, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM waitlist
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.WaitlistEntry
	for rows.Next() {
		var e models.WaitlistEntry
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *waitlistRepo) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(*) FROM waitlist;`
	var total int64
	if err := r.db.QueryRow(ctx, query).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
