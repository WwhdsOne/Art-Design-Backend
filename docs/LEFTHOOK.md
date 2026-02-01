# Lefthook - Git Hooks 管理

本项目使用 [Lefthook](https://github.com/evilmartians/lefthook) 管理 Git hooks，确保代码质量和提交信息规范。

## 功能

### Pre-commit Hook
在提交前自动运行以下检查：
1. ✅ `make wire` - 生成依赖注入代码
2. ✅ `make lint` - 代码检查（revive）
3. ✅ `make test` - 运行测试

### Commit-msg Hook
检查提交信息格式，符合 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

**格式要求：**
```
<type>(<scope>): <subject>
```

或

```
<type>: <subject>
```

**允许的类型（type）：**
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

## 安装

### 首次安装

1. **安装 Lefthook 工具**

```bash
go get -tool github.com/evilmartians/lefthook
```

确保 `$GOPATH/bin` 在你的 PATH 中：

```bash
# macOS/Linux
export PATH=$PATH:$(go env GOPATH)/bin

# Windows
# 将 $(go env GOPATH)\bin 添加到系统 PATH
```

2. **安装 Git hooks**

```bash
./scripts/setup-lefthook.sh
```

或手动运行：

```bash
go tool lefthook install
```

### 验证安装

检查 hooks 是否已安装：

```bash
lefthook run pre-commit
lefthook run commit-msg "feat: test message"
```

## 使用

### 正常提交流程

```bash
git add .
git commit -m "feat: add new feature"
```

Lefthook 会自动：
1. 运行 pre-commit 检查
2. 检查提交信息格式
3. 如果全部通过，提交成功

### 手动运行检查

```bash
# 运行 pre-commit hooks
lefthook run pre-commit

# 测试提交信息检查
lefthook run commit-msg "feat: test message"
```

### 查看配置的 hooks

```bash
lefthook list
```

### 跳过检查（不推荐）

如果确实需要跳过检查，可以使用 Git 的 `--no-verify` 标志：

```bash
# 跳过所有 hooks
git commit --no-verify -m "your message"

# 跳过特定 hook
SKIP=wire git commit -m "your message"
SKIP=lint git commit -m "your message"
SKIP=conventional-commits git commit -m "your message"
```

⚠️ **警告**：跳过检查可能导致低质量代码或不规范的提交信息进入仓库。

## 配置

配置文件位于项目根目录的 `lefthook.yml`：

```yaml
pre-commit:
  parallel: true
  commands:
    wire:
      run: make wire
      glob: "**/*.go"
    lint:
      run: make lint
      glob: "**/*.go"
    test:
      run: make test
      glob: "**/*.go"

commit-msg:
  commands:
    conventional-commits:
      run: ./scripts/lint-commit-msg.sh {1}
```

### 配置说明

- `parallel: true` - 并行执行命令，提高速度
- `glob: "**/*.go"` - 只对 Go 文件运行检查
- `{1}` - commit-msg hook 的第一个参数是提交信息文件路径

## 自定义

### 添加新的 hook

在 `lefthook.yml` 中添加新的 hook 配置：

```yaml
pre-push:
  commands:
    full-test:
      run: go test -race ./...
```

### 修改提交信息规则

编辑 `scripts/lint-commit-msg.sh` 文件，修改正则表达式或类型列表。

### 禁用某个 hook

在 `lefthook.yml` 中注释掉对应的配置：

```yaml
pre-commit:
  parallel: true
  commands:
    # wire:
    #   run: make wire
    #   glob: "**/*.go"
    lint:
      run: make lint
      glob: "**/*.go"
    test:
      run: make test
      glob: "**/*.go"
```

## 常见问题

### Lefthook 命令不存在

**问题**：`command not found: lefthook`

**解决**：
1. 确认已安装：`go list -m github.com/evilmartians/lefthook`
2. 检查 PATH：`echo $PATH | grep $(go env GOPATH)/bin`
3. 重新安装：`go get -tool github.com/evilmartians/lefthook`

### Hook 没有生效

**问题**：提交时没有运行检查

**解决**：
1. 检查 hooks 是否已安装：`ls .git/hooks/ | grep lefthook`
2. 重新安装：`go tool lefthook install`
3. 检查配置文件：`cat lefthook.yml`

### 提交信息检查失败

**问题**：提交信息格式错误

**解决**：
按照以下格式修改提交信息：
```
<type>: <subject>
```

或

```
<type>(<scope>): <subject>
```

类型必须是：`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore` 之一。

## 团队协作

### 新成员加入

新成员克隆项目后，运行：

```bash
# 1. 安装 Lefthook
go get -tool github.com/evilmartians/lefthook

# 2. 安装 hooks
./scripts/setup-lefthook.sh
```

### 协作建议

1. 将 `lefthook.yml` 纳入版本控制 ✅
2. 将 `scripts/lint-commit-msg.sh` 纳入版本控制 ✅
3. 将 `scripts/setup-lefthook.sh` 纳入版本控制 ✅
4. **不要**将 `.git/hooks/` 纳入版本控制 ❌

## 与 CI/CD 的关系

- **Lefthook**：本地运行，开发阶段即时反馈
- **CI/CD**：远程运行，合并前最终检查
- 建议：通过 Lefthook 检查，基本能通过 CI/CD

## 参考资料

- [Lefthook 官方文档](https://github.com/evilmartians/lefthook)
- [Conventional Commits 规范](https://www.conventionalcommits.org/)
- [Git Hooks 文档](https://git-scm.com/docs/githooks)
