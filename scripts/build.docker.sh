#!/usr/bin/env bash
set -e

# 切换到项目根目录（假设脚本在scripts目录）
cd "$(dirname "$0")/.."

# 设置代理
export GOPROXY=https://goproxy.cn,direct

# 更新依赖
go get -u ./... && go mod tidy

# 生成 wire 依赖注入代码
go run github.com/google/wire/cmd/wire ./...

# 配置
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

# 构建参数
LD_FLAGS="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
TAGS="sonic,avx"

# 构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "${LD_FLAGS}" \
    -tags "${TAGS},netgo,osusergo" \
    -o "${APP_NAME}" \
    ./cmd/app
