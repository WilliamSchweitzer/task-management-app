#!/bin/bash
set -e

VERSION=${1:-latest}
SCRIPT_DIR="$(dirname "$0")"

echo "🚀 Starting task-service deployment..."

# Step 1: Apply infrastructure
echo "📋 Step 1: Applying Terraform infrastructure..."
cd "$SCRIPT_DIR/../terraform"
terraform apply -auto-approve

# Step 2: Build and push images
echo "🐳 Step 2: Building and pushing Docker images..."
cd ../scripts
cd "$SCRIPT_DIR"
./push-task-service.sh $VERSION

# Step 3: Force ECS deployment
echo "🔄 Step 3: Updating task-service ECS..."
aws ecs update-service --cluster task-management-cluster --service task-service --force-new-deployment --region us-east-1

# Step 4: Wait for services to be stable
echo "⏳ Step 4: Waiting for task-serviceservice to stabilize..."
aws ecs wait services-stable --cluster task-management-cluster --services task-service --region us-east-1

# Step 5: Configure Kong - TODO
echo "🔧 Step 5: Configuring Kong routes..."
./configure-kong.sh

echo "✅ Deployment complete!"
cd "$SCRIPT_DIR/../terraform"
echo "🌐 Access your app at: http://$(terraform output -raw nlb_dns_name)"