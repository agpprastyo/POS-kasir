-- Drop trigger
DROP TRIGGER IF EXISTS trigger_set_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS set_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Drop table
DROP TABLE IF EXISTS users;

-- Drop enum type
DROP TYPE IF EXISTS user_role;