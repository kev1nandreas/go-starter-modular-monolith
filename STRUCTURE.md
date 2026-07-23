# Project Architecture & Structure

This document provides an in-depth explanation of the modular monolith architecture used in this project.

---

## Table of Contents

- [Architecture Pattern](#architecture-pattern)
- [Directory Structure](#directory-structure)
- [Layer Responsibilities](#layer-responsibilities)
- [Request Flow](#request-flow)
- [Module Anatomy](#module-anatomy)
- [Shared Packages](#shared-packages)
- [Factory Pattern](#factory-pattern)
- [Dependency Injection](#dependency-injection)
- [Adding a New Module](#adding-a-new-module)

---

## Architecture Pattern

This project implements a **Modular Monolith** — a single deployable application where the codebase is organized into independent, loosely-coupled modules. Each module encapsulates a specific business domain.

### Why Modular Monolith?

| Benefit | Description |
|---------|-------------|
| **Clear boundaries** | Each module owns its domain logic, entities, and data access |
| **Easy to test** | Interfaces at every boundary enable unit testing with mocks |
| **Simple deployment** | Single binary — no service mesh or network overhead |
| **Migration-ready** | Modules can be extracted into microservices later if needed |
| **Team scalability** | Multiple developers can work on different modules independently |

### Clean Architecture Layers

```
┌─────────────────────────────────────────────┐
│                  main.go                     │  ← Entry point & DI wiring
├─────────────────────────────────────────────┤
│                 Handler                      │  ← HTTP layer (Gin)
├─────────────────────────────────────────────┤
│                 Usecase                      │  ← Business logic
├─────────────────────────────────────────────┤
│                Repository                    │  ← Data access abstraction
├─────────────────────────────────────────────┤
│            Database / Cache                  │  ← Infrastructure
└─────────────────────────────────────────────┘
```

Dependencies point **inward** — outer layers depend on inner layers, never the reverse. The usecase layer defines interfaces that the repository layer implements.

---

## Directory Structure

```text
.
├── main.go                             # Application entry point
├── config/                             # Application configuration
├── src/
│   ├── factory/                        # Infrastructure initialization
│   ├── handler/                        # HTTP handlers & routing
│   ├── modules/                        # Domain modules
│   │   ├── book/                       # Book module
│   │   └── user/                       # User module
│   └── shared/                         # Cross-cutting utilities
├── scripts/                            # CLI tools & seeders
├── docs/                               # Swagger generated documentation
└── tests/                              # E2E tests
```

### Top-Level Files

| File | Purpose |
|------|---------|
| `main.go` | Application bootstrap — initializes infrastructure, creates modules, starts server |
| `Makefile` | Build, run, test, and seed commands |
| `Dockerfile` | Multi-stage production build |
| `docker-compose.yml` | Service orchestration (app, PostgreSQL, Redis, MongoDB) |
| `.air.toml` | Hot-reload configuration for development |
| `.env.example` | Environment variable template |

---

## Layer Responsibilities

### `config/`

```text
config/
└── config.go
```

Provides typed environment variable accessors (`GetEnvString`, `GetEnvInt`, `GetEnvBool`). This replaces hardcoded values and centralizes all configuration reading. No business logic lives here.

---

### `src/factory/`

```text
src/factory/
├── db.go           # PostgreSQL connection, retry logic, auto-migration
├── cache.go        # Redis client initialization
├── mongo.go        # MongoDB connection for logging
└── middleware.go   # All HTTP middleware definitions
```

The factory package is responsible for **infrastructure initialization**. It creates and configures the external dependencies that modules need:

- **db.go** — Connects to PostgreSQL with retry logic. Runs GORM `AutoMigrate` when `APP_MIGRATIONS=true`.
- **cache.go** — Creates a Redis client from environment config.
- **mongo.go** — Connects to MongoDB for request logging.
- **middleware.go** — Defines all Gin middleware:
  - `APIKeyAuth()` — Validates `X-API-Key` header
  - `JWTAuth()` — Validates Bearer JWT tokens
  - `Cors()` — CORS configuration
  - `Security()` — Security headers (HSTS, CSP, X-Frame-Options)
  - `Xss()` — XSS protection
  - `RateLimiter()` — Token bucket rate limiting
  - `Logger()` — Request logging to MongoDB

---

### `src/handler/`

```text
src/handler/
├── router.go       # Route definitions, middleware chain, Swagger setup
├── book.go         # Book HTTP handler (CRUD + cache)
└── user.go         # User HTTP handler (login, register)
```

Handlers are the **HTTP transport layer**. They:

1. Parse and validate incoming requests
2. Delegate business logic to usecases
3. Format and return responses

Handlers know about HTTP concepts (status codes, headers, JSON binding) but do NOT contain business logic. They receive usecases via constructor injection.

**router.go** wires everything together:
- Applies global middleware (CORS, rate limiting, security)
- Defines route groups with appropriate auth middleware
- Mounts Swagger UI

---

### `src/modules/`

Each module is a self-contained domain boundary:

```text
src/modules/<module>/
├── entity/             # Domain models (GORM structs)
├── dto/                # Data Transfer Objects (request/response shapes)
├── repository/         # Data access layer
│   ├── <name>_interface.go     # Interface definition
│   ├── <name>_postgresql.go    # PostgreSQL implementation
│   └── mocks/                  # Generated test mocks
├── usecase/            # Business logic layer
│   ├── <name>_interface.go     # Interface definition
│   ├── <name>.go               # Implementation
│   └── mocks/                  # Generated test mocks
└── factory/            # Module-level dependency wiring
    └── <name>.go
```

Modules do NOT import from each other. If two modules need to communicate, they should do so through shared interfaces or events (not direct imports).

---

### `src/shared/`

```text
src/shared/
├── auth/
│   └── auth.go         # JWT generation, password hashing, key generation
├── response/
│   └── response.go     # Standardized JSON response structure
└── pagination/
    └── pagination.go   # Pagination parameter parsing & GORM integration
```

Shared packages provide **cross-cutting utilities** that multiple modules need. Rules:

- Shared packages must NOT import from any module
- They contain only generic, reusable logic
- They define no business rules specific to a domain

---

## Request Flow

Here's what happens when a request hits `POST /api/v1/books`:

```
Client
  │
  ▼
┌──────────────────────────────────────────────────┐
│  Gin Router (src/handler/router.go)              │
│  ├── factory.APIKeyAuth()    → validates API key │
│  ├── factory.JWTAuth()       → validates JWT     │
│  └── bookHandler.CreateBook  → handler method    │
└──────────────────────────────────────────────────┘
  │
  ▼
┌──────────────────────────────────────────────────┐
│  BookHandler (src/handler/book.go)               │
│  ├── Binds JSON to dto.CreateBookRequest         │
│  ├── Calls usecase.Create(request)               │
│  ├── Invalidates Redis cache                     │
│  └── Returns JSON response                       │
└──────────────────────────────────────────────────┘
  │
  ▼
┌──────────────────────────────────────────────────┐
│  BookUsecase (src/modules/book/usecase/book.go)  │
│  ├── Constructs entity.Book from DTO            │
│  ├── Calls repository.Create(book)               │
│  └── Returns created book                        │
└──────────────────────────────────────────────────┘
  │
  ▼
┌──────────────────────────────────────────────────┐
│  BookRepository (src/modules/book/repository/)   │
│  └── Executes db.Create(book) via GORM          │
└──────────────────────────────────────────────────┘
  │
  ▼
┌──────────────────────────────────────────────────┐
│  PostgreSQL                                      │
└──────────────────────────────────────────────────┘
```

---

## Module Anatomy

### Entity

The domain model. Maps directly to a database table via GORM tags.

```go
// src/modules/book/entity/book.go
type Book struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Title     string
    Author    string
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

### DTO (Data Transfer Object)

Shapes data for API input/output. Separates internal representation from external contracts.

```go
// src/modules/book/dto/book_request.go
type CreateBookRequest struct {
    Title  string `json:"title" binding:"required"`
    Author string `json:"author" binding:"required"`
}
```

### Repository

Defines **what** data operations are available (interface) and **how** they're implemented (PostgreSQL, MySQL, etc.).

```go
// Interface — no implementation details
type BookRepository interface {
    FindAll(offset, limit int) ([]entity.Book, error)
    FindByID(id string) (*entity.Book, error)
    Create(book *entity.Book) error
    Update(book *entity.Book) error
    Delete(book *entity.Book) error
}
```

```go
// Implementation — specific to PostgreSQL/GORM
type bookPostgresRepository struct {
    db *gorm.DB
}

func (r *bookPostgresRepository) FindByID(id string) (*entity.Book, error) {
    var book entity.Book
    if err := r.db.Where("id = ?", id).First(&book).Error; err != nil {
        return nil, err
    }
    return &book, nil
}
```

### Usecase

Contains **business logic**. Orchestrates repository calls and enforces domain rules.

```go
type bookUsecase struct {
    repo repository.BookRepository
}

func (u *bookUsecase) Create(req dto.CreateBookRequest) (*entity.Book, error) {
    book := &entity.Book{
        Title:  req.Title,
        Author: req.Author,
    }
    if err := u.repo.Create(book); err != nil {
        return nil, err
    }
    return book, nil
}
```

### Module Factory

Wires the module's dependencies together. Called from `main.go`.

```go
// src/modules/book/factory/book.go
func NewBookUsecase(db *gorm.DB) usecase.BookUsecase {
    repo := repository.NewBookPostgresRepository(db)
    return usecase.NewBookUsecase(repo)
}
```

---

## Shared Packages

### `shared/auth`

| Function | Description |
|----------|-------------|
| `GenerateToken(username)` | Creates a signed JWT with 5-minute expiry |
| `HashPassword(password)` | Bcrypt hash (cost 14) |
| `GenerateRandomKey()` | 256-bit random key (base64 encoded) |

### `shared/response`

Provides a consistent JSON envelope for all API responses:

```json
{
  "success": true,
  "message": "Books retrieved successfully",
  "data": [...],
  "meta": { "page": 1, "total_items": 50 }
}
```

### `shared/pagination`

Parses `page`, `page_size`, `sort`, `order` query params and applies them to GORM queries. Returns metadata (total items, total pages, has_next, has_prev).

---

## Factory Pattern

The factory pattern is used at two levels:

### 1. Infrastructure Factory (`src/factory/`)

Creates infrastructure clients. Called once at startup in `main.go`:

```go
db := factory.NewDatabase()          // PostgreSQL
redisClient := factory.NewRedisClient()  // Redis
mongoCol := factory.NewMongoDB()     // MongoDB
```

### 2. Module Factory (`src/modules/<name>/factory/`)

Wires module internals. Accepts infrastructure and returns a ready-to-use usecase:

```go
bookUsecase := bookFactory.NewBookUsecase(db)
userUsecase := userFactory.NewUserUsecase(db)
```

This keeps `main.go` clean — it doesn't need to know about repositories or internal module structure.

---

## Dependency Injection

Dependencies flow **top-down** from `main.go`:

```
main.go
  ├── Creates infrastructure (db, redis, mongo)
  ├── Creates usecases via module factories
  ├── Creates handlers (injecting usecases + redis)
  └── Creates router (injecting handlers + logger)
```

No global variables. No service locator. Each component declares its dependencies explicitly in its constructor:

```go
func NewBookHandler(uc usecase.BookUsecase, redis *redis.Client, ctx context.Context) *BookHandler
func NewUserHandler(uc usecase.UserUsecase) *UserHandler
func NewRouter(logger, mongo, bookHandler, userHandler) *gin.Engine
```

---

## Adding a New Module

Example: adding a `category` module.

### 1. Create module structure

```bash
mkdir -p src/modules/category/{entity,dto,repository/mocks,usecase/mocks,factory}
```

### 2. Define the entity

```go
// src/modules/category/entity/category.go
package entity

type Category struct {
    ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name string    `gorm:"unique"`
}
```

### 3. Define DTOs

```go
// src/modules/category/dto/category_request.go
package dto

type CreateCategoryRequest struct {
    Name string `json:"name" binding:"required"`
}
```

### 4. Define repository interface & implementation

```go
// src/modules/category/repository/category_interface.go
type CategoryRepository interface {
    FindAll() ([]entity.Category, error)
    Create(category *entity.Category) error
}
```

```go
// src/modules/category/repository/category_postgresql.go
func NewCategoryPostgresRepository(db *gorm.DB) CategoryRepository { ... }
```

### 5. Define usecase interface & implementation

```go
// src/modules/category/usecase/category_interface.go
type CategoryUsecase interface {
    FindAll() ([]entity.Category, error)
    Create(req dto.CreateCategoryRequest) (*entity.Category, error)
}
```

### 6. Create module factory

```go
// src/modules/category/factory/category.go
func NewCategoryUsecase(db *gorm.DB) usecase.CategoryUsecase {
    repo := repository.NewCategoryPostgresRepository(db)
    return usecase.NewCategoryUsecase(repo)
}
```

### 7. Register entity for migration

In `src/factory/db.go`, add to the `Migrate` function:

```go
import categoryEntity "github.com/kev1nandreas/go-rest-api-template/src/modules/category/entity"

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&bookEntity.Book{}, &userEntity.User{}, &categoryEntity.Category{})
}
```

### 8. Create handler

```go
// src/handler/category.go
type CategoryHandler struct { usecase usecase.CategoryUsecase }
func NewCategoryHandler(uc usecase.CategoryUsecase) *CategoryHandler { ... }
```

### 9. Wire in router

In `src/handler/router.go`, add routes:

```go
v1.GET("/categories", factory.APIKeyAuth(), categoryHandler.FindAll)
v1.POST("/categories", factory.APIKeyAuth(), factory.JWTAuth(), categoryHandler.Create)
```

### 10. Wire in main.go

```go
categoryUsecase := categoryFactory.NewCategoryUsecase(db)
categoryHandler := handler.NewCategoryHandler(categoryUsecase)
// Pass to NewRouter
```

---

## Summary

| Concept | Location | Responsibility |
|---------|----------|----------------|
| Configuration | `config/` | Read environment variables |
| Infrastructure | `src/factory/` | Create DB, cache, middleware |
| Transport | `src/handler/` | HTTP routing & request handling |
| Business Logic | `src/modules/*/usecase/` | Domain rules & orchestration |
| Data Access | `src/modules/*/repository/` | Database queries |
| Domain Models | `src/modules/*/entity/` | Struct definitions |
| API Contracts | `src/modules/*/dto/` | Request/response shapes |
| Module Wiring | `src/modules/*/factory/` | Connect repo → usecase |
| Utilities | `src/shared/` | Auth, response, pagination |
