# Define all services including Kong
locals {
  services = ["auth-service", "task-service", "kong-gateway"]
}

# Create ECR repository for each service
resource "aws_ecr_repository" "services" {
  for_each = toset(local.services)

  name                 = each.key
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = each.key
    Environment = "production"
    Project     = "task-management"
  }
}

# Lifecycle policy for all repositories
resource "aws_ecr_lifecycle_policy" "services" {
  for_each = aws_ecr_repository.services

  repository = each.value.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 tagged images"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["v"]
          countType     = "imageCountMoreThan"
          countNumber   = 10
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Remove untagged images after 7 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Outputs
output "ecr_repository_urls" {
  description = "ECR repository URLs for all services"
  value = {
    for service, repo in aws_ecr_repository.services :
    service => repo.repository_url
  }
}

output "ecr_auth_service_url" {
  description = "ECR repository URL for auth-service"
  value       = aws_ecr_repository.services["auth-service"].repository_url
}

output "ecr_task_service_url" {
  description = "ECR repository URL for task-service"
  value       = aws_ecr_repository.services["task-service"].repository_url
}

output "ecr_kong_gateway_url" {
  description = "ECR repository URL for kong-gateway"
  value       = aws_ecr_repository.services["kong-gateway"].repository_url
}

output "ecr_registry_id" {
  description = "ECR registry ID"
  value       = aws_ecr_repository.services["auth-service"].registry_id
}