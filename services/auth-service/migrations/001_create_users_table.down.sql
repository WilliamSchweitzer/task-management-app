DROP TRIGGER IF EXISTS update_users_updated_at ON auth.users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth CASCADE;
