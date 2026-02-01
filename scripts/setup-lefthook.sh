#!/bin/bash
# Lefthook 安装脚本

set -e

echo "🔧 安装 Lefthook..."

# 检查 lefthook 是否已安装
if ! command -v lefthook &> /dev/null; then
    echo "📦 Lefthook 未安装，正在安装..."

    # 尝试使用 go install
    if command -v go &> /dev/null; then
        go get -tool github.com/evilmartians/lefthook
        echo "✅ Lefthook 安装成功"
        echo "📝 请确保 \$GOPATH/bin 在你的 PATH 中"
    else
        echo "❌ 错误：未找到 Go，无法安装 Lefthook"
        echo "请先安装 Go: https://golang.org/dl/"
        exit 1
    fi
else
    echo "✅ Lefthook 已安装"
fi

# 安装 hooks
echo "🔗 安装 Git hooks..."
if command -v lefthook &> /dev/null; then
    go tool lefthook install
    echo "✅ Git hooks 安装成功"
else
    echo "⚠️  Lefthook 未在 PATH 中，尝试使用完整路径..."

    # 尝试找到 lefthook 二进制
    LEFTHOOK_PATH="$GOPATH/bin/lefthook"
    if [ -z "$GOPATH" ]; then
        LEFTHOOK_PATH="$HOME/go/bin/lefthook"
    fi

    if [ -f "$LEFTHOOK_PATH" ]; then
        "$LEFTHOOK_PATH" install
        echo "✅ Git hooks 安装成功"
    else
        echo "❌ 错误：无法找到 Lefthook"
        exit 1
    fi
fi

echo ""
echo "🎉 Lefthook 配置完成！"
echo "📝 配置文件: lefthook.yml"
echo ""
echo "现在以下 hooks 将自动运行："
echo "  - pre-commit: wire → lint → test"
echo "  - commit-msg: 检查提交信息格式"
