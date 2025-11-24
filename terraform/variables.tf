variable "region" {
  default     = "us-east-1"
  description = "AWS region"
}

variable "db_password" {
  description = "RDS root user password"
  type        = string
  sensitive   = true
}

variable "kong_db_password" {
  description = "Kong database password"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret key for signing tokens"
  type        = string
  sensitive   = true
}