#!/bin/bash
set -e

# Navigate to terraform directory to get outputs
cd "$(dirname "$0")/../terraform"

NLB_DNS=$(terraform output -raw nlb_dns_name)
KONG_ADMIN_URL="http://$NLB_DNS:8001"

echo "🔧 Configuring Kong routes..."

# Wait for Kong to be ready
echo "⏳ Waiting for Kong Admin API..."
until curl -s $KONG_ADMIN_URL/status > /dev/null 2>&1; do
  echo "Waiting for Kong..."
  sleep 5
done

# Add auth-service
echo "Adding auth-service route..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=auth-service" \
  --data "url=http://auth-service.taskmanagement.local:8080"

curl -s -X POST $KONG_ADMIN_URL/services/auth-service/routes \
  --data "paths[]=/auth" \
  --data "paths[]=/api/auth" \
  --data "strip_path=false"

# Add task-service
echo "Adding task-service route..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=task-service" \
  --data "url=http://task-service.taskmanagement.local:8081"

curl -s -X POST $KONG_ADMIN_URL/services/task-service/routes \
  --data "paths[]=/tasks" \
  --data "paths[]=/api/tasks" \
  --data "strip_path=false"

# Add Cors plugin
echo "Adding CORS plugin to services..."

EXISTING_CORS=$(curl -s $KONG_ADMIN_URL/plugins | jq -r '.data[] | select(.name=="cors") | .id')
if [ -n "$EXISTING_CORS" ]; then
  echo "Removing existing CORS plugin..."
  curl -s -X DELETE $KONG_ADMIN_URL/plugins/$EXISTING_CORS
fi

CORS_ORIGINS=${CORS_ORIGINS:-'["http://localhost:3000", "https://task-management.wschweitzer.com"]'}

curl -s -X POST $KONG_ADMIN_URL/plugins \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"cors\",
    \"config\": {
      \"origins\": $CORS_ORIGINS,
      \"methods\": [\"GET\", \"POST\", \"PUT\", \"DELETE\", \"PATCH\", \"OPTIONS\"],
      \"headers\": [\"Accept\", \"Authorization\", \"Content-Type\"],
      \"credentials\": true,
      \"max_age\": 3600
    }
  }" | jq . 2>/dev/null || echo ""

echo "✅ Kong routes configured!"