# Art-Design-Pro个人后端项目

[前端项目地址](https://github.com/Daymychen/art-design-pro)

# 项目结构

```shell
.
├── README.md                # 项目说明文档，通常包含项目简介、使用方式等
├── build.sh                 # 项目构建脚本，用于编译或部署
├── cmd                      # 命令行程序入口目录
│   └── app                  # 主程序入口（可能是 main.go 或其他启动文件）
│
├── config                   # 配置初始化相关代码
│   ├── initConfig.go        # 配置文件加载逻辑
│   ├── initGorm.go          # GORM（数据库ORM）初始化
│   ├── initHttpServer.go    # HTTP 服务初始化
│   ├── initJwt.go           # JWT 认证初始化
│   ├── initRedis.go         # Redis 客户端初始化
│   ├── initValidator.go     # 请求参数验证器初始化
│   └── initZaplog.go        # Zap 日志库初始化
│
├── configs                  # 配置文件目录（YAML/JSON等）
│   ├── config-prod.yaml     # 生产环境配置
│   └── config.yaml          # 开发环境默认配置
│
├── go.mod                   # Go 模块定义文件（依赖管理）
├── go.sum                   # Go 模块校验文件
│
├── internal                 # 内部实现代码（遵循Go 1.4+的 internal 规则）
│   ├── controller           # HTTP 请求控制器层
│   ├── model                # 数据模型定义（结构体）
│   ├── repository           # 数据持久层（数据库操作）
│   └── service              # 业务逻辑层
│
└── pkg                      # 公共库代码（可被外部引用）
    ├── constant             # 全局常量定义
    ├── errorTypes           # 自定义错误类型
    ├── jwt                  # JWT 相关工具
    ├── loginUtils           # 登录相关辅助函数
    ├── middleware           # HTTP 中间件
    ├── redisx               # Redis 扩展工具
    ├── response             # HTTP 响应格式化
    └── utils                # 通用工具函数
```

# 注意事项

运行前记得生成依赖注入的wire代码

```shell
go tool github.com/google/wire/cmd/wire ./...
```

