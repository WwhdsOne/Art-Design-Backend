#!/bin/bash
# 提交信息格式检查
# 符合 Conventional Commits 规范

# 获取所有参数（处理空格问题）
# shellcheck disable=SC2124
ARGS="$@"

# 判断参数是文件路径还是直接的提交消息
if [ -f "$1" ] && [ $# -eq 1 ]; then
    # Git hook 会传递文件路径
    COMMIT_MSG=$(cat "$1")
else
    # 测试或手动运行时，直接使用所有参数作为提交消息
    COMMIT_MSG="$ARGS"
fi

# 定义允许的提交类型
TYPES="feat|fix|docs|style|refactor|perf|test|chore"

# 正则表达式匹配 Conventional Commits 格式
# 格式: type(scope): subject
# 示例: feat(auth): add JWT authentication
PATTERN="^($TYPES)(\(.+\))?!?: .+"

# 检查是否符合规范
if [[ ! "$COMMIT_MSG" =~ $PATTERN ]]; then
    cat << EOF
❌ 提交信息格式错误！

提交信息必须符合 Conventional Commits 规范，格式如下：

  <type>(<scope>): <subject>

或

  <type>: <subject>

允许的类型（type）：
  - feat     新功能
  - fix      修复 bug
  - docs     文档修改
  - style    格式（不影响代码逻辑）
  - refactor 重构
  - perf     性能优化
  - test     增加测试
  - chore    构建过程或辅助工具的变动

示例：
  ✓ feat(auth): add JWT authentication
  ✓ fix: resolve memory leak in user service
  ✓ docs: update API documentation
  ✓ style: format code with revive

你的提交信息：
  $COMMIT_MSG

💡 提示：请修改提交信息后再次尝试提交
EOF
    exit 1
fi

# 可选：检查标题长度（72 字符限制，默认注释）
# SUBJECT=$(echo "$COMMIT_MSG" | head -n1)
# if [ ${#SUBJECT} -gt 72 ]; then
#     echo "⚠️  警告：提交信息标题超过 72 字符，建议缩短"
# fi

echo "✅ 提交信息格式检查通过"
exit 0
