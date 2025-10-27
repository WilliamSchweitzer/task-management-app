# Auth Service

Authentication and authorization microservice for the Task Management App.

## Features

- User registration (signup)
- User authentication (login)
- JWT token generation and validation
- Refresh token mechanism
- Password hashing with bcrypt
- Token verification endpoint

## API Endpoints

### Health Check
```
GET /health
```

### Authentication Endpoints

#### Signup
```
POST /auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

#### Login
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### Refresh Token
```
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Verify Token
```
GET /auth/verify
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

#### Logout
```
POST /auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

## Environment Variables

```bash
AUTH_SERVICE_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanagement
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h
LOG_LEVEL=info
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP
);
```

## Running Locally

### Prerequisites
- Go 1.21+
- PostgreSQL 14+

### Setup
```bash
# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go up

# Run the service
go run cmd/main.go
```

### Running with Docker
```bash
# Build image
docker build -t auth-service .

# Run container
docker run -p 8080:8080 --env-file .env auth-service
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Project Structure

```
auth-service/
├── cmd/
│   ├── main.go           # Application entry point
│   └── migrate/          # Database migration tool
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data models
│   ├── repository/      # Database layer
│   └── service/         # Business logic
├── migrations/          # SQL migration files
├── Dockerfile
├── go.mod
└── go.sum
```

## TODO

- [ ] Implement user registration with email validation
- [ ] Add password reset functionality
- [ ] Implement rate limiting for login attempts
- [ ] Add OAuth 2.0 integration (Google, GitHub)
- [ ] Add email verification
- [ ] Implement 2FA (Two-Factor Authentication)
- [ ] Add user profile management endpoints
- [ ] Implement account deletion
