CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS waitlist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT email_format CHECK (position('@' in email) > 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_waitlist_email_unique ON waitlist (lower(email));
