# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.networking.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of public subnets"
  value       = module.networking.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of private subnets"
  value       = module.networking.private_subnet_ids
}

# Database Outputs
output "db_endpoint" {
  description = "RDS instance endpoint"
  value       = module.database.db_endpoint
}

output "db_port" {
  description = "RDS instance port"
  value       = module.database.db_port
}

output "db_name" {
  description = "Database name"
  value       = module.database.db_name
}

# ECR Outputs
output "auth_service_repository_url" {
  description = "URL of auth service ECR repository"
  value       = module.ecr.auth_service_repository_url
}

output "task_service_repository_url" {
  description = "URL of task service ECR repository"
  value       = module.ecr.task_service_repository_url
}

output "kong_repository_url" {
  description = "URL of Kong ECR repository"
  value       = module.ecr.kong_repository_url
}

# ECS Outputs
output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "auth_service_name" {
  description = "Name of auth service"
  value       = module.ecs.auth_service_name
}

output "task_service_name" {
  description = "Name of task service"
  value       = module.ecs.task_service_name
}

output "kong_service_name" {
  description = "Name of Kong service"
  value       = module.ecs.kong_service_name
}

# ALB Outputs
output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.alb.alb_dns_name
}

output "alb_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = module.alb.alb_zone_id
}

output "alb_url" {
  description = "URL of the Application Load Balancer"
  value       = "http://${module.alb.alb_dns_name}"
}

# Monitoring Outputs
output "cloudwatch_log_groups" {
  description = "CloudWatch log groups created"
  value       = module.monitoring.log_group_names
}

# Connection String (for reference, not for production use)
output "database_connection_info" {
  description = "Database connection information"
  value = {
    host     = module.database.db_endpoint
    port     = module.database.db_port
    database = module.database.db_name
    # Note: username and password should be retrieved from AWS Secrets Manager
  }
  sensitive = true
}

# Useful Commands
output "useful_commands" {
  description = "Useful commands for working with the infrastructure"
  value = {
    ecr_login          = "aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
    ecs_list_services  = "aws ecs list-services --cluster ${module.ecs.cluster_name}"
    alb_url            = "http://${module.alb.alb_dns_name}"
    db_connect         = "psql -h ${module.database.db_endpoint} -p ${module.database.db_port} -U ${var.db_username} -d ${var.db_name}"
  }
}
