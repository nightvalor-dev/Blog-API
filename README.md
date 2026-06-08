# Modular Monolithic Blog System API

A robust, high-performance CRUD RESTful API for a blogging platform built with **Go (Golang)**. The project follows a **Modular Monolithic Architecture** — clear separation of concerns per domain, single deployable unit, horizontally scalable.

Features Role-Based Access Control (RBAC), JWT authentication, structured logging, database migrations, connection pooling, Redis cache-aside strategy, advanced filtering/sorting/pagination, and multi-container Docker deployment.

---

## Table of Contents

- [Features](#-features)
- [Tech Stack](#️-tech-stack--dependencies)
- [Project Structure](#-project-structure)
- [Architecture](#-domain-module-architecture)
- [Configuration](#️-local-environment-configuration)
- [Getting Started](#-getting-started)
- [API Reference](#️-api-endpoint-reference)
- [Advanced Paradigms](#-advanced-paradigms)
- [Docker Deployment](#-docker-deployment)
- [License](#-license)

---

## ✨ Features

- **Modular Monolithic Design** — Codebase organized tightly by domain (`auth`, `user`, `blog`, `category`, `comment`, `tag`), minimizing cross-module coupling while remaining single-deployable.
- **Advanced Query Support** — Out-of-the-box pagination, dynamic field sorting, and relational filtering on all content endpoints.
- **Dual-Layer Storage Architecture**
  - **Primary DB:** PostgreSQL via high-performance connection pooling (`pgx/v5`)
  - **Cache Layer:** Redis (`go-redis/v9`) with a cache-aside pattern for heavy read traffic on blog entries
- **Security & Authorization**
  - Secure password hashing via `bcrypt`
  - Stateless authentication with JSON Web Tokens (JWT)
  - Granular **Role-Based Access Control (RBAC)** middleware securing admin and moderation endpoints
- **Production-Grade Middleware Stack** — Global Panic Recovery, CORS, Structured Request-ID Logging, Auth context injection
- **Asynchronous Event Communication** — Built-in in-memory event bus for loosely coupled inter-module notifications
- **Containerized Infrastructure** — Multi-service Docker Compose environment with raw SQL sequential schema migrations

---

## 🛠️ Tech Stack & Dependencies

### Core Runtime

| Component | Detail |
|-----------|--------|
| Language | Go `v1.25.9` |

### Frameworks & Drivers

| Dependency | Purpose |
|------------|---------|
| [go-chi/chi/v5](https://github.com/go-chi/chi) | Lightweight, idiomatic, context-aware HTTP router |
| [jackc/pgx/v5](https://github.com/jackc/pgx) | Low-allocation, high-performance PostgreSQL driver + connection pool |
| [redis/go-redis/v9](https://github.com/redis/go-redis) | Type-safe Redis client |
| [go-playground/validator/v10](https://github.com/go-playground/validator) | Struct validation using tags |
| [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | Secure JWT generation and parsing |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Bcrypt password hashing |
| [joho/godotenv](https://github.com/joho/godotenv) | `.env` file runtime reader |

---

## 📂 Project Structure

```
├── api/
│   └── routes.go                # Central HTTP router (aggregates all module routes)
├── cmd/
│   └── server/
│       └── main.go              # Entrypoint — initializes config, DBs, modules, HTTP server
├── config/
│   └── config.go                # Strongly-typed environment config loader
├── docker-compose.yml           # Multi-container orchestration (PostgreSQL, Redis)
├── go.mod
├── go.sum
├── internal/                    # Private application code
│   ├── modules/
│   │   ├── auth/                # Sign-in, registration, credentials
│   │   ├── blog/                # Post CRUD + Redis caching layers
│   │   ├── category/            # Categorization taxonomy management
│   │   ├── comment/             # User discussions and feedback
│   │   ├── tag/                 # Many-to-many metadata tagging
│   │   └── user/                # Profile management, identities, permissions
│   └── shared/
│       ├── db/                  # PostgreSQL connection pool engine
│       ├── events/              # Internal pub/sub event bus
│       ├── middleware/          # HTTP interceptors (Auth, CORS, Logger, Recovery)
│       ├── redis/               # Redis client setup and key builders
│       └── roles/               # RBAC definitions and privilege matrices
├── logs/
│   └── app.log                  # Rolling structured application log
├── migrations/                  # Sequential raw SQL schema migration scripts (Up/Down)
├── pkg/
│   ├── jwt/                     # JWT generation and token validation
│   └── utils/                   # HTTP utilities (Pagination, Parsers, Validators, JSON Writers)
└── README.md
```

---

## 🧱 Domain Module Architecture

Every module under `internal/modules/` adheres strictly to a layered architecture with explicit boundaries:

| Layer | File | Responsibility |
|-------|------|----------------|
| **Domain Models** | `entity.go` | Primary data structures representing core DB tables |
| **DTOs** | `*DTO.go` | Input validation rules and masked output schemas |
| **Repository Interface** | `repository.go` | Contracts for persistent storage interactions |
| **Repository Impl** | `repositoryImpl.go` | Raw SQL queries via `pgx` |
| **Service** | `service.go` | Business logic, transactions, cache invalidation |
| **Handler** | `handler.go` | HTTP request parsing, validation, service invocation, status mapping |
| **Router** | `routes.go` | Module endpoint bindings and local RBAC guards |

---

## ⚙️ Local Environment Configuration

Create a `.env` file in the project root:

```env
# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Security
JWT_SECRET=your_super_secret_high_entropy_ecdsa_or_hmac_string
JWT_EXPIRATION_HOURS=24

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=blog_admin
DB_PASSWORD=secure_postgres_pass
DB_NAME=blog_monolith
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
```

---

## 🚀 Getting Started

### Prerequisites

- Go SDK `v1.25.9` or newer
- Docker Engine & Docker Compose CLI
- `golang-migrate` CLI _(optional, for manual schema management)_

### 1. Provision Infrastructure

Spin up PostgreSQL and Redis containers in detached mode:

```bash
docker compose up -d
```

### 2. Run Database Migrations

```bash
migrate -path ./migrations \
  -database "postgres://blog_admin:secure_postgres_pass@localhost:5432/blog_monolith?sslmode=disable" \
  up
```

> Alternatively, run `migrations/000001_init_schema.up.sql` manually via your preferred database client.

### 3. Start the Server

```bash
# Fetch dependencies
go mod download

# Run the application
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`.

---

## 🛰️ API Endpoint Reference

### 🔐 Authentication

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/auth/register` | Public | Register a new user account |
| `POST` | `/api/v1/auth/login` | Public | Authenticate and receive a JWT token |

### 👤 Users

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `GET` | `/api/v1/users/profile` | Authenticated | Retrieve current session profile |
| `PUT` | `/api/v1/users/profile` | Authenticated | Update active account attributes |

### 📝 Blogs _(Redis-cached public reads)_

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/blogs` | Author / Admin | Create a new blog post |
| `GET` | `/api/v1/blogs` | Public | Paginated list with filter/sort support |
| `GET` | `/api/v1/blogs/{id}` | Public | Fetch full blog entry (cache-aside) |
| `PUT` | `/api/v1/blogs/{id}` | Owner / Admin | Update post (invalidates Redis cache) |
| `DELETE` | `/api/v1/blogs/{id}` | Owner / Admin | Remove post from DB and cache |

### 🏷️ Categories & Tags

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/categories` | Admin | Create a taxonomy category |
| `GET` | `/api/v1/categories` | Public | List all categories |
| `POST` | `/api/v1/tags` | Author / Admin | Create a new tag |
| `GET` | `/api/v1/tags` | Public | List all global tags |

### 💬 Comments

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| `POST` | `/api/v1/blogs/{id}/comments` | Authenticated | Append a comment to a post |
| `GET` | `/api/v1/blogs/{id}/comments` | Public | Retrieve comment thread for a post |
| `DELETE` | `/api/v1/comments/{id}` | Owner / Moderator | Delete a comment |

---

## 📊 Advanced Paradigms

### ⚡ Query Parameters

All list endpoints (`/blogs`, `/comments`) support standard URL query parameters:

| Parameter | Example | Description |
|-----------|---------|-------------|
| Pagination | `?page=1&limit=10` | Predictable resource chunking |
| Sorting | `?sort_by=created_at&order=desc` | Dynamic field sorting |
| Filtering | `?search=golang&category_id=3` | Relational content filtering |

### 💾 Cache-Aside Pattern

When a `GET /api/v1/blogs/{id}` request arrives:

1. System builds a lookup key via `internal/shared/redis/key.go` (e.g., `blog:post:{id}`)
2. **Cache Hit** → data returned instantly from Redis, bypassing the database
3. **Cache Miss** → PostgreSQL queried, result written to Redis with an explicit TTL, then served
4. **Mutations** (`PUT`, `DELETE`) → database updated first, then cache keys immediately evicted to prevent stale reads

### 🛡️ Middleware Chain

```
Request ──► [Recovery] ──► [CORS] ──► [Logger] ──► [Auth Context] ──► [RBAC Evaluator] ──► Handler
```

Defined in `internal/shared/middleware/`.

---

## 🐳 Docker Deployment

The environment uses isolated Docker networking:

- **`postgres:16-alpine`** — Persistent volume at `.docker/db_data`
- **`redis:7-alpine`** — In-memory store on default ports

```bash
# Start infrastructure
docker compose up -d

# Tail live logs
docker compose logs -f

# Teardown and prune volumes
docker compose down -v
```

---

## 📜 License

This project is open-source software licensed under the **MIT License**.
