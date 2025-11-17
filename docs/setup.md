# Backend Setup Guide

This guide walks you through running the API locally.

## Prerequisites
- Go 1.24+
- PostgreSQL reachable via a connection string

## Environment
At minimum, set `DATABASE_URL`. Example for local Postgres:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable
```
The app will auto-load a `.env` file in the repo root if it exists (local/dev). Otherwise, export variables in your shell.
Other useful variables (see `.env.example` for defaults):
- `PORT` (default 8080)
- `APP_NAME`
- `ADMIN_API_KEY` (optional; enables listing)
- `ALLOWED_ORIGINS`
- `EMAIL_ENABLED`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM`
- `LOG_LEVEL`

On Windows PowerShell:
```powershell
$Env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable"
go run main.go
```

On macOS/Linux:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable"
go run main.go
```

## Database migration
Before running the API, apply the migration:
```bash
psql "$DATABASE_URL" -f internal/database/migrations/001_create_waitlist_table.up.sql
```

## Docker-compose option
```bash
cp .env.example .env
docker-compose up --build
# In another shell:
psql "postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable" \
  -f internal/database/migrations/001_create_waitlist_table.up.sql
```

## Common error
- `config error err="DATABASE_URL is required"`  
  Set `DATABASE_URL` in your environment (see examples above) before running `go run main.go`.
