# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build & Test
```bash
make wire            # Generate dependency injection code (Wire)
make build           # Build the project
make test            # Run tests
```

### Code Quality
```bash
make fix             # Fix code using go fix (Go 1.26+)
make lint            # Quick lint with revive
make lint-full       # Full lint with golangci-lint
make pre-commit      # Run all pre-commit checks (fix, wire, lint, test)
```

### Dependency Management
```bash
make install-tools   # Install all tools declared in go.mod
make tidy            # Tidy dependencies (no upgrades)
make upgrade-patch   # Upgrade dependencies (patch versions only, recommended)
make upgrade         # Upgrade all dependencies (may break compatibility)
```

### Running the Application
```bash
make wire            # First, generate Wire code
make build           # Build the binary
./bin/art-design-backend  # Run the server
```

## Project Architecture

### Dependency Injection (Wire)
This project uses **Google Wire** for dependency injection. Wire files are located at:
- `internal/bootstrap/wire_bootstrap.go`
- `internal/repository/wire_repository.go`
- `internal/controller/wire_controller.go`
- `cmd/app/wire.go`

**IMPORTANT:** After modifying any Wire provider functions or adding new dependencies, run `make wire` to regenerate `wire_gen.go` files.

### Layered Architecture
```
cmd/app/                          # Application entry point
  ├── main.go                     # Entry point
  ├── wire.go                     # Wire injector definition
  └── wire_gen.go                 # Generated Wire code (DO NOT EDIT)

internal/
  ├── bootstrap/                  # Initialization providers
  │   ├── init_*.go               # Individual component initialization
  │   └── wire_bootstrap.go       # Wire providers for bootstrap
  ├── controller/                 # HTTP handlers (routing)
  │   ├── *.go                    # Individual controllers
  │   └── wire_controller.go      # Wire providers for controllers
  ├── service/                    # Business logic layer
  │   └── *.go                    # Service implementations
  ├── repository/                 # Data access layer
  │   ├── db/                     # Database repositories (GORM)
  │   ├── cache/                  # Cache repositories (Redis)
  │   └── wire_repository.go      # Wire providers for repositories
  └── model/                      # Domain models

pkg/                              # Shared packages
  ├── middleware/                 # Gin middleware
  ├── ai/                         # AI client abstraction
  ├── redisx/                     # Redis wrapper
  ├── ws/                         # WebSocket hub
  ├── errors/                     # Custom errors
  ├── result/                     # HTTP response helpers
  └── utils/                      # Utilities

config/                           # Configuration loading (Consul)
```

### Controller Pattern
Controllers are responsible for:
1. Registering routes in their `New*Controller` constructor
2. Binding request data to structs
3. Calling service layer
4. Returning responses using `result` package

Example from `internal/controller/auth.go`:
```go
func NewAuthController(engine *gin.Engine, mws *middleware.Middlewares, svc *service.AuthService) *AuthController {
    authCtrl := &AuthController{authService: svc}
    r := engine.Group("/api").Group("/auth")
    {
        r.POST("/login", authCtrl.login)
        r.POST("/logout", mws.AuthMiddleware(), authCtrl.logout)
    }
    return authCtrl
}
```

### Service Layer
Services contain business logic and are injected with repositories. Services should:
- Use transaction managers for multi-step database operations
- Return domain errors (from `pkg/errors`)
- Not depend on Gin context

### Repository Pattern
Repositories split concerns between:
- **DB repositories** (`internal/repository/db/`) - PostgreSQL via GORM
- **Cache repositories** (`internal/repository/cache/`) - Redis caching
- **Repository interface** (`internal/repository/*.go`) - Combines DB + Cache

Example:
```go
type UserRepo struct {
    UserDB    *db.UserDB
    UserCache *cache.UserCache
    RoleCache *cache.RoleCache  // Can use multiple caches
}
```

## Core Features

### Browser Agent
The Browser Agent is an LLM-powered browser automation system:
- **WebSocket-based** real-time communication with browser clients
- **LLM Decision Loop**: LLM analyzes page state → generates action → executes → repeats
- **Supported actions**: goto, click, input, select, scroll, wait, close_browser
- **Models**: DeepSeek, Zhipu AI for decision-making

Key files:
- `internal/service/browser_agent.go` - Core agent logic
- `internal/controller/browser_agent.go` - WebSocket endpoint
- `pkg/ws/` - WebSocket hub for client connections

### AI Services
Multi-provider AI abstraction supporting:
- **Providers**: Zhipu AI, DeepSeek, Tongyi Qianwen
- **Models**: Chat completion, embeddings (text-embedding-v4, 1024 dimensions)
- **Features**: Streaming responses (SSE), conversation history

### Knowledge Base (RAG)
Document Q&A system with:
- **Vector search** using pgvector (PostgreSQL extension)
- **Hybrid search**: Vector + keyword retrieval
- **Reranking**: SiliconFlow Rerank API
- **Chunking**: External slicer service

### Authentication & Authorization
- **RBAC**: User → Role → Menu permissions
- **JWT**: Token-based auth with Redis session storage
- **Middleware**: `AuthMiddleware()` for protected routes
- **Password validation**: Custom `strongpassword` validator (uppercase + lowercase + numbers)

## Configuration

Uses **Consul** for configuration management:
1. Start Consul: `docker run -d -p 8500:8500 consul`
2. Set environment variables:
   ```bash
   export CONSUL_ADDR=localhost:8500
   export CONSUL_CONFIG_KEY=art-design-backend
   ```
3. Upload config to Consul key `art-design-backend`
4. Example config: `configs/config.example.yaml`

Key config sections:
- `server`: HTTP server settings (port, timeouts)
- `postgre_sql`: PostgreSQL connection
- `redis`: Redis connection
- `jwt`: JWT signing key and expiration
- `ai`: AI provider API keys and base URLs

## Git Hooks (Lefthook)

Pre-commit hooks run automatically:
1. `make fix` - Auto-fix code issues
2. `make wire` - Regenerate Wire code
3. `make lint` - Run revive linter
4. `make test` - Run tests

Commit messages must follow [Conventional Commits](https://www.conventionalcommits.org/):
- Format: `<type>(<scope>): <subject>`
- Types: feat, fix, docs, style, refactor, perf, test, chore

## Database

### Technologies
- **PostgreSQL** with `pgvector` extension for vector operations
- **Redis** for caching and session storage
- **GORM** as ORM

### Migrations
No automatic migration tool is configured. Migrations are likely manual or SQL-based.

### Transaction Management
Use `db.GormTransactionManager` for multi-repository transactions:
```go
txManager.GormTransaction(func(tx *gorm.DB) error {
    // Multiple operations using tx
    return nil
})
```

## Testing

Tests use standard Go testing. Run with:
```bash
make test
```

Note: As of the last check, no test files exist in the codebase.

## Linting

Two linters are configured:
1. **Revive** (`make lint`) - Fast, recommended for development
2. **golangci-lint** (`make lint-full`) - Comprehensive, slower

Configuration: `.golangci.yml` and `revive.toml`

## Common Gotchas

1. **Wire not regenerated**: Always run `make wire` after:
   - Adding new dependencies to constructors
   - Modifying provider functions
   - Creating new controllers/services/repositories

2. **Consul config**: The app will fail to start without Consul configuration. Ensure:
   - Consul is running
   - Environment variables are set
   - Config is uploaded to Consul

3. **Route registration**: Routes are registered in controller constructors (`New*Controller`), not in a central router file

4. **Middleware order**: Defined in `internal/bootstrap/init_gin.go`:
   - Gzip → Logger → Recovery → ErrorHandler → OperationLogger → RateLimiter → Auth (per-route)

5. **Strong password validator**: Custom validator requires passwords to have uppercase, lowercase, and numbers
