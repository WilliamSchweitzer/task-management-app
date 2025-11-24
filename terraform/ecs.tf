# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "task-management-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"  # Cost saving
  }

  tags = {
    Name = "task-management-cluster"
  }
}

# Network Load Balancer (cheaper than ALB for Kong)
resource "aws_lb" "main" {
  name               = "task-management-nlb"
  internal           = false
  load_balancer_type = "network"
  subnets            = module.vpc.public_subnets

  tags = {
    Name = "task-management-nlb"
  }
}

# Target Group for Kong
resource "aws_lb_target_group" "kong" {
  name        = "kong-tg"
  port        = 8000
  protocol    = "TCP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"

  health_check {
    protocol            = "HTTP"
    path                = "/status"
    port                = "8001"
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 10
    interval            = 30
  }

  deregistration_delay = 30
}

# NLB Listener
resource "aws_lb_listener" "kong" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.kong.arn
  }
}

# ECS Task Execution Role
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Additional policy for Secrets Manager
resource "aws_iam_role_policy" "ecs_secrets_policy" {
  name = "ecs-secrets-policy"
  role = aws_iam_role.ecs_task_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          aws_secretsmanager_secret.db_password.arn,
          aws_secretsmanager_secret.kong_db_password.arn,
          aws_secretsmanager_secret.jwt_secret.arn
        ]
      }
    ]
  })
}

# Security Group for ECS Tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "task-management-ecs-tasks"
  description = "Allow inbound traffic for ECS tasks"
  vpc_id      = module.vpc.vpc_id

  # Allow traffic from NLB
  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow tasks to communicate with each other
  ingress {
    from_port = 0
    to_port   = 65535
    protocol  = "tcp"
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "task-management-ecs-tasks"
  }
}

# Update RDS security group to allow ECS tasks
resource "aws_security_group_rule" "rds_from_ecs" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.ecs_tasks.id
  security_group_id        = aws_security_group.rds.id
}

# CloudWatch Log Groups
resource "aws_cloudwatch_log_group" "auth_service" {
  name              = "/ecs/auth-service"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "task_service" {
  name              = "/ecs/task-service"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "kong" {
  name              = "/ecs/kong-gateway"
  retention_in_days = 7
}

# Secrets Manager
resource "aws_secretsmanager_secret" "db_password" {
  name = "task-management-db-password"
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

resource "aws_secretsmanager_secret" "kong_db_password" {
  name = "kong-db-password"
}

resource "aws_secretsmanager_secret_version" "kong_db_password" {
  secret_id     = aws_secretsmanager_secret.kong_db_password.id
  secret_string = var.kong_db_password
}

resource "aws_secretsmanager_secret" "jwt_secret" {
  name = "task-management-jwt-secret"
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = var.jwt_secret
}

# Service Discovery Namespace
resource "aws_service_discovery_private_dns_namespace" "main" {
  name = "taskmanagement.local"
  vpc  = module.vpc.vpc_id
}

resource "aws_service_discovery_service" "auth_service" {
  name = "auth-service"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

resource "aws_service_discovery_service" "task_service" {
  name = "task-service"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

# Kong Migration Task (One-time)
resource "aws_ecs_task_definition" "kong_migration" {
  family                   = "kong-migration"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name  = "kong-migration"
      image = "kong:3.4"
      
      command = ["kong", "migrations", "bootstrap"]

      environment = [
        {
          name  = "KONG_DATABASE"
          value = "postgres"
        },
        {
          name  = "KONG_PG_HOST"
          value = aws_db_instance.taskmanagement.address
        },
        {
          name  = "KONG_PG_PORT"
          value = "5432"
        },
        {
          name  = "KONG_PG_USER"
          value = "kong"
        },
        {
          name  = "KONG_PG_DATABASE"
          value = "kong"
        }
      ]

      secrets = [
        {
          name      = "KONG_PG_PASSWORD"
          valueFrom = aws_secretsmanager_secret.kong_db_password.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.kong.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "migration"
        }
      }
    }
  ])
}

# ECS Task Definition - Kong Gateway
resource "aws_ecs_task_definition" "kong" {
  family                   = "kong-gateway"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"  # Kong needs a bit more resources
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name  = "kong-gateway"
      image = "kong:3.4"
      
      portMappings = [
        {
          containerPort = 8000
          protocol      = "tcp"
        },
        {
          containerPort = 8001
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "KONG_DATABASE"
          value = "postgres"
        },
        {
          name  = "KONG_PG_HOST"
          value = aws_db_instance.taskmanagement.address
        },
        {
          name  = "KONG_PG_PORT"
          value = "5432"
        },
        {
          name  = "KONG_PG_USER"
          value = "kong"
        },
        {
          name  = "KONG_PG_DATABASE"
          value = "kong"
        },
        {
          name  = "KONG_PROXY_ACCESS_LOG"
          value = "/dev/stdout"
        },
        {
          name  = "KONG_ADMIN_ACCESS_LOG"
          value = "/dev/stdout"
        },
        {
          name  = "KONG_PROXY_ERROR_LOG"
          value = "/dev/stderr"
        },
        {
          name  = "KONG_ADMIN_ERROR_LOG"
          value = "/dev/stderr"
        },
        {
          name  = "KONG_ADMIN_LISTEN"
          value = "0.0.0.0:8001"
        }
      ]

      secrets = [
        {
          name      = "KONG_PG_PASSWORD"
          valueFrom = aws_secretsmanager_secret.kong_db_password.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.kong.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# ECS Task Definition - Auth Service
resource "aws_ecs_task_definition" "auth_service" {
  family                   = "auth-service"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name  = "auth-service"
      image = "${aws_ecr_repository.services["auth-service"].repository_url}:latest"
      
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "DB_HOST"
          value = aws_db_instance.taskmanagement.address
        },
        {
          name  = "DB_PORT"
          value = "5432"
        },
        {
          name  = "DB_USER"
          value = aws_db_instance.taskmanagement.username
        },
        {
          name  = "DB_NAME"
          value = "postgres"
        },
        {
          name  = "DB_SSLMODE"
          value = "require"
        },
        {
          name  = "PORT"
          value = "8080"
        },
        {
          name  = "ACCESS_TOKEN_EXPIRY"
          value = "15m"
        },
        {
          name  = "REFRESH_TOKEN_EXPIRY"
          value = "168h"
        }
      ]

      secrets = [
        {
          name      = "DB_PASSWORD"
          valueFrom = aws_secretsmanager_secret.db_password.arn
        },
        {
          name      = "JWT_SECRET"
          valueFrom = aws_secretsmanager_secret.jwt_secret.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.auth_service.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# ECS Task Definition - Task Service
resource "aws_ecs_task_definition" "task_service" {
  family                   = "task-service"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name  = "task-service"
      image = "${aws_ecr_repository.services["task-service"].repository_url}:latest"
      
      portMappings = [
        {
          containerPort = 8081
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "DB_HOST"
          value = aws_db_instance.taskmanagement.address
        },
        {
          name  = "DB_PORT"
          value = "5432"
        },
        {
          name  = "DB_USER"
          value = aws_db_instance.taskmanagement.username
        },
        {
          name  = "DB_NAME"
          value = "postgres"
        },
        {
          name  = "DB_SSLMODE"
          value = "require"
        },
        {
          name  = "PORT"
          value = "8081"
        }
      ]

      secrets = [
        {
          name      = "DB_PASSWORD"
          valueFrom = aws_secretsmanager_secret.db_password.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.task_service.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# ECS Service - Kong (with Fargate Spot)
resource "aws_ecs_service" "kong" {
  name            = "kong-gateway"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.kong.arn
  desired_count   = 1

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
    base              = 1
  }

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.kong.arn
    container_name   = "kong-gateway"
    container_port   = 8000
  }

  # Add this second load balancer block
  load_balancer {
    target_group_arn = aws_lb_target_group.kong_admin.arn
    container_name   = "kong-gateway"
    container_port   = 8001
  }

  depends_on = [
    aws_lb_listener.kong,
    aws_lb_listener.kong_admin,
    null_resource.kong_migration_runner
  ]
}

# ECS Service - Auth Service (with Fargate Spot)
resource "aws_ecs_service" "auth_service" {
  name            = "auth-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.auth_service.arn
  desired_count   = 1

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
    base              = 1
  }

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = true
  }

  service_registries {
    registry_arn = aws_service_discovery_service.auth_service.arn
  }
}

# ECS Service - Task Service (with Fargate Spot)
resource "aws_ecs_service" "task_service" {
  name            = "task-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.task_service.arn
  desired_count   = 1

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
    base              = 1
  }

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = true
  }

  service_registries {
    registry_arn = aws_service_discovery_service.task_service.arn
  }
}

# Add target group for Kong Admin API
resource "aws_lb_target_group" "kong_admin" {
  name        = "kong-admin-tg"
  port        = 8001
  protocol    = "TCP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"

  health_check {
    protocol            = "HTTP"
    path                = "/status"
    port                = "8001"
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 10
    interval            = 30
  }

  deregistration_delay = 30
}

# Add NLB listener for Kong Admin API
resource "aws_lb_listener" "kong_admin" {
  load_balancer_arn = aws_lb.main.arn
  port              = 8001
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.kong_admin.arn
  }
}

# Run Kong migrations
resource "null_resource" "kong_migration_runner" {
  depends_on = [
    aws_db_instance.taskmanagement,
    null_resource.create_kong_database
  ]

  provisioner "local-exec" {
    command = <<-EOT
      aws ecs run-task \
        --cluster ${aws_ecs_cluster.main.name} \
        --task-definition ${aws_ecs_task_definition.kong_migration.family} \
        --launch-type FARGATE \
        --network-configuration "awsvpcConfiguration={subnets=[${join(",", module.vpc.public_subnets)}],securityGroups=[${aws_security_group.ecs_tasks.id}],assignPublicIp=ENABLED}"
    EOT
  }
}

# Outputs
output "nlb_dns_name" {
  description = "NLB DNS name"
  value       = aws_lb.main.dns_name
}

output "kong_gateway_url" {
  description = "Kong Gateway URL"
  value       = "http://${aws_lb.main.dns_name}"
}

output "kong_admin_url" {
  description = "Kong Admin API URL (requires port forwarding or VPN)"
  value       = "Access Kong Admin at port 8001 via ECS task"
}