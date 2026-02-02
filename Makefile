# Makefile for Art-Design-Backend
# Code quality, dependency management and build

.PHONY: help install-tools lint lint-full test build tidy upgrade upgrade-patch wire pre-commit check

help:
	@echo "可用命令:"
	@echo "  make install-tools  - 安装 go.mod 中声明的所有工具"
	@echo "  make lint           - 快速代码检查（revive）"
	@echo "  make lint-full      - 全量代码检查（golangci-lint）"
	@echo "  make wire           - 生成依赖注入代码（wire）"
	@echo "  make test           - 运行测试"
	@echo "  make build          - 构建项目"
	@echo "  make tidy           - 整理依赖（不升级）"
	@echo "  make upgrade        - 升级全部依赖（⚠️ 有破坏性）"
	@echo "  make upgrade-patch  - 升级全部依赖（仅 patch 版本，推荐）"
	@echo "  make pre-commit     - 提交前检查"
	@echo "  make check          - 本地完整检查"

# 安装所有 tool 依赖（Go 1.22+）
install-tools:
	@echo "安装开发工具..."
	go get -tool github.com/mgechev/revive
	go get -tool github.com/google/wire/cmd/wire
	go get -tool github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -tool github.com/evilmartians/lefthook
	@echo "工具安装完成！"

# 快速 lint（revive）
lint:
	@echo "运行 revive..."
	go tool revive -config revive.toml ./...

# 全量 lint（golangci-lint）
# todo 日后逐步支持golangci-lint
lint-full:
	@echo "运行 golangci-lint..."
	go tool golangci-lint run ./...

# 依赖注入代码生成（Wire）
wire:
	@echo "运行 wire 生成依赖注入代码..."
	go tool wire ./...

test:
	@echo "运行测试..."
	go test ./...

build:
	@echo "构建项目..."
	bash ./scripts/build.sh

# 整理依赖（安全）
tidy:
	@echo "整理依赖..."
	go mod tidy

# ⚠️ 升级全部依赖（主版本 / 次版本 / patch）
upgrade:
	@echo "升级全部依赖（可能破坏兼容性）..."
	go get -u ./...
	go mod tidy

# ✅ 推荐：只升级 patch 版本
upgrade-patch:
	@echo "升级全部依赖（仅 patch 版本）..."
	go get -u=patch ./...
	go mod tidy

pre-commit: wire lint test
