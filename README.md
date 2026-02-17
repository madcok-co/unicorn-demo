# 🦄 Unicorn Demo API

A full-featured enterprise-grade REST API demonstrating all capabilities of the [Unicorn Framework](https://github.com/madcok-co/unicorn).

## 🌟 Features

This demo showcases all 6 enterprise features of the Unicorn framework:

- ✅ **OAuth2/OIDC Authentication** - Google, GitHub, Microsoft providers
- ✅ **RBAC Authorization** - 5 roles with hierarchical permissions
- ✅ **Multi-Tenancy** - Subdomain-based tenant isolation
- ✅ **Configuration Management** - Environment-based config with hot reload
- ✅ **Pagination** - Offset-based and cursor-based pagination
- ✅ **API Versioning** - v1 (offset) and v2 (cursor + batch operations)

### Additional Features

- 🗄️ **PostgreSQL/SQLite** - Production-ready database with GORM
- 🔐 **Security** - JWT authentication, CORS, request validation
- 📊 **Analytics** - Project and user statistics endpoints
- 📝 **Audit Logging** - Complete audit trail of all actions
- 🐳 **Docker Ready** - Full Docker Compose setup
- 🏥 **Health Checks** - Kubernetes-ready health endpoints

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Architecture](#-architecture)
- [API Documentation](#-api-documentation)
- [Multi-Tenancy](#-multi-tenancy)
- [Authorization](#-authorization)
- [Development](#-development)
- [Deployment](#-deployment)
- [Testing](#-testing)

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- (Optional) Docker & Docker Compose
- (Optional) PostgreSQL 15+

### Local Development

```bash
# Clone the repository
git clone https://github.com/madcok-co/unicorn-demo
cd unicorn-demo

# Install dependencies
make deps

# Copy environment file
cp .env.example .env

# Run the application (SQLite - no setup needed)
make run

# Or with Docker (PostgreSQL + Redis)
make docker-up
```

The API will be available at:
- http://localhost:8080/health
- http://localhost:8080/api/v1/projects
- http://localhost:8080/api/v2/projects

### Setup Multi-Tenant Hosts

For subdomain-based multi-tenancy:

```bash
# Add to /etc/hosts (Linux/Mac)
sudo sh -c 'echo "127.0.0.1  acme.localhost" >> /etc/hosts'
sudo sh -c 'echo "127.0.0.1  techcorp.localhost" >> /etc/hosts'

# Or use the Makefile
make setup-hosts
```

Now you can access tenant-specific URLs:
- http://acme.localhost:8080/api/v1/projects
- http://techcorp.localhost:8080/api/v1/projects

## 🏗️ Architecture

### Project Structure

```
unicorn-demo/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/
│   │   ├── models.go              # Database models
│   │   └── dtos.go                # Request/Response DTOs
│   ├── handlers/
│   │   ├── auth.go                # Authentication handlers
│   │   ├── projects.go            # Project CRUD
│   │   ├── users.go               # User management
│   │   ├── tasks.go               # Task management
│   │   ├── health.go              # Health checks
│   │   └── helpers.go             # Shared utilities
│   ├── middleware/
│   │   ├── auth.go                # JWT middleware
│   │   └── tenant.go              # Tenant resolution
│   ├── database/
│   │   └── database.go            # Database connection & migrations
│   └── services/                  # Business logic services
├── config/                        # Configuration files
├── Dockerfile                     # Docker image
├── docker-compose.yml             # Multi-container setup
├── Makefile                       # Build automation
└── README.md                      # This file
```

### Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | [Unicorn](https://github.com/madcok-co/unicorn) |
| Language | Go 1.21+ |
| Database | PostgreSQL 15 / SQLite |
| Cache | Redis 7 |
| ORM | GORM |
| Auth | OAuth2 + JWT |
| Container | Docker + Docker Compose |

## 📚 API Documentation

### Authentication Endpoints

#### POST /api/v1/auth/login
Email/password login (demo only)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "dev@acme.com",
    "password": "password"
  }'
```

#### GET /api/v1/auth/callback
OAuth2 callback endpoint

```bash
# After OAuth provider redirects
curl "http://localhost:8080/api/v1/auth/callback?code=AUTH_CODE&state=STATE"
```

#### POST /api/v1/auth/refresh
Refresh access token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

#### GET /api/v1/auth/me
Get current user

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Project Endpoints

#### GET /api/v1/projects
List projects with pagination

```bash
curl "http://localhost:8080/api/v1/projects?page=1&limit=20&sort=created_at&order=desc" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### POST /api/v1/projects
Create new project

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Project",
    "description": "Project description",
    "tags": ["demo", "test"]
  }'
```

#### GET /api/v1/projects/:id
Get project by ID

```bash
curl http://localhost:8080/api/v1/projects/proj-1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### PUT /api/v1/projects/:id
Update project

```bash
curl -X PUT http://localhost:8080/api/v1/projects/proj-1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Project Name",
    "status": "inactive"
  }'
```

#### DELETE /api/v1/projects/:id
Delete project

```bash
curl -X DELETE http://localhost:8080/api/v1/projects/proj-1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### GET /api/v1/projects/stats
Get project statistics

```bash
curl http://localhost:8080/api/v1/projects/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### User Endpoints

#### GET /api/v1/users
List users (admin only)

```bash
curl "http://localhost:8080/api/v1/users?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### POST /api/v1/users
Create user (admin only)

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@acme.com",
    "name": "New User",
    "roles": ["developer"]
  }'
```

### Task Endpoints

#### GET /api/v1/tasks
List tasks

```bash
curl "http://localhost:8080/api/v1/tasks?project_id=proj-1&status=todo" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### POST /api/v1/tasks
Create task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "proj-1",
    "title": "New Task",
    "description": "Task description",
    "priority": "high",
    "assignee_id": "user-dev-acme"
  }'
```

### V2 Endpoints (Enhanced Features)

#### POST /api/v2/projects/batch
Batch update projects

```bash
curl -X POST http://localhost:8080/api/v2/projects/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_ids": ["proj-1", "proj-2"],
    "updates": {
      "status": "archived"
    }
  }'
```

### Health Check

#### GET /health
Application health status

```bash
curl http://localhost:8080/health
```

## 🏢 Multi-Tenancy

This demo uses **subdomain-based** tenant isolation:

### Demo Tenants

| Tenant | Subdomain | Plan | Max Users |
|--------|-----------|------|-----------|
| Acme Corporation | acme.localhost | Enterprise | 100 |
| Tech Corp | techcorp.localhost | Professional | 50 |

### Testing Multi-Tenancy

```bash
# Acme tenant
curl -H "Host: acme.localhost" http://localhost:8080/api/v1/projects

# Tech Corp tenant
curl -H "Host: techcorp.localhost" http://localhost:8080/api/v1/projects

# Data is isolated - Acme cannot see Tech Corp's projects
```

### Alternative Strategies

The demo supports multiple tenant identification strategies:

**Header-based:**
```bash
export APP_MULTITENANCY_STRATEGY=header
curl -H "X-Tenant-ID: acme" http://localhost:8080/api/v1/projects
```

**Path-based:**
```bash
export APP_MULTITENANCY_STRATEGY=path
curl http://localhost:8080/acme/api/v1/projects
```

## 🔒 Authorization

### Roles & Permissions

The demo includes 5 pre-configured roles with hierarchical inheritance:

```
super_admin (*)
    ↓
tenant_admin (projects:*, users:*, tasks:*, comments:*)
    ↓
project_manager (projects:read/create/update, tasks:*, users:read, comments:*)
    ↓
developer (projects:read, tasks:read/update, comments:*)
    ↓
viewer (*:read)
```

### Permission Format

Permissions follow the format: `resource:action`

Examples:
- `projects:read` - Read projects
- `projects:create` - Create projects
- `projects:*` - All project operations
- `*:read` - Read all resources
- `*` - All permissions

### Demo Users

| Email | Tenant | Role | Can Do |
|-------|--------|------|--------|
| admin@acme.com | acme | tenant_admin | Everything in tenant |
| pm@acme.com | acme | project_manager | Manage projects & tasks |
| dev@acme.com | acme | developer | Update tasks, read projects |

### Testing Authorization

```bash
# Login as developer
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@acme.com","password":"password"}' | jq -r '.token')

# Can read projects ✅
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/projects

# Can update tasks ✅
curl -X PUT -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tasks/task-1 \
  -d '{"status":"done"}'

# Cannot delete projects ❌
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/projects/proj-1
# Returns: 403 Forbidden
```

## 💻 Development

### Running Locally

```bash
# With SQLite (default - no setup)
make run

# With PostgreSQL
export APP_DB_DRIVER=postgres
export APP_DB_HOST=localhost
export APP_DB_USER=unicorn
export APP_DB_PASSWORD=unicorn_secret
make run
```

### Hot Reload

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Database Migrations

```bash
# Run migrations
make db-migrate

# Seed demo data
make db-seed
```

## 🐳 Deployment

### Docker Compose

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Environment Variables

Key environment variables:

```bash
# Application
APP_HTTP_HOST=0.0.0.0
APP_HTTP_PORT=8080

# Database
APP_DB_DRIVER=postgres
APP_DB_HOST=postgres
APP_DB_PORT=5432
APP_DB_USER=unicorn
APP_DB_PASSWORD=unicorn_secret
APP_DB_DATABASE=unicorn_demo

# OAuth2
APP_OAUTH_PROVIDER=google
APP_OAUTH_CLIENT_ID=your-client-id
APP_OAUTH_CLIENT_SECRET=your-client-secret
APP_OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/auth/callback

# Multi-tenancy
APP_MULTITENANCY_STRATEGY=subdomain
APP_MULTITENANCY_DOMAIN=localhost
```

### Kubernetes

Coming soon - Kubernetes manifests for production deployment.

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package
go test -v ./internal/handlers/
```

## 📊 Performance

Expected performance metrics:

- **Throughput**: 10K+ requests/second (single instance)
- **Latency**: < 10ms (p95) for simple CRUD operations
- **Memory**: ~50 MB idle, ~200 MB under load
- **Database**: Optimized with connection pooling (10 idle, 100 max)

## 🤝 Contributing

Contributions are welcome! This is a demo project to showcase the Unicorn framework.

## 📄 License

MIT License

## 🔗 Links

- [Unicorn Framework](https://github.com/madcok-co/unicorn)
- [Documentation](https://github.com/madcok-co/unicorn/docs)
- [Report Issues](https://github.com/madcok-co/unicorn-demo/issues)

## 🙏 Acknowledgments

Built with the [Unicorn Framework](https://github.com/madcok-co/unicorn) - A batteries-included Go framework for building enterprise-grade APIs.

---

**Built with ❤️ using Unicorn Framework 🦄**
