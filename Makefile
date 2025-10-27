.PHONY: help docker-up docker-down docker-build docker-push tf-init tf-plan tf-apply tf-destroy migrate-up migrate-down migrate-create test test-coverage lint dev-auth dev-task dev-frontend clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker commands
docker-up: ## Start all services with docker compose
	docker compose up -d

docker-down: ## Stop all services
	docker compose down

docker-build: ## Build all Docker images
	docker compose build

docker-push: ## Push images to ECR (requires AWS login)
	@echo "Building and pushing images to ECR..."
	@./scripts/push-to-ecr.sh

docker-logs: ## View logs from all containers
	docker compose logs -f

docker-clean: ## Remove all containers, volumes, and images
	docker compose down -v --rmi all

# Terraform commands
tf-init: ## Initialize Terraform
	cd terraform && terraform init

tf-plan: ## Plan Terraform changes (dev environment)
	cd terraform && terraform plan -var-file="environments/dev/terraform.tfvars"

tf-apply: ## Apply Terraform changes (dev environment)
	cd terraform && terraform apply -var-file="environments/dev/terraform.tfvars"

tf-destroy: ## Destroy Terraform infrastructure (dev environment)
	cd terraform && terraform destroy -var-file="environments/dev/terraform.tfvars"

tf-fmt: ## Format Terraform files
	cd terraform && terraform fmt -recursive

tf-validate: ## Validate Terraform configuration
	cd terraform && terraform validate

# Database migration commands
migrate-up: ## Run all database migrations
	@echo "Running migrations for auth service..."
	cd services/auth-service && go run cmd/migrate/main.go up
	@echo "Running migrations for task service..."
	cd services/task-service && go run cmd/migrate/main.go up

migrate-down: ## Rollback last migration
	@echo "Rolling back migrations..."
	cd services/auth-service && go run cmd/migrate/main.go down
	cd services/task-service && go run cmd/migrate/main.go down

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_users_table SERVICE=auth)
	@if [ -z "$(NAME)" ] || [ -z "$(SERVICE)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name SERVICE=auth|task"; \
		exit 1; \
	fi
	@cd services/$(SERVICE)-service && \
		migrate create -ext sql -dir migrations -seq $(NAME)

# Testing commands
test: ## Run all tests
	@echo "Running auth service tests..."
	cd services/auth-service && go test ./... -v
	@echo "Running task service tests..."
	cd services/task-service && go test ./... -v
	@echo "Running frontend tests..."
	cd frontend && npm test

test-coverage: ## Run tests with coverage
	@echo "Running auth service tests with coverage..."
	cd services/auth-service && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@echo "Running task service tests with coverage..."
	cd services/task-service && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# Linting commands
lint: ## Run linters for all services
	@echo "Linting auth service..."
	cd services/auth-service && golangci-lint run
	@echo "Linting task service..."
	cd services/task-service && golangci-lint run
	@echo "Linting frontend..."
	cd frontend && npm run lint

# Development commands
dev-auth: ## Run auth service locally
	cd services/auth-service && go run cmd/main.go

dev-task: ## Run task service locally
	cd services/task-service && go run cmd/main.go

dev-frontend: ## Run frontend locally
	cd frontend && npm run dev

dev-kong: ## Configure Kong for local development
	@./scripts/setup-kong.sh

# Setup commands
setup: ## Initial project setup
	@echo "Setting up project..."
	cp .env.example .env
	@echo "Installing Go dependencies..."
	cd services/auth-service && go mod download
	cd services/task-service && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Setup complete! Run 'make docker-up' to start services."

# Cleanup commands
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning build artifacts..."
	find . -name "*.out" -delete
	find . -name "*.test" -delete
	find . -name "coverage.html" -delete
	cd services/auth-service && go clean
	cd services/task-service && go clean
	cd frontend && rm -rf .next node_modules

# AWS commands
aws-login: ## Login to AWS ECR
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.us-east-1.amazonaws.com

# Database commands
db-shell: ## Connect to PostgreSQL database
	docker compose exec postgres psql -U postgres -d taskmanagement

db-reset: ## Reset database (WARNING: destroys all data)
	docker compose down -v
	docker compose up -d postgres
	@echo "Waiting for database to be ready..."
	sleep 5
	make migrate-up

# Monitoring commands
logs-auth: ## View auth service logs
	docker compose logs -f auth-service

logs-task: ## View task service logs
	docker compose logs -f task-service

logs-kong: ## View Kong logs
	docker compose logs -f kong

logs-all: ## View all service logs
	docker compose logs -f
