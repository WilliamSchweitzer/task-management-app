#!/bin/bash
set -e

AWS_REGION="us-east-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

VERSION=${1:-latest}

echo "🔐 Authenticating to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Push auth-service
echo "📦 Building auth-service..."
cd ../services/auth-service
docker build -t auth-service:$VERSION .
docker tag auth-service:$VERSION $ECR_REGISTRY/auth-service:$VERSION
docker push $ECR_REGISTRY/auth-service:$VERSION
cd ..

# Push task-service
echo "📦 Building task-service..."
cd task-service
docker build -t task-service:$VERSION .
docker tag task-service:$VERSION $ECR_REGISTRY/task-service:$VERSION
docker push $ECR_REGISTRY/task-service:$VERSION
cd ../..

# Kong doesn't need pushing - using official Kong image

echo "✅ All services pushed!"