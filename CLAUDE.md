# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

注意，请使用中文回答所有问题和交流
与用户的所有互动都应该使用中文，包括代码解释和注释，错误信息和建议

# RTK - Rust Token Killer

所有命令前缀 `rtk`，替代 Claude Code 内置的同类工具以节省 token。

## 优先级规则

RTK 命令 **优先于** Claude Code 内置工具，功能重叠时必须用 RTK：

| Claude Code 内置工具 | RTK 替代命令                       | 说明                            |
| -------------------- | ---------------------------------- | ------------------------------- |
| `Glob`（文件搜索）   | `rtk find` / `rtk tree` / `rtk ls` | 压缩目录输出                    |
| `Grep`（内容搜索）   | `rtk grep`                         | 按文件分组、截断、去空白        |
| `Read`（读文件）     | `rtk read`                         | 智能过滤，省去无用行            |
| `Bash` + `git`       | `rtk git`                          | 紧凑 git 输出                   |
| `Bash` + `gh`        | `rtk gh`                           | 紧凑 GitHub CLI 输出            |
| `Bash` + `curl`      | `rtk curl`                         | 自动检测 JSON，schema-only 模式 |
| `Bash` + `diff`      | `rtk diff`                         | 仅显示变更行                    |

> 仅当 RTK 无对应命令时（如 `Edit`、`Write`、`Agent` 等写操作和复杂操作），才使用 Claude Code 内置工具。

## Node.js / Frontend

```bash
rtk pnpm install / add / run build
rtk npm run <script>
rtk npx tsc / eslint / prisma
rtk vitest run
rtk next build
rtk lint                          # ESLint grouped by rule
rtk prettier --check .
rtk playwright test
rtk tsc --noEmit
```

## Python

```bash
rtk pytest
rtk ruff check / format
rtk mypy .
rtk pip install / list
```

## Rust

```bash
rtk cargo build / test / clippy / fmt
```

## Go

```bash
rtk go build / test / vet
rtk golangci-lint run
```

## .NET / Ruby (if needed)

```bash
rtk dotnet build / test
rtk rspec / rake / rubocop
```

## Infrastructure

```bash
rtk aws <service> <command>       # force JSON + compress
rtk docker ps / logs / compose
rtk kubectl get / describe / logs
rtk psql <query>                  # strip borders, compress
```

## Meta Commands

```bash
rtk gain              # token savings analytics
rtk gain --history    # usage history with savings
rtk discover          # find missed opportunities in session history
```

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
  │   ├── *.go                    # Repository interfaces (Service 层只依赖此层)
  │   ├── db/                     # Database repositories (GORM)
  │   ├── cache/                  # Cache repositories (Redis)
  │   ├── cachex/                 # 泛型缓存工具 GetOrNil/GetSlice/Set
  │   ├── dupcheck/               # 泛型去重检查工具
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
- **Repository interface** (`internal/repository/*.go`) - Combines DB + Cache, **Service 层只依赖 `repository` 包**（`db.GormTX` 除外）
- **Generic utilities** (`internal/repository/cachex/`, `internal/repository/dupcheck/`) - 泛型缓存/去重工具

**Important**: Service 层统一通过 `repository` 包访问数据，不直接依赖 `db/*`（`db.GormTX` 事务管理器除外）。

Example:
```go
type UserRepo struct {
    *db.UserDB
    *db.UserRolesDB
    *cache.UserCache
    *cache.RoleCache
}
```

**cachex 注意事项**：`cachex.GetOrNil`/`GetSlice` 在缓存未命中时返回 `nil, redis.Nil`，调用方必须用 `errors.Is(err, redis.Nil)` 判断缓存 miss，而非检查 `err == nil`。

## Core Features

### Browser Agent
The Browser Agent is an LLM-powered browser automation system:
- **WebSocket-based** real-time communication with browser clients
- **LLM Decision Loop**: LLM analyzes page state → generates action → executes → repeats
- **Vision Fallback**: Text model (DeepSeek) can request `need_vision: true`, triggering Qwen 3.5 Flash vision model with labeled screenshots
- **Supported actions**: goto, click, input, select, scroll, wait, close_browser
- **Models**: DeepSeek (text decision), Qwen 3.5 Flash via 通义千问 (vision), Zhipu AI (legacy)

Key files:
- `internal/service/browser_agent.go` - Core agent logic + vision fallback
- `internal/controller/browser_agent.go` - WebSocket endpoint
- `pkg/ws/` - WebSocket hub for client connections
- `pkg/constant/llmid/` - Model ID constants
- `pkg/constant/prompt/` - System prompts + vision prompt

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

6. **Message sorting (2026-04-10 fixed)**: `ListMessagesByConversationID` now returns messages in `created_at ASC` order (oldest first, newest last) to match natural conversation flow. The related `ListMessagesPage` uses `DESC` order for pagination, which is intentional.

## Go 版本新特性（1.25 & 1.26）

项目当前使用 Go 1.26。以下是 1.25（2025.8）和 1.26（2026.2）中常见的、可在本项目中使用的 API 和语法更新。写代码时优先使用新写法。

### Go 1.25 新特性

#### `testing/synctest` 包（稳定版）
用于测试涉及并发代码（goroutine、channel、定时器）的场景，可精确控制时间推进。
```go
import "testing/synctest"

func TestWithTime(t *testing.T) {
    synctest.Test(t, func(d time.Duration) {
        // d 是虚拟时钟，可精确控制时间流逝
        // 适用于测试超时、重试、缓存过期等逻辑
    })
}
```

#### `sync.WaitGroup.Go` 方法
简化 goroutine 启动模式，替代 `wg.Add(1); go func() { defer wg.Done(); ... }()`。
```go
var wg sync.WaitGroup
wg.Go(func() {
    // 自动处理 Add/Done
    doWork()
})
wg.Wait()
```

#### Container 感知的 `GOMAXPROCS`
在 Docker/Kubernetes 中运行时，`runtime.GOMAXPROCS(0)` 自动感知 cgroup CPU 限制，不再需要手动设置 `automaxprocs`。

#### `os.Root` 类型
提供安全的文件系统操作根目录，防止路径穿越（path traversal）。
```go
root, _ := os.OpenRoot("/data")
f, _ := root.Open("safe/path.txt")   // 不允许 "../" 等逃逸路径
```

#### `reflect.TypeAssert` 方法
类型断言的反射版本，比传统的 `Type.Implements()` 更直接。

### Go 1.26 新特性

#### 表达式 `new()` — 最实用
`new()` 现在可以接受任意表达式，不再仅限于类型名。大幅简化复合字面量中的指针创建。
```go
// 旧写法
x := 42
p := &x

// 新写法 — 直接 new 表达式
p := new(42)
p := new(3.14)
p := new(someFunc())
p := new(User{Name: "test"})
```

#### 自引用泛型
类型参数可以引用自身，实现递归类型约束。
```go
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

#### `errors.AsType` — 泛型错误匹配
`errors.As` 的泛型版本，无需声明变量、无需类型断言。
```go
// 旧写法
var timeoutErr *net.OpError
if errors.As(err, &timeoutErr) {
    fmt.Println(timeoutErr.Op)
}

// 新写法 — 更简洁
if timeoutErr, ok := errors.AsType[*net.OpError](err); ok {
    fmt.Println(timeoutErr.Op)
}
```

#### Green Tea GC（默认启用）
新一代垃圾回收器，降低 CPU 开销 10-40%。无需代码改动，自动生效。

#### `reflect` 迭代器
新增 `Type.Fields()`、`Type.Methods()`、`Value.Fields()` 等迭代器方法，替代繁琐的手动遍历。
```go
// 遍历结构体字段
for field, ok := t.Fields(); ok; {
    fmt.Println(field.Name)
}
```

#### `io.ReadAll` 性能提升
`io.ReadAll` 速度提升约 2x，无需任何代码改动。

#### `bytes.Buffer.Peek` 方法
`bytes.Buffer` 新增 `Peek(n)` 方法，查看缓冲区前 n 字节但不消耗。
```go
buf := bytes.NewBuffer([]byte("hello world"))
peek, _ := buf.Peek(5)  // "hello"，buf 中的数据不被消耗
```

#### `log/slog.NewMultiHandler`
将日志同时输出到多个 handler。
```go
handler := slog.NewMultiHandler(fileHandler, consoleHandler)
logger := slog.New(handler)
```

#### `go fix` 现代化工具
`go fix` 全面升级，内置 modernizers 自动将旧写法更新为新写法（如 `errors.AsType` 替换 `errors.As`）。项目 `make fix` 已集成。
