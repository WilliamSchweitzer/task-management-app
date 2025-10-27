# General Variables
variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "task-management"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }
}

variable "owner" {
  description = "Owner of the resources"
  type        = string
  default     = "william-schweitzer"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

# Networking Variables
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

# Database Variables
variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"  # Free tier eligible
}

variable "db_allocated_storage" {
  description = "Allocated storage for RDS in GB"
  type        = number
  default     = 20  # Free tier limit
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "15.4"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "taskmanagement"
}

variable "db_username" {
  description = "Database master username"
  type        = string
  default     = "dbadmin"
  sensitive   = true
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

# ECS Service Variables
variable "auth_service_image" {
  description = "Docker image for auth service"
  type        = string
  default     = "auth-service:latest"
}

variable "task_service_image" {
  description = "Docker image for task service"
  type        = string
  default     = "task-service:latest"
}

variable "kong_image" {
  description = "Docker image for Kong"
  type        = string
  default     = "kong:3.4"
}

variable "auth_service_port" {
  description = "Port for auth service"
  type        = number
  default     = 8080
}

variable "task_service_port" {
  description = "Port for task service"
  type        = number
  default     = 8081
}

variable "kong_proxy_port" {
  description = "Kong proxy port"
  type        = number
  default     = 8000
}

# ECS Task Resources (Free tier: 600 vCPU hours/month, 1200 GB hours/month)
variable "auth_service_cpu" {
  description = "CPU units for auth service (1024 = 1 vCPU)"
  type        = number
  default     = 256  # 0.25 vCPU
}

variable "auth_service_memory" {
  description = "Memory for auth service in MB"
  type        = number
  default     = 512  # 0.5 GB
}

variable "task_service_cpu" {
  description = "CPU units for task service (1024 = 1 vCPU)"
  type        = number
  default     = 256  # 0.25 vCPU
}

variable "task_service_memory" {
  description = "Memory for task service in MB"
  type        = number
  default     = 512  # 0.5 GB
}

variable "kong_cpu" {
  description = "CPU units for Kong (1024 = 1 vCPU)"
  type        = number
  default     = 256  # 0.25 vCPU
}

variable "kong_memory" {
  description = "Memory for Kong in MB"
  type        = number
  default     = 512  # 0.5 GB
}

# Security Variables
variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

variable "certificate_arn" {
  description = "ARN of ACM certificate for HTTPS (optional)"
  type        = string
  default     = ""
}

# Monitoring Variables
variable "enable_detailed_monitoring" {
  description = "Enable detailed CloudWatch monitoring"
  type        = bool
  default     = false  # Set to true for production
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}
