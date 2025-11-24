#!/bin/bash
set -e

SCRIPT_DIR="$(dirname "$0")"

# Get RDS connection details from Terraform
cd "$SCRIPT_DIR/../terraform"
RDS_HOST=$(terraform output -raw rds_hostname)
RDS_PORT=$(terraform output -raw rds_port)
RDS_USER=$(terraform output -raw rds_username)

echo "🗄️  Running database migrations..."

# Run task-service migrations in order
cd "$SCRIPT_DIR/../services/task-service/migrations"
for file in *.up.sql; do
  echo "Running $file..."
  psql -h $RDS_HOST -p $RDS_PORT -U $RDS_USER postgres -f "$file"
done

# Run auth-service migrations in order
cd "$SCRIPT_DIR/../../auth-service/migrations"
for file in *.up.sql; do
  echo "Running $file..."
  psql -h $RDS_HOST -p $RDS_PORT -U $RDS_USER postgres -f "$file"
done

echo "✅ Migrations completed!"