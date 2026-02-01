# 项目代码质量检查

本项目使用完整的代码质量检查体系，确保代码质量和团队协作规范。

## 检查体系

### 1. Lefthook - Git Hooks 管理

自动化 Git hooks，在提交前进行检查。

**文档：** [LEFTHOOK.md](LEFTHOOK.md)

**功能：**
- Pre-commit: wire → lint → test
- Commit-msg: 提交信息格式检查（Conventional Commits）

### 2. Revive - Go 代码检查

使用 Revive 进行代码质量检查。

**配置文件：** `revive.toml`

**运行：**
```bash
make lint
```

### 3. Makefile - 任务管理

统一的构建和检查命令。

**常用命令：**
```bash
make help           # 查看所有命令
make lint           # 快速代码检查
make test           # 运行测试
make build          # 构建项目
make pre-commit     # 提交前检查（包含 wire、lint、test）
make check          # 本地完整检查（格式化+lint+test+构建）
```

### 4. GitHub Actions - CI/CD

自动化 CI/CD 流程，确保代码在合并前通过检查。

**配置文件：** `.github/workflows/`

**触发条件：**
- Push 到分支
- Pull Request

## 提交流程

### 开发阶段

1. 编写代码
2. 运行 `make pre-commit` 检查（可选，hook 会自动运行）
3. 提交：`git commit -m "feat: your message"`
4. Lefthook 自动运行检查

### 推送阶段

1. 推送到远程：`git push`
2. GitHub Actions 自动运行 CI/CD
3. 检查通过后，可以合并 PR

## 代码规范

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>
```

**类型（type）：**
- `feat` - 新功能
- `fix` - 修复 bug
- `docs` - 文档修改
- `style` - 格式（不影响代码逻辑）
- `refactor` - 重构
- `perf` - 性能优化
- `test` - 增加测试
- `chore` - 构建过程或辅助工具的变动

**示例：**
```bash
feat(auth): add JWT authentication
fix: resolve memory leak in user service
docs: update API documentation
style: format code with revive
```

### 代码风格

遵循 Go 官方代码规范：
- 使用 `gofmt` 格式化
- 使用 `revive` 检查代码质量
- 导入排序：goimports
- 注释清晰，必要时添加中文注释

## 检查失败处理

### 本地检查失败（Lefthook）

1. 查看错误信息
2. 修复代码问题
3. 重新提交
4. 或运行 `make pre-commit` 查看详细错误

### CI/CD 检查失败（GitHub Actions）

1. 查看 GitHub Actions 日志
2. 本地复现问题
3. 修复后重新推送
4. 确保本地能通过 `make check`

## 跳过检查（不推荐）

### 跳过 Lefthook

```bash
git commit --no-verify -m "your message"
```

### 跳过特定 Hook

```bash
SKIP=wire git commit -m "your message"
SKIP=lint git commit -m "your message"
SKIP=conventional-commits git commit -m "your message"
```

⚠️ **警告**：跳过检查可能导致低质量代码进入仓库。

## 工具安装

### Lefthook

```bash
go get -tool github.com/evilmartians/lefthook
./scripts/setup-lefthook.sh
```

### Revive

```bash
go get -tool github.com/mgechev/revive
```

### Wire

```bash
go get -tool github.com/google/wire/cmd/wire
```

### 所有工具

```bash
make install-tools
```

## 配置文件

| 文件 | 说明 |
|------|------|
| `lefthook.yml` | Lefthook 配置 |
| `revive.toml` | Revive 代码检查配置 |
| `Makefile` | 构建和检查任务 |
| `.github/workflows/` | CI/CD 配置 |

## 参考文档

- [Lefthook 使用指南](LEFTHOOK.md)
- [Revive 配置说明](https://github.com/mgechev/revive#configuration)
- [Conventional Commits 规范](https://www.conventionalcommits.org/)
- [Go 代码规范](https://go.dev/doc/effective_go)
