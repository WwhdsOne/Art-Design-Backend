#!/usr/bin/env bash
# ❗ 一旦脚本中的任意一条命令返回非零状态（即执行失败），整个脚本立即终止执行。
set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 生成 wire 依赖注入代码
go run github.com/google/wire/cmd/wire ./...

# 配置
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

# 构建参数
LD_FLAGS="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
TAGS="sonic,avx"

# 执行构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "${LD_FLAGS}" \
    -tags "${TAGS},netgo,osusergo" \
    -o "${APP_NAME}" \
    ./cmd/app

# 使用 UPX 压缩
upx --lzma --best "${APP_NAME}"
