DROP INDEX IF EXISTS uq_users_google_subject;

ALTER TABLE users
  DROP COLUMN IF EXISTS google_subject;

UPDATE users
SET password_hash = 'oauth-only-account'
WHERE password_hash IS NULL;

ALTER TABLE users
  ALTER COLUMN password_hash SET NOT NULL;
