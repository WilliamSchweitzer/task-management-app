terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Owner       = var.owner
    }
  }
}

# Data sources
data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available" {
  state = "available"
}

# Networking Module
module "networking" {
  source = "./modules/networking"

  project_name        = var.project_name
  environment         = var.environment
  vpc_cidr            = var.vpc_cidr
  availability_zones  = slice(data.aws_availability_zones.available.names, 0, 2)
  public_subnet_cidrs = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}

# Database Module
module "database" {
  source = "./modules/database"

  project_name           = var.project_name
  environment            = var.environment
  vpc_id                 = module.networking.vpc_id
  private_subnet_ids     = module.networking.private_subnet_ids
  db_security_group_id   = module.networking.db_security_group_id
  db_instance_class      = var.db_instance_class
  db_allocated_storage   = var.db_allocated_storage
  db_engine_version      = var.db_engine_version
  db_name                = var.db_name
  db_username            = var.db_username
  db_password            = var.db_password
}

# ECR Module (Container Registry)
module "ecr" {
  source = "./modules/ecr"

  project_name = var.project_name
  environment  = var.environment
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"

  project_name                = var.project_name
  environment                 = var.environment
  vpc_id                      = module.networking.vpc_id
  public_subnet_ids           = module.networking.public_subnet_ids
  ecs_security_group_id       = module.networking.ecs_security_group_id
  
  # Service configurations
  auth_service_image          = var.auth_service_image
  task_service_image          = var.task_service_image
  kong_image                  = var.kong_image
  
  auth_service_port           = var.auth_service_port
  task_service_port           = var.task_service_port
  kong_proxy_port             = var.kong_proxy_port
  
  auth_service_cpu            = var.auth_service_cpu
  auth_service_memory         = var.auth_service_memory
  task_service_cpu            = var.task_service_cpu
  task_service_memory         = var.task_service_memory
  kong_cpu                    = var.kong_cpu
  kong_memory                 = var.kong_memory
  
  # Database connection
  db_host                     = module.database.db_endpoint
  db_port                     = module.database.db_port
  db_name                     = var.db_name
  db_username                 = var.db_username
  db_password                 = var.db_password
  
  # JWT secret
  jwt_secret                  = var.jwt_secret
}

# Application Load Balancer Module
module "alb" {
  source = "./modules/alb"

  project_name          = var.project_name
  environment           = var.environment
  vpc_id                = module.networking.vpc_id
  public_subnet_ids     = module.networking.public_subnet_ids
  alb_security_group_id = module.networking.alb_security_group_id
  
  # Target groups from ECS module
  kong_target_group_arn = module.ecs.kong_target_group_arn
  
  # Certificate ARN (if using HTTPS)
  certificate_arn       = var.certificate_arn
}

# Monitoring Module
module "monitoring" {
  source = "./modules/monitoring"

  project_name = var.project_name
  environment  = var.environment
  
  # ECS resources to monitor
  ecs_cluster_name     = module.ecs.cluster_name
  auth_service_name    = module.ecs.auth_service_name
  task_service_name    = module.ecs.task_service_name
  kong_service_name    = module.ecs.kong_service_name
  
  # ALB to monitor
  alb_arn_suffix       = module.alb.alb_arn_suffix
  
  # Database to monitor
  db_instance_id       = module.database.db_instance_id
}
