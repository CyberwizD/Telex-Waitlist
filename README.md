# Telex Waitlist API

Lightweight waitlist backend built with Gin and PostgreSQL. Public submissions are stored and an optional thank-you email is sent via SMTP. An admin token can be used to page through entries.

## Stack
- Go 1.24, Gin HTTP framework
- PostgreSQL + pgx pool
- SMTP for transactional email
- Dockerfile + docker-compose for local/dev

## Features
- Public `POST /api/v1/waitlist` for name/email (no auth)
- Optional admin listing `GET /api/v1/waitlist` protected by `X-Admin-Token`
- Thank-you email on submission (toggle with `EMAIL_ENABLED`)
- CORS support with configurable origins
- Healthcheck at `/health`

## Project structure
```
.
├── main.go                       # App bootstrap: config, DB pool, router
├── go.mod / go.sum               # Module and deps
├── Dockerfile                    # Multi-stage build for container image
├── docker-compose.yml            # App + Postgres for local/dev
├── .env.example                  # Sample environment configuration
├── internal
│   ├── config/                   # Environment loading
│   ├── database/                 # pgx pool helper + SQL migrations
│   │   └── migrations/           # SQL migration files
│   ├── handlers/                 # HTTP handlers (waitlist)
│   ├── middleware/               # CORS and future middlewares
│   ├── models/                   # Domain models (waitlist entry)
│   ├── repository/               # Postgres data access (waitlist)
│   ├── routes/                   # Routes registration
│   └── services/                 # Business logic + email sender
└── pkg/                          # Reserved for shared libs (currently empty)
```

Additional frontend-focused docs: `docs/frontend.md`.

## Quickstart (Docker)
1) Copy `.env.example` to `.env` and fill values (set `DATABASE_URL`, SMTP creds, `ADMIN_API_KEY` if you want listing).
2) Start services: `docker-compose up --build`
3) Run migration into the Postgres container:
   ```bash
   psql "postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable" \
     -f internal/database/migrations/001_create_waitlist_table.up.sql
   ```
4) API available on `http://localhost:8080`.

## Running locally without Docker
```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs)   # or set vars manually on Windows
psql "$DATABASE_URL" -f internal/database/migrations/001_create_waitlist_table.up.sql
go run ./...
```

## Environment variables
- `PORT` – HTTP port (default `8080`)
- `APP_NAME` – used in email From/subject (default `Telex Waitlist`)
- `DATABASE_URL` – Postgres connection string (required)
- `ADMIN_API_KEY` – token for listing endpoint (leave empty to disable listing entirely)
- `ALLOWED_ORIGINS` – comma-separated CORS origins; empty allows any
- `EMAIL_ENABLED` – `true`/`false` to toggle sending emails
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` – SMTP settings
- `LOG_LEVEL` – `debug`, `info`, `warn`, or `error` (default `info`)
  - Note: the app will auto-load a `.env` file in the project root if present (for local dev). In production, set env vars via your process manager.

## API reference

### Health
`GET /health` → `200 {"status":"ok"}`

### Submit to waitlist (public)
- `POST /api/v1/waitlist`
- Body:
  ```json
  { "name": "Jane Doe", "email": "jane@example.com" }
  ```
- Responses:
  - `201` with saved entry
  - `400` validation error
  - `409` if email already exists

Example:
```bash
curl -X POST http://localhost:8080/api/v1/waitlist \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

### List waitlist entries (admin)
- `GET /api/v1/waitlist?limit=50&offset=0`
- Header: `X-Admin-Token: <ADMIN_API_KEY>`
- Responses:
  - `200` with `{ data, total, limit, offset }`
  - `401` invalid token
  - `403` if `ADMIN_API_KEY` not set in server env

Example:
```bash
curl "http://localhost:8080/api/v1/waitlist?limit=20&offset=0" \
  -H "X-Admin-Token: $ADMIN_API_KEY"
```

## Email behavior
- On successful submission, a thank-you email is sent via SMTP.
- Disable for local dev with `EMAIL_ENABLED=false`.
- Ensure `SMTP_FROM` is a valid sender for your provider.

## Database
Migration: `internal/database/migrations/001_create_waitlist_table.up.sql`  
Table: `waitlist` with UUID `id`, `name`, `email` (unique, case-insensitive), `created_at`.

## Testing
```bash
GOWORK=off go test ./...
```
(No tests yet; add unit tests around services/handlers as needed.)

## Deployment notes
- Dockerfile builds a minimal Alpine image running the compiled binary.
- Set `DATABASE_URL`, SMTP vars, and (optionally) `ADMIN_API_KEY` in your deployment env.
- Run migrations before/with deploy.

## SEO/front-end tips
- Allow CORS to your landing/waitlist domain.
- Ensure `/health` is wired into your uptime checks.
