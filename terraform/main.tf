terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
  
  backend "s3" {
    bucket         = "task-management-s3-bucket"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "task-management-state"
  }
}

provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.16.0"  # Changed from 2.77.0 to 5.16.0

  name                 = "taskmanagement"
  cidr                 = "10.0.0.0/16"
  azs                  = data.aws_availability_zones.available.names
  public_subnets       = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_db_subnet_group" "taskmanagement" {
  name       = "taskmanagement"
  subnet_ids = module.vpc.public_subnets

  tags = {
    Name = "Taskmanagement"
  }
}

resource "aws_db_instance" "taskmanagement" {
  identifier             = "taskmanagement"
  instance_class         = "db.t3.micro"
  allocated_storage      = 5
  engine                 = "postgres"
  engine_version         = "14"
  username               = "postgres"
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.taskmanagement.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.taskmanagement.name
  publicly_accessible    = true
  skip_final_snapshot    = true
}

resource "aws_security_group" "rds" {
  name   = "taskmanagement_rds"
  vpc_id = module.vpc.vpc_id

  # Allow from ECS tasks
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_tasks.id]  # Reference ECS security group
  }

  # Public access IS needed for testing/development purposes
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "taskmanagement_rds"
  }
}

resource "aws_db_parameter_group" "taskmanagement" {
  name   = "taskmanagement"
  family = "postgres14"

  parameter {
    name  = "log_connections"
    value = "1"
  }
}