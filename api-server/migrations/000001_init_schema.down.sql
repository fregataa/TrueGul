DROP TRIGGER IF EXISTS update_analyses_updated_at ON analyses;
DROP TRIGGER IF EXISTS update_writings_updated_at ON writings;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP TABLE IF EXISTS analysis_logs;
DROP TABLE IF EXISTS analyses;
DROP TABLE IF EXISTS writings;
DROP TABLE IF EXISTS users;
