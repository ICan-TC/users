# NeoICan Users Service

A reusable template for building backend APIs in Go using the [Huma](https://huma.rocks/) framework. This service provides user authentication, management, and token handling with a clean, scalable architecture.

## Table of Contents

- [Purpose](#purpose)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
- [Build, Run, and Test](#build-run-and-test)
- [Code Style Guidelines](#code-style-guidelines)
- [Huma Framework](#huma-framework)
- [Docker & Environment](#docker--environment)
- [Contributing](#contributing)
- [References](#references)

---

## Purpose

- Kickstart new Go API projects with Huma.
- Encourage clean architecture and modular code.
- Provide working examples for authentication and user management.
- Serve as a foundation for scalable backend development.

---

## Quick Start

```sh
# Clone the repository
git clone https://github.com/Mhirii/rest-template.git
cd rest-template

# Initialize Go modules
go mod tidy

# Update module name in go.mod and import paths as needed

# Run the API server
make run ARGS="--port=8000 --config=config.example.yaml"
```

Visit [http://localhost:8000/docs](http://localhost:8000/docs) for API documentation.

---

## Configuration

- All configuration options are defined in `internal/config/config.go` (`Config` struct).
- Configuration sources (in priority order):
  1. CLI flags (e.g. `--port=8000`)
  2. YAML config file (`config.yaml`)
  3. Environment variables (e.g. `SERVICE_PORT`)
  4. Code defaults

**Example config file:**
```yaml
Server:
  Port: 8000
  LogLevel: debug
DB:
  Host: localhost
  Port: 5432
Auth:
  Secret: secret
  AccessTokenTTL: 86400
  RefreshTokenTTL: 86400
  RateLimit: 1000
```

---

## Project Structure

```
cmd/
  server/
    main.go                # Application entrypoint
internal/
  config/                  # Centralized config loader
  dto/                     # Data Transfer Objects (request/response models)
  handlers/                # HTTP route registration and handler logic
  middleware/              # Middleware (auth, logging, etc.)
  migrations/              # SQL migrations
  models/                  # Database models
  observability/           # OpenTelemetry tracing
docs/
  huma/                    # Official Huma documentation
Makefile                   # Build, run, test commands
docker-compose.yml         # Local development services (Postgres, Redis)
config.example.yaml        # Example configuration
```

---

## API Overview

Auto Generated and Available at /docs, e.g. localhost:8000/docs

**DTOs and models are defined in `internal/dto/` and `internal/models/`.**

---

## Build, Run, and Test

- **Build:** `make build` (binary in `bin/`)
- **Run:** `make run` or `cd cmd/server && go run main.go`
- **Test all:** `make test` or `go test ./...`
- **Test single package:** `go test ./internal/handlers`
- **Test single file:** `go test -run TestFuncName ./internal/handlers/auth.handler.go`
- **Tidy modules:** `make tidy`

---

## Code Style Guidelines

- **Imports:** Group stdlib, third-party, and local packages.
- **Formatting:** Always run `gofmt` or `go fmt ./...` before committing.
- **Types:** Prefer explicit types; use structs for request/response models.
- **Naming:** Use `CamelCase` for types, `mixedCaps` for variables/functions; acronyms are capitalized (e.g., `API`, `ID`).
- **Error Handling:** Return errors, use context-aware logging (`zerolog.Ctx(ctx)`), and prefer Huma error helpers.
- **Comments:** Use Go doc comments for exported types/functions.
- **Concurrency:** Use sync primitives for shared state.
- **File Structure:** Handlers in `internal/handlers`, DTOs in `internal/dto`, services in `internal/service`.
- **Tests:** Place unit tests in the same package; integration/e2e tests in `test/` or scripts.

---

## Huma Framework

This project uses [Huma](https://huma.rocks/), a modern, fast, and flexible micro framework for building HTTP APIs in Go, backed by OpenAPI 3 and JSON Schema.

- Official documentation: [`docs/huma/`](docs/huma/)
- API docs available at `/docs` when the server is running.
- Features include OpenAPI generation, middleware, validation, and more.

---

## Docker & Environment

**docker-compose.yml** provides local development services:
- **Postgres** (`postgres:17-alpine`) on port 5432
- **Redis** (`redis:8.0-alpine`) on port 6379

**Environment variables** can be used for configuration (see `config.example.yaml`).

---

## Contributing

- Follow code style and structure guidelines.
- Update documentation and configuration as needed.
- See [AGENTS.md](AGENTS.md) for contribution conventions.
