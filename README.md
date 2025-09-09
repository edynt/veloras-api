# Veloras E-commerce API

A clean, modern Go-based e-commerce API built with clean architecture principles.

## Features

- **Authentication & Authorization**: User registration, login, JWT tokens
- **Role-Based Access Control (RBAC)**: Roles and permissions management
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Database**: PostgreSQL with auto-increment integer IDs
- **API Documentation**: Swagger/OpenAPI documentation
- **Email Verification**: OTP-based email verification

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **ORM**: SQLC for type-safe SQL queries
- **Web Framework**: Gin
- **Authentication**: JWT
- **Documentation**: Swagger
- **Migration**: Goose

## Project Structure

```
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── swag/              # Swagger documentation
├── internal/              # Private application code
│   ├── auth/              # Authentication module
│   │   ├── application/   # Application layer (services, DTOs)
│   │   ├── controller/    # Presentation layer (handlers, routers)
│   │   ├── domain/        # Domain layer (entities, repositories)
│   │   └── infrastructure/ # Infrastructure layer (database)
│   ├── initialize/        # Application initialization
│   ├── middleware/        # HTTP middleware
│   └── shared/            # Shared utilities and database
├── pkg/                   # Public packages
│   ├── config/            # Configuration management
│   ├── response/          # HTTP response utilities
│   └── utils/             # Utility functions
└── templates-email/       # Email templates
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd veloras-version-ai
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp pkg/environment/config.yaml.example pkg/environment/config.yaml
# Edit config.yaml with your database and other settings
```

4. Run database migrations:
```bash
make migrate-up
```

5. Generate SQLC code:
```bash
make sqlc-generate
```

6. Run the application:
```bash
go run cmd/server/main.go
```

### API Documentation

Once the server is running, visit:
- Swagger UI: http://localhost:8080/swagger/index.html

## Database Schema

The application uses PostgreSQL with auto-increment integer IDs for all entities:

- **users**: User accounts with authentication
- **roles**: User roles for RBAC
- **permissions**: Granular permissions
- **user_roles**: Many-to-many relationship between users and roles
- **role_permissions**: Many-to-many relationship between roles and permissions
- **sessions**: JWT refresh token storage
- **email_verifications**: OTP verification codes
- **password_resets**: Password reset tokens

## Development

### Available Make Commands

```bash
make help              # Show available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make sqlc-generate     # Generate SQLC code
make swag              # Generate Swagger documentation
```

### Code Generation

The project uses SQLC for type-safe database queries:

```bash
# After modifying SQL files in internal/shared/queries/
make sqlc-generate
```

## License

This project is licensed under the MIT License.
