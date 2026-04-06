-- Milestone 2 foundation: auth schema and core user table.
-- Source of truth: docs/Technical_Architecture.md

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  avatar_url TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_users_email_not_blank CHECK (btrim(email) <> ''),
  CONSTRAINT chk_users_name_not_blank CHECK (btrim(name) <> ''),
  CONSTRAINT chk_users_name_len CHECK (char_length(name) BETWEEN 1 AND 255)
);

CREATE UNIQUE INDEX uq_users_lower_email ON users (LOWER(email));

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
