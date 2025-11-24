# Task Management App

A full-stack task management application showcasing microservices architecture, infrastructure as code, and modern DevOps practices.

## 📚 Primary Learning Goals

- Infrastrucute setup and management via Terraform (AWS, Kong)
- Backend development with Go, PostgreSQL (REST, JWT Authentication) 

## Kong Note

**Current State (Development):**
Microservices are currently publicly accessible for easier development and testing.

**Production Recommendation:**
In a production environment, auth-service and task-service should be:
1. Deployed in private subnets with no public IPs
2. Accessible only through Kong Gateway
3. Security groups restricting traffic to Kong only
4. Using VPC endpoints for AWS services

## 🏗️ Architecture

This project demonstrates a production-ready microservices architecture with:

- **Infrastructure as Code**: Terraform for AWS resource provisioning
- **API Gateway**: Kong for authentication, rate limiting, and routing
- **Microservices**: Go-based auth and task services
- **Database**: PostgreSQL with proper schema separation
- **Frontend**: React/Next.js with TypeScript
- **CI/CD**: GitHub Actions for automated testing and deployment

### Architecture Diagram

```
User Browser
    ↓
Route 53 (DNS)
    ↓
Application Load Balancer (ALB)
    ↓
Kong API Gateway (ECS Fargate)
    ├─ JWT Authentication Plugin
    ├─ Rate Limiting Plugin
    └─ CORS Plugin
    ↓
    ├─→ Auth Service (Go, ECS Fargate)
    └─→ Task Service (Go, ECS Fargate)
            ↓
    RDS PostgreSQL (Multi-schema)
```

## 🛠️ Tech Stack

### Infrastructure
- **Terraform** - Infrastructure as Code
- **AWS ECS Fargate** - Serverless container orchestration
- **AWS RDS PostgreSQL** - Managed database
- **AWS ALB** - Load balancing and SSL termination
- **AWS ECR** - Container registry
- **AWS CloudWatch** - Logging and monitoring

### Backend
- **Go 1.21+** - Primary language
- **Chi/Gin** - HTTP router framework
- **GORM** - ORM for database operations
- **golang-migrate** - Database migrations
- **JWT** - Authentication tokens

### API Gateway
- **Kong** - API Gateway
  - JWT Plugin
  - Rate Limiting Plugin
  - CORS Plugin
  - Request Transformer Plugin

### Frontend
- **React 18** - UI framework
- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **TanStack Query** - Server state management

### DevOps
- **Docker** - Containerization
- **GitHub Actions** - CI/CD
- **docker-compose** - Local development

## 📋 Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **Docker** and Docker Compose
- **Terraform** 1.5 or higher
- **AWS CLI** configured with credentials
- **PostgreSQL** 14+ (for local development)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/WilliamSchweitzer/task-management-app.git
cd task-management-app
```

### 2. Local Development Setup

```bash
# Copy environment variables
cp .env.example .env

# Start all services with docker-compose
docker-compose up -d

# Run database migrations
make migrate-up

# The app will be available at:
# - Frontend: http://localhost:3000
# - Kong Gateway: http://localhost:8000
# - Kong Admin API: http://localhost:8001
```

### 3. Development Workflow

```bash
# Start individual services for development
cd services/auth-service
go run cmd/main.go

cd services/task-service
go run cmd/main.go

cd frontend
npm run dev
```

## 🏗️ Infrastructure Deployment

### Prerequisites

1. AWS account with appropriate permissions
2. Terraform installed locally
3. AWS CLI configured

### Deploy to AWS

```bash
# Initialize Terraform
cd terraform
terraform init

# Plan infrastructure changes
terraform plan -var-file="environments/dev/terraform.tfvars"

# Apply infrastructure
terraform apply -var-file="environments/dev/terraform.tfvars"
```

## 📖 API Documentation

### Auth Service Endpoints

```
POST   /auth/signup     - Create new user account
POST   /auth/login      - Login and receive JWT tokens
POST   /auth/refresh    - Refresh access token
GET    /auth/verify     - Verify token validity
POST   /auth/logout     - Logout (invalidate refresh token)
```

### Task Service Endpoints

```
GET    /tasks           - Get all tasks for authenticated user
POST   /tasks           - Create a new task
GET    /tasks/:id       - Get specific task
PUT    /tasks/:id       - Update task
DELETE /tasks/:id       - Delete task
PATCH  /tasks/:id/complete - Mark task as complete
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific service tests
cd services/auth-service
go test ./...

# Run frontend tests
cd frontend
npm test
```

## 📁 Project Structure

```
task-management-app/
├── terraform/              # Infrastructure as Code
│   ├── modules/           # Reusable Terraform modules
│   └── environments/      # Environment-specific configs
├── services/              # Backend microservices
│   ├── auth-service/     # Authentication service
│   └── task-service/     # Task management service
├── frontend/             # React frontend application
├── kong/                 # Kong API Gateway configuration
├── scripts/              # Utility scripts
├── docs/                 # Documentation
└── .github/workflows/    # CI/CD pipelines
```

## 🔧 Development Commands

```bash
# Infrastructure
make tf-init          # Initialize Terraform
make tf-plan          # Plan infrastructure changes
make tf-apply         # Apply infrastructure changes
make tf-destroy       # Destroy infrastructure

# Docker
make docker-build     # Build all Docker images
make docker-push      # Push images to ECR
make docker-up        # Start docker-compose stack
make docker-down      # Stop docker-compose stack

# Database
make migrate-up       # Run all migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration

# Testing
make test             # Run all tests
make test-coverage    # Run tests with coverage
make lint             # Run linters

# Development
make dev-auth         # Run auth service locally
make dev-task         # Run task service locally
make dev-frontend     # Run frontend locally
```

## 🔐 Security Features

- JWT-based authentication with refresh tokens
- Password hashing with bcrypt
- SQL injection prevention via parameterized queries
- CORS configuration
- Rate limiting to prevent abuse
- Security groups for network isolation
- Secrets management via AWS Secrets Manager
- HTTPS enforcement

## 📊 Monitoring & Observability

- CloudWatch Logs for all services
- CloudWatch Metrics for resource utilization
- ALB access logs
- Kong request/response logging
- Database query logging

## 🌟 Key Features

### For Interviews
This project demonstrates:
- ✅ Microservices architecture
- ✅ Infrastructure as Code (Terraform)
- ✅ API Gateway patterns (Kong)
- ✅ Container orchestration (ECS Fargate)
- ✅ Database design and migrations
- ✅ JWT authentication
- ✅ RESTful API design
- ✅ CI/CD pipelines
- ✅ Security best practices
- ✅ Cost optimization (public subnets, minimal resources)

### Technical Highlights
- Clean architecture with separation of concerns
- Proper error handling and logging
- Database migrations for version control
- Environment-based configuration
- Comprehensive testing
- API documentation with OpenAPI/Swagger

## 🚧 Roadmap

- [ ] Add OAuth 2.0 login (Google, GitHub)
- [ ] Implement WebSocket for real-time updates
- [ ] Add task categories and tags
- [ ] Implement task sharing between users
- [ ] Add email notifications
- [ ] Create mobile app (React Native)
- [ ] Add GraphQL API option
- [ ] Implement caching layer (Redis)
- [ ] Add comprehensive metrics dashboard

## 🤝 Contributing

This is a personal portfolio project, but suggestions and feedback are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📝 License

MIT License - See [LICENSE](LICENSE) file for details

## 👤 Author

**William Schweitzer**

- Website: [wschweitzer.com](https://wschweitzer.com)
- GitHub: [@WilliamSchweitzer](https://github.com/WilliamSchweitzer)

## 🙏 Acknowledgments

- Scaffolded with Claude AI
- Functionality to be implemented alongside Claude AI
- Infrastructure managed and setup by hand
- Built as a portfolio project to demonstrate full-stack and DevOps skills
- Architecture inspired by production microservices patterns
- Designed with job interview showcasing in mind

---

**Note**: This project is optimized for learning and demonstration. In a production environment with sensitive data, additional security measures would be recommended (private subnets with NAT Gateway, WAF, additional monitoring, etc.).
