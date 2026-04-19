ALTER TABLE users
  ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE users
  ADD COLUMN google_subject VARCHAR(255) NULL;

CREATE UNIQUE INDEX uq_users_google_subject
  ON users (google_subject)
  WHERE google_subject IS NOT NULL;
