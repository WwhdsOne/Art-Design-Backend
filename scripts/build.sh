#!/usr/bin/env bash
set -e

# 切换到项目根目录（假设脚本在scripts目录）
cd "$(dirname "$0")/.."

# 配置
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

# 构建参数
LD_FLAGS="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
TAGS="sonic,avx"

# 构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
    -trimpath \
    -buildvcs=false \
    -ldflags "${LD_FLAGS}" \
    -tags "${TAGS},netgo,osusergo" \
    -o "${APP_NAME}" \
    ./cmd/app

if command -v upx >/dev/null 2>&1; then
    echo "🔧 使用 UPX 压缩可执行文件..."
    upx --lzma --best "${APP_NAME}"
else
    echo "⚠️  未找到 upx，跳过压缩。"
fi
