DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP INDEX IF EXISTS uq_users_lower_email;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS set_updated_at();
