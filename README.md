# Art-Design-Pro个人后端项目

[前端项目地址](https://github.com/Daymychen/art-design-pro)

# 项目结构

```shell
.
├── LICENSE                         # 开源许可文件
├── README.md                       # 项目说明文档
├── build.sh                        # 构建或部署脚本
├── cmd                             # 主程序入口目录
│   └── app                         # 启动应用的 main 函数等入口代码
├── commitlint.config.js           # commitlint 配置文件（校验 Git 提交信息）
├── config                          # 应用初始化配置（如数据库、Redis、日志等初始化逻辑）
├── configs                         # 配置文件目录（如 YAML 格式的配置信息）
├── go.mod                          # Go 依赖管理文件
├── go.sum                          # Go 依赖的完整校验和
├── internal                        # 核心业务逻辑（controller、service、repository、model 等）
│   ├── controller                 # 路由控制层（处理 HTTP 请求）
│   ├── model                      # 数据模型定义（Entity、Request、Response、Query 等）
│   ├── repository                 # 数据访问层（对数据库的操作）
│   └── service                    # 服务层（处理业务逻辑）
├── package.json                   # Node 项目配置（用于前端工具 Git 钩子，Husky、Commitlint）
├── pkg                             # 通用库模块（可复用的工具库或封装组件）
│   ├── aliyun                    # 阿里云 OSS 客户端封装
│   ├── authutils                 # 权限认证相关工具
│   ├── client                    # 外部服务客户端（如推理服务）
│   ├── constant                  # 常量定义（如 Redis key、表名等）
│   ├── container                 # 并发安全的容器封装
│   ├── errors                    # 错误类型定义
│   ├── jwt                       # JWT 生成与验证
│   ├── middleware                # Gin 中间件（如认证、限流、日志记录等）
│   ├── redisx                    # Redis 操作封装（支持 Lua、集合、缓存等功能）
│   ├── result                    # 通用响应格式封装
│   └── utils                     # 常用工具方法（如 UUID、时间解析等）
```

# 注意事项

运行前记得生成依赖注入的wire代码

```shell
go get -u ./... && go mod tidy && go tool github.com/google/wire/cmd/wire ./...
```

