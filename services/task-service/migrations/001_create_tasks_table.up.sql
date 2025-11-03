-- Enable UUID extension (if not already enabled)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tasks schema
CREATE SCHEMA IF NOT EXISTS tasks;

-- Create tasks table
CREATE TABLE tasks.tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    priority VARCHAR(50) DEFAULT 'medium',
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_status CHECK (status IN ('todo', 'in-progress', 'done')),
    CONSTRAINT chk_priority CHECK (priority IN ('low', 'medium', 'high'))
);

-- Indexes
CREATE INDEX idx_tasks_user_id ON tasks.tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks.tasks(status);
CREATE INDEX idx_tasks_user_status ON tasks.tasks(user_id, status);
CREATE INDEX idx_tasks_due_date ON tasks.tasks(due_date) WHERE due_date IS NOT NULL;

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_tasks_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks.tasks
    FOR EACH ROW EXECUTE FUNCTION update_tasks_updated_at_column();
