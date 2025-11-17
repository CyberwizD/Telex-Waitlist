# Frontend Integration Guide

This guide is for frontend developers integrating the Telex Waitlist API.

## Base URLs
- Local (docker-compose): `http://localhost:8080`
- Production: set in your frontend env (e.g., `VITE_API_BASE_URL=https://api.example.com`)

All endpoints are prefixed with `/api/v1` unless noted.

## Endpoints

### Health
- `GET /health`
- Success: `200 {"status":"ok"}` (no authentication)

### Submit to waitlist (public)
- `POST /api/v1/waitlist`
- Headers: `Content-Type: application/json`
- Body:
  ```json
  { "name": "Jane Doe", "email": "jane@example.com" }
  ```
- Success: `201` with payload:
  ```json
  {
    "data": {
      "id": "uuid",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "created_at": "2025-11-17T00:00:00Z"
    }
  }
  ```
- Errors:
  - `400` invalid payload (missing/invalid name or email)
  - `409` email already exists on the waitlist

### List waitlist entries (admin-only)
- `GET /api/v1/waitlist?limit=50&offset=0`
- Header: `X-Admin-Token: <ADMIN_API_KEY>`
- Success: `200` with:
  ```json
  {
    "data": [ { "id": "...", "name": "...", "email": "...", "created_at": "..." } ],
    "total": 120,
    "limit": 50,
    "offset": 0
  }
  ```
- Errors:
  - `401` invalid/missing admin token
  - `403` listing disabled (server not configured with `ADMIN_API_KEY`)

## Validation rules
- `name`: required, non-empty string.
- `email`: required, valid email format; case-insensitive unique.

## Email behavior
- On successful submission, the backend triggers a thank-you email (if `EMAIL_ENABLED=true` and SMTP is configured).
- Nothing is required from the frontend beyond calling the submit endpoint.

## CORS
- The backend allows origins configured via `ALLOWED_ORIGINS`. Leaving it empty allows any origin (development).
- Ensure your frontend domain is included in `ALLOWED_ORIGINS` in production.

## Example usage (fetch)
```js
const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

async function joinWaitlist(name, email) {
  const res = await fetch(`${apiBase}/api/v1/waitlist`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email }),
  });
  if (res.status === 409) throw new Error('This email is already on the waitlist');
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Failed to join waitlist');
  }
  const data = await res.json();
  return data.data;
}
```

## Local dev workflow for FE
1) Backend: `docker-compose up --build` (or `go run ./...` with a local Postgres).
2) Run migration: `psql "postgres://postgres:postgres@localhost:5432/telex_waitlist?sslmode=disable" -f internal/database/migrations/001_create_waitlist_table.up.sql`.
3) Point the frontend to `http://localhost:8080` via your env (e.g., `VITE_API_BASE_URL`).
4) Submit test entries with the endpoint above; check DB or enable admin listing with `ADMIN_API_KEY` if you need to fetch them.

## Error payloads
Errors are returned as:
```json
{ "error": "human-readable message" }
```

## Production checklist
- Confirm `ALLOWED_ORIGINS` includes the frontend origin.
- Set `ADMIN_API_KEY` if the FE needs a listing view; otherwise, leave empty to disable.
- Ensure SMTP environment variables are configured if thank-you emails should be sent.
- Run migrations before serving traffic.
