#!/bin/bash
set -e

AWS_REGION="us-east-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

VERSION=${1:-latest}

echo "🔐 Authenticating to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Push task-service
echo "📦 Building task-service..."
cd ../services/task-service
docker build -t task-service:$VERSION .
docker tag task-service:$VERSION $ECR_REGISTRY/task-service:$VERSION
docker push $ECR_REGISTRY/task-service:$VERSION
cd ../..

# Kong doesn't need pushing - using official Kong image

echo "✅ task-service pushed!"