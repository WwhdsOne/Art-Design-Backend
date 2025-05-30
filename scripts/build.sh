#!/usr/bin/env bash
# ❗ 一旦脚本中的任意一条命令返回非零状态（即执行失败），整个脚本立即终止执行。
set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 配置
APP_NAME="myapp"
LD_FLAGS="-w -s"
TAGS="sonic,avx"

# 检测平台是否为 Linux
if [[ "$(uname)" == "Linux" ]]; then
    echo "🧠 Linux 平台，限制为单核编译"
    export GOMAXPROCS=1
fi

# 构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
    -trimpath \
    -buildvcs=false \
    -ldflags "${LD_FLAGS}" \
    -tags "${TAGS},netgo,osusergo" \
    -o "${APP_NAME}" \
    ./cmd/app

# 如果安装了 upx，则执行压缩
if command -v upx >/dev/null 2>&1; then
    echo "🔧 使用 UPX 压缩可执行文件..."
    upx --lzma --best "${APP_NAME}"
else
    echo "⚠️  未找到 upx，跳过压缩。"
fi
