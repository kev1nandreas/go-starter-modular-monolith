# Go Modular Monolith REST API Template

A production-ready **Go REST API** built with a **Modular Monolith** architecture using **Gin**, featuring **JWT Authentication**, **Rate Limiting**, **Swagger**, **PostgreSQL**, **Redis**, **MongoDB Logging**, and **database seeders**. Fully Dockerized and ready for local development or deployment.

---

## Features

* Modular Monolith architecture (clean separation per domain)
* RESTful API (CRUD-ready)
* JWT Authentication
* API Key security
* Rate Limiting
* Swagger / OpenAPI documentation
* PostgreSQL with GORM
* Redis cache layer
* MongoDB for request logging
* Database migration & seeding (gofakeit)
* End-to-End (E2E) tests
* Docker & Docker Compose
* Hot reload with Air
* Pre-commit hooks & conventional commits

---

## Project Structure

```text
.
├── main.go                         # Application entry point
├── config
│   └── config.go                   # Environment variable helpers
├── src
│   ├── factory
│   │   ├── db.go                   # PostgreSQL connection & migration
│   │   ├── cache.go                # Redis client initialization
│   │   ├── mongo.go                # MongoDB connection
│   │   └── middleware.go           # All middleware (JWT, API Key, CORS, etc.)
│   ├── handler
│   │   ├── router.go              # Route definitions & middleware wiring
│   │   ├── book.go                # Book HTTP handlers
│   │   └── user.go                # User HTTP handlers
│   ├── modules
│   │   ├── book
│   │   │   ├── entity
│   │   │   │   └── book.go        # Book domain entity (GORM model)
│   │   │   ├── dto
│   │   │   │   ├── book_request.go     # Request DTOs
│   │   │   │   └── book_response.go    # Response DTOs
│   │   │   ├── repository
│   │   │   │   ├── book_interface.go   # Repository interface
│   │   │   │   ├── book_postgresql.go  # PostgreSQL implementation
│   │   │   │   └── mocks/             # Generated mocks
│   │   │   ├── usecase
│   │   │   │   ├── book_interface.go   # Usecase interface
│   │   │   │   ├── book.go            # Usecase implementation
│   │   │   │   └── mocks/             # Generated mocks
│   │   │   └── factory
│   │   │       └── book.go            # Module dependency wiring
│   │   └── user
│   │       ├── entity
│   │       │   └── user.go
│   │       ├── dto
│   │       │   ├── user_request.go
│   │       │   └── user_response.go
│   │       ├── repository
│   │       │   ├── user_interface.go
│   │       │   ├── user_postgresql.go
│   │       │   └── mocks/
│   │       ├── usecase
│   │       │   ├── user_interface.go
│   │       │   ├── user.go
│   │       │   └── mocks/
│   │       └── factory
│   │           └── user.go
│   └── shared
│       ├── auth
│       │   └── auth.go             # JWT & password hashing utilities
│       ├── response
│       │   └── response.go         # Standardized JSON response helpers
│       └── pagination
│           └── pagination.go       # Pagination utilities
├── scripts
│   ├── generate_key.go             # Random key generator
│   └── seeders
│       ├── cmd
│       │   └── main.go             # Seeder CLI entry point
│       ├── seeder.go               # Seeder orchestration
│       ├── book_seeder.go          # Book fake data seeder
│       └── user_seeder.go          # User fake data seeder
├── docs                            # Swagger generated docs
├── tests
│   ├── e2e.py                      # E2E tests (Python/pytest)
│   └── requirements.txt
├── Dockerfile
├── Makefile
├── docker-compose.yml
├── .air.toml
├── .env.example
└── .pre-commit-config.yaml
```

---

## Getting Started

### Prerequisites

* Go **1.21+**
* Docker & Docker Compose

### Installation

```bash
git clone https://github.com/kev1nandreas/go-rest-api-template.git
cd go-rest-api-template
cp .env.example .env
```

### Run with Docker

```bash
make up
```

### Run Locally (Development)

```bash
# Start infrastructure containers
docker compose up -d db redis mongo

# Run the application
make run-local
```

Or with hot reload:

```bash
air
```

---

## Environment Variables

Copy `.env.example` to `.env` and configure:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `go-rest-api-template` |
| `APP_PORT` | Server port | `8080` |
| `APP_DEBUG` | Enable debug mode | `true` |
| `APP_MIGRATIONS` | Auto-run migrations on startup | `true` |
| `APP_MONGO_LOGGING` | Enable MongoDB request logging | `true` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_DB` | Database name | — |
| `POSTGRES_USER` | Database user | — |
| `POSTGRES_PASSWORD` | Database password | — |
| `POSTGRES_PORT` | Database port | `5435` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `MONGO_HOST` | MongoDB host | `localhost` |
| `MONGO_PORT` | MongoDB port | — |
| `MONGO_DATABASE` | MongoDB database name | — |
| `MONGO_COLLECTION` | MongoDB collection name | — |
| `JWT_SECRET_KEY` | JWT signing key | — |
| `API_SECRET_KEY` | API key for header auth | — |

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make setup` | Install Swagger tools and generate docs |
| `make build-docker` | Build Docker image (no cache) |
| `make run-local` | Start infra containers + run app locally |
| `make up` | Start all services via Docker Compose |
| `make down` | Stop all containers |
| `make restart` | Restart all containers |
| `make build` | Build the Go application |
| `make test` | Run tests with race detection and coverage |
| `make seed` | Seed database with fake data |
| `make seed-clear` | Clear all seeded data |
| `make clean` | Remove containers and database volumes |

---

## Database Migration

Migrations run automatically on startup when `APP_MIGRATIONS=true`. The migration logic lives in `src/factory/db.go` and uses GORM's `AutoMigrate` with the registered entities.

---

## Database Seeding

Seeders use [gofakeit](https://github.com/brianvoe/gofakeit) for realistic fake data.

```bash
# Seed the database (10 users, 20 books)
make seed

# Clear all seeded data
make seed-clear
```

Manual usage:

```bash
go run scripts/seeders/cmd/main.go up
go run scripts/seeders/cmd/main.go down
```

---

## API Documentation (Swagger)

Once the app is running:

```
http://localhost:8080/swagger/index.html
```

Regenerate docs after modifying annotations:

```bash
make setup
```

---

## API Endpoints

### Auth

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/register` | API Key | Register a new user |
| POST | `/api/v1/login` | API Key | Login and receive JWT |

### Books

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/` | — | Health check |
| GET | `/api/v1/books` | API Key | List books (paginated, cached) |
| GET | `/api/v1/books/:id` | API Key | Get book by ID |
| POST | `/api/v1/books` | API Key + JWT | Create a book |
| PUT | `/api/v1/books/:id` | API Key | Update a book |
| DELETE | `/api/v1/books/:id` | API Key | Delete a book |

### Authentication Headers

```bash
X-API-Key: <your-api-key>
Authorization: Bearer <your-jwt-token>
```

---

## Adding a New Module

To add a new domain module (e.g., `category`):

1. Create the module structure:

```text
src/modules/category/
├── entity/category.go
├── dto/
│   ├── category_request.go
│   └── category_response.go
├── repository/
│   ├── category_interface.go
│   └── category_postgresql.go
├── usecase/
│   ├── category_interface.go
│   └── category.go
└── factory/
    └── category.go
```

2. Register the entity in `src/factory/db.go` for migration.
3. Create a handler in `src/handler/category.go`.
4. Wire routes in `src/handler/router.go`.
5. Initialize via factory in `main.go`.

---

## Testing

### Unit Tests

```bash
make test
```

### E2E Tests (Python)

```bash
cd tests
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt

export BASE_URL=http://localhost:8080/api/v1
export API_KEY=your-api-key

pytest e2e.py
```

---

## Clean Up Script

Remove boilerplate entities (e.g., books) using the cleanup script:

```bash
bash cleanup.sh
```

> Review the script before running — it removes files by name pattern.
