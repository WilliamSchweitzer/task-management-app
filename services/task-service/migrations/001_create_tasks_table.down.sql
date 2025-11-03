DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks.tasks;
DROP FUNCTION IF EXISTS update_tasks_updated_at_column();
DROP TABLE IF EXISTS tasks.tasks;
DROP SCHEMA IF EXISTS tasks CASCADE;
