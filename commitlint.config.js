module.exports = {
    extends: ['@commitlint/config-conventional'],
    rules: {
        'type-enum': [
            2,
            'always',
            [
                'feat',     // 新功能
                'fix',      // 修复 bug
                'docs',     // 文档修改
                'style',    // 格式（不影响代码逻辑）
                'refactor', // 重构
                'perf',     // 性能优化
                'test',     // 增加测试
                'chore',    // 构建过程或辅助工具的变动
            ],
        ],
        // 如果你想更严格一点，也可以开启以下规则
        // 'header-max-length': [2, 'always', 72],
    },
};
