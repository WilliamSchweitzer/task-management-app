-- Initialize database with separate schemas for auth and tasks

-- Create schemas
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS tasks;

-- Set search path
SET search_path TO auth, tasks, public;

-- Auth schema tables will be created by migrations
-- Tasks schema tables will be created by migrations

-- Grant permissions (adjust as needed for production)
GRANT ALL PRIVILEGES ON SCHEMA auth TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA tasks TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA tasks TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA tasks TO postgres;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a read-only user for reporting (optional)
-- CREATE USER readonly_user WITH PASSWORD 'change_me_in_production';
-- GRANT CONNECT ON DATABASE taskmanagement TO readonly_user;
-- GRANT USAGE ON SCHEMA auth, tasks TO readonly_user;
-- GRANT SELECT ON ALL TABLES IN SCHEMA auth, tasks TO readonly_user;
