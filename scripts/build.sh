#!/usr/bin/env bash
set -e

# 切换到项目根目录（假设脚本在scripts目录）
cd "$(dirname "$0")/.."

# 配置
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date '+%Y-%m-%d_%H:%M:%S')

# 检测当前平台
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

echo "🔹 构建信息: OS=$OS ARCH=$ARCH VERSION=$VERSION BUILD_TIME=$BUILD_TIME"

# 构建参数，启用 Greentea GC
CGO_ENABLED=0 \
GOEXPERIMENT=greenteagc \
go build \
  -trimpath \
  -ldflags "-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
  -o "bin/${APP_NAME}" \
  ./cmd/app

# UPX 压缩（仅 Linux/Windows，macOS 跳过）
if command -v upx >/dev/null 2>&1; then
    if [[ "$OS" == "darwin" ]]; then
        echo "⚠️ macOS 不支持 UPX 压缩，跳过..."
    else
        echo "🔧 使用 UPX 压缩可执行文件..."
        upx --lzma --best "bin/${APP_NAME}" || true
    fi
else
    echo "⚠️ 未找到 UPX，跳过压缩。"
fi

echo "✅ 构建完成，输出：bin/${APP_NAME}"
