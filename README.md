# BrowersFC REST API

A production-ready REST API for managing a football league system. It is built with Go and organized around Clean Architecture and Hexagonal Architecture principles.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [Docker Deployment](#docker-deployment)
- [API Documentation](#api-documentation)
- [Database](#database)
- [Security Model](#security-model)
- [Project Structure](#project-structure)
- [Development](#development)
- [CI/CD](#cicd)

## Overview

BrowersFC REST API is a backend service for managing football league operations, including:

- **User management** - Authentication, authorization, and user profiles.
- **Team management** - Create, update, and manage teams.
- **Player management** - Manage player information and statistics.
- **Match management** - Schedule, update, and track matches.
- **Season management** - Organize leagues by season.
- **Lineup management** - Create and manage match lineups.
- **Article management** - Publish club news and articles.

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters) and clean code principles to keep business logic isolated from frameworks and infrastructure concerns.

```text
┌─────────────────────────────────────────────────────┐
│                   API Layer                         │
│                 (api/handler/)                      │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Domain Services Layer                  │
│         (internal/domain/service/)                  │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│                 Domain Models                       │
│                   (domain/)                         │
└──────────────────┬──────────────────────────────────┘
                   │
        ┌──────────┴──────────┐
        │                     │
┌───────▼────────┐   ┌────────▼────────┐
│  Repositories  │   │   DTOs/Mappers  │
│    (Ports)     │   │   (Adapters)    │
└───────┬────────┘   └────────┬────────┘
        │                     │
┌───────▼─────────────────────▼───────────────────────┐
│            Infrastructure / Persistence             │
│      (internal/infrastructure/persistence/)         │
└─────────────────────────────────────────────────────┘
```

### Design Principles

- **SOLID principles** - Clear separation of responsibilities and maintainable abstractions.
- **Dependency injection** - Services depend on interfaces, not concrete implementations.
- **Repository pattern** - Data access is abstracted through domain interfaces.
- **DTOs** - API request and response models are separated from domain entities.
- **Minimal exports** - Only required types and functions are exposed from packages.

## Prerequisites

Make sure you have the following installed:

- **Go** 1.23.0 or later
- **PostgreSQL** 15+
- **Docker** and **Docker Compose** (optional, for containerized deployment)
- **Git**

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/EdwinRincon/browersfc-api.git
cd browersfc-api
```

### 2. Set up environment variables

Create a local `.env` file:

```bash
cp .env.example .env
```

See the [Environment Variables](#environment-variables) section for details.

### 3. Install dependencies

```bash
go mod download
```

### 4. Configure secrets

Create a `secrets/` directory and add the required files:

```bash
mkdir -p secrets
echo "your-jwt-secret" > secrets/jwt_secret.txt
echo "your-oauth-client-secret" > secrets/oauth_client_secret.txt
echo "postgresql://user:password@localhost:5432/browersfc" > secrets/db_url.txt
```

### 5. Run database migrations

Migrations are applied automatically on startup. Make sure the database connection in `secrets/db_url.txt` is valid and reachable.

### 6. Start the application

```bash
go run ./cmd/browersfc/main.go
```

The API will be available at `http://localhost:5050`.

## Environment Variables

### Core configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `PORT` | string | `5050` | API server port |
| `GIN_MODE` | string | `debug` | Gin mode (`debug` or `release`) |
| `LOG_LEVEL` | string | `debug` in development / `info` in production | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| `LOG_FORMAT` | string | `text` in development / `json` in production | Log output format |

### Database configuration

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `DB_URL_FILE` | string | Yes | Path to the file containing the database URL |

Default connection example:

```text
postgresql://user:password@localhost:5432/browersfc
```

### Authentication and OAuth2

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `JWT_SECRET_FILE` | string | Yes | Path to the JWT signing secret file (minimum 32 characters) |
| `OAUTH_CLIENT_ID` | string | Yes | Google OAuth2 client ID |
| `OAUTH_CLIENT_SECRET_FILE` | string | Yes | Path to the Google OAuth2 client secret file |
| `OAUTH_REDIRECT_URL` | string | Yes | OAuth2 callback URL, for example `http://localhost:3000/auth/google/callback` |

### Security headers

Security headers are configured automatically based on the environment:

- **Content-Security-Policy** - Restricts resource loading.
- **X-Content-Type-Options** - Prevents MIME type sniffing.
- **X-Frame-Options** - Protects against clickjacking.
- **Strict-Transport-Security** - Enforced in production over HTTPS.

## Running the Application

### Local development

```bash
# Standard run
go run ./cmd/browersfc/main.go

# Run with database seeding (development only)
go run ./cmd/browersfc/main.go -seed
```

### Hot reload

Install `air` for live reload during development:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Build the binary

```bash
# Build for the current platform
go build -o browersfc ./cmd/browersfc

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o browersfc ./cmd/browersfc

# Run the binary
./browersfc
```

## Docker Deployment

### Docker Compose

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

### Multi-stage build

The `Dockerfile` uses a multi-stage build:

1. **Builder stage** - Compiles the Go application.
2. **Runtime stage** - Ships only the compiled binary in a lightweight image.

### Services

| Service | Port | Description |
|---------|------|-------------|
| `api` | `5050` | BrowersFC REST API |
| `nginx` | `80`, `443` | Reverse proxy |
| `postgres` | `5432` | PostgreSQL database (if enabled) |

### Docker secrets

Sensitive data is managed with Docker secrets:

```yaml
secrets:
  db_url:
  jwt_secret:
  oauth_client_secret:
```

### Container resource limits

```yaml
api:
  limits:
    cpus: '1'
    memory: 1G
  reservations:
    cpus: '0.25'
    memory: 512M
```

## API Documentation

### Swagger / OpenAPI

Interactive API documentation is available at:

- **Development**: `http://localhost:5050/swagger/index.html`
- **Production**: `https://api.yourdomain.com/swagger/index.html`

To regenerate Swagger documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### Base path

All endpoints are served under:

```text
/api
```

### Authentication routes

**Public routes**

```text
GET  /api/users/auth/google
GET  /api/users/auth/google/callback
```

**Protected routes**

```text
GET  /api/users/me
GET  /api/users
GET  /api/users/:username
```

**Admin routes**

```text
POST   /api/admin/users
PUT    /api/admin/users/:id
DELETE /api/admin/users/:id
```

### Main resource endpoints

- `/api/teams` - Team management
- `/api/players` - Player management
- `/api/matches` - Match management
- `/api/seasons` - Season management
- `/api/lineups` - Lineup management
- `/api/articles` - News and articles
- `/api/roles` - Role management

## Database

### Database engine

- **Primary database**: PostgreSQL 15+
- **ORM**: GORM with parameterized queries

### Schema management

Migrations run automatically on startup. Models in the `domain/` package include:

- Embedded `gorm.Model` fields for IDs and timestamps
- Performance-oriented indexes
- Foreign key relationships

### Soft deletes

All entities support soft deletes through the `deleted_at` timestamp. Deleted records are excluded from normal queries.

### Connection string example

```text
postgresql://browersfc:password@localhost:5432/browersfc
```

### Database setup

```sql
CREATE DATABASE browersfc;

CREATE USER browersfc WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE browersfc TO browersfc;
```

## Security Model

### Authentication

**Google OAuth2**

- Users can sign in with their Google account.
- Credentials are not stored in the application database.
- Users can be created automatically on first login.

**JWT**

- Issued after successful authentication
- Signed using the secret from `JWT_SECRET_FILE`
- Sent in the `Authorization: Bearer <token>` header
- Contains user identity, roles, and expiration metadata
- Validated on every protected route

### Authorization

Role-based access control (**RBAC**) is used throughout the API.

| Role | Permissions |
|------|-------------|
| **Admin** | Full access to all endpoints |
| **Coach** | Manage teams and players, create lineups |
| **Player** | View own profile and matches |
| **User** | Read-only access to public resources |

### Security headers

```text
Content-Security-Policy: default-src 'self'; script-src 'self' https://apis.google.com
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

### Input validation

- All request inputs are validated at the handler boundary.
- Validation uses `github.com/go-playground/validator/v10`.
- Parameterized queries help prevent SQL injection.
- CORS is configurable for cross-origin access.

### Middleware stack

1. **CORS middleware** - Controls cross-origin access.
2. **Security headers middleware** - Applies HTTP security headers.
3. **JWT authentication middleware** - Validates and parses tokens.
4. **RBAC middleware** - Enforces role permissions.
5. **Request logging middleware** - Provides structured logs.

## Project Structure

```text
browersfc-api/
├── adapter/                          # Hexagonal adapters
│   ├── http/                         # HTTP response mappers
│   ├── mapper/                       # Authentication mappers
│   └── persistence/                  # Entity-to-model mappers
├── api/                              # HTTP API layer
│   ├── constants/                    # Constants and configuration
│   ├── dto/                          # Data Transfer Objects
│   ├── handler/                      # HTTP request handlers
│   ├── middleware/                   # HTTP middleware
│   └── router/                       # Route definitions
├── cmd/browersfc/                    # Application entry point
├── config/                           # Configuration loading
│   ├── app.go                        # App configuration
│   ├── oauth.go                      # OAuth2 setup
│   └── security.go                   # Security configuration
├── docs/                             # Swagger/OpenAPI documentation
├── domain/                           # Business logic and models
│   ├── *_repository.go               # Repository interfaces (ports)
│   └── *.go                          # Domain models
├── helper/                           # Utility functions
├── internal/                         # Internal-only packages
│   ├── domain/service/               # Domain services
│   └── infrastructure/persistence/   # Database adapters
├── pkg/                              # Shared packages
│   ├── jwt/                          # JWT token handling
│   ├── logger/                       # Structured logging
│   ├── orm/                          # GORM setup
│   ├── security/                     # Security helpers
│   ├── seed/                         # Database seeding
│   └── validation/                   # Custom validators
├── server/                           # Server initialization
├── secrets/                          # Gitignored secret files
├── .env                              # Gitignored environment variables
├── docker-compose.yml                # Docker Compose configuration
├── Dockerfile                        # Multi-stage Docker build
├── go.mod
├── go.sum
└── README.md
```

### Package overview

**`domain/`**

- Contains business models and repository interfaces.
- Has no direct dependency on frameworks.

**`adapter/`**

- Implements ports and maps data between layers.
- Includes persistence and HTTP-facing adapters.

**`internal/`**

- Contains domain services and infrastructure implementations.
- Prevents external packages from depending on internal application details.

**`api/`**

- Contains handlers, routes, DTOs, and middleware.
- Defines the HTTP boundary of the system.

## Development

### Code standards

- Prefer simple, readable code over clever abstractions.
- Write tests for business logic and services.
- Write comments to explain **why**, not **what**.
- Use clear and descriptive names.
- Handle errors explicitly.
- Follow SOLID principles consistently.

### Run tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate an HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the linter
golangci-lint run
```

### Seed development data

```bash
go run ./cmd/browersfc/main.go -seed
```

This command populates the database with sample:

- Users and roles
- Teams and players
- Seasons and matches
- Articles

## CI/CD

### Pre-commit checks

Before committing, run:

```bash
go fmt ./...
go vet ./...
go test ./...
golangci-lint run
```

### Production build

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o browersfc ./cmd/browersfc

# Build Docker image
docker build -t browersfc-api:1.0 .
```

---

Built for **BrowersFC**.
