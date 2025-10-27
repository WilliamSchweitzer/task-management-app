#!/bin/bash

# Setup Kong for local development
# This script configures Kong routes and plugins

set -e

KONG_ADMIN_URL="${KONG_ADMIN_URL:-http://localhost:8001}"

echo "Setting up Kong configuration..."
echo "Kong Admin URL: $KONG_ADMIN_URL"

# Wait for Kong to be ready
echo "Waiting for Kong to be ready..."
until curl -s "${KONG_ADMIN_URL}" > /dev/null; do
    echo "Kong is not ready yet. Retrying in 5 seconds..."
    sleep 5
done

echo "Kong is ready!"

# Create Auth Service
echo "Creating auth-service..."
curl -i -X POST "${KONG_ADMIN_URL}/services" \
  --data name=auth-service \
  --data url='http://auth-service:8080'

# Create Auth Routes
echo "Creating auth routes..."
curl -i -X POST "${KONG_ADMIN_URL}/services/auth-service/routes" \
  --data 'paths[]=/auth' \
  --data 'strip_path=false' \
  --data 'methods[]=GET' \
  --data 'methods[]=POST' \
  --data 'methods[]=PUT' \
  --data 'methods[]=DELETE' \
  --data 'methods[]=OPTIONS'

# Add CORS plugin to auth routes
echo "Adding CORS plugin to auth routes..."
curl -i -X POST "${KONG_ADMIN_URL}/services/auth-service/plugins" \
  --data "name=cors" \
  --data "config.origins=http://localhost:3000" \
  --data "config.origins=http://localhost:8000" \
  --data "config.methods=GET" \
  --data "config.methods=POST" \
  --data "config.methods=PUT" \
  --data "config.methods=DELETE" \
  --data "config.methods=OPTIONS" \
  --data "config.credentials=true" \
  --data "config.max_age=3600"

# Add rate limiting to auth service
echo "Adding rate limiting to auth service..."
curl -i -X POST "${KONG_ADMIN_URL}/services/auth-service/plugins" \
  --data "name=rate-limiting" \
  --data "config.minute=100" \
  --data "config.hour=1000" \
  --data "config.policy=local"

# Create Task Service
echo "Creating task-service..."
curl -i -X POST "${KONG_ADMIN_URL}/services" \
  --data name=task-service \
  --data url='http://task-service:8081'

# Create Task Routes
echo "Creating task routes..."
curl -i -X POST "${KONG_ADMIN_URL}/services/task-service/routes" \
  --data 'paths[]=/tasks' \
  --data 'strip_path=false' \
  --data 'methods[]=GET' \
  --data 'methods[]=POST' \
  --data 'methods[]=PUT' \
  --data 'methods[]=DELETE' \
  --data 'methods[]=PATCH' \
  --data 'methods[]=OPTIONS'

# Add JWT plugin to task routes
echo "Adding JWT authentication to task routes..."
# Note: You'll need to configure JWT plugin with your JWT secret
# This is a placeholder - adjust based on your JWT implementation
curl -i -X POST "${KONG_ADMIN_URL}/services/task-service/plugins" \
  --data "name=jwt"

# Add CORS plugin to task routes
echo "Adding CORS plugin to task routes..."
curl -i -X POST "${KONG_ADMIN_URL}/services/task-service/plugins" \
  --data "name=cors" \
  --data "config.origins=http://localhost:3000" \
  --data "config.origins=http://localhost:8000" \
  --data "config.methods=GET" \
  --data "config.methods=POST" \
  --data "config.methods=PUT" \
  --data "config.methods=DELETE" \
  --data "config.methods=PATCH" \
  --data "config.methods=OPTIONS" \
  --data "config.credentials=true" \
  --data "config.max_age=3600"

# Add rate limiting to task service
echo "Adding rate limiting to task service..."
curl -i -X POST "${KONG_ADMIN_URL}/services/task-service/plugins" \
  --data "name=rate-limiting" \
  --data "config.minute=100" \
  --data "config.hour=5000" \
  --data "config.policy=local"

echo "Kong setup complete!"
echo ""
echo "Services:"
echo "  - Auth Service: ${KONG_ADMIN_URL/8001/8000}/auth"
echo "  - Task Service: ${KONG_ADMIN_URL/8001/8000}/tasks"
echo ""
echo "Kong Admin: ${KONG_ADMIN_URL}"
