# Art-Design-Pro个人后端项目

[前端项目地址](https://github.com/Daymychen/art-design-pro)

# 项目结构

```shell
├── cmd                         # 应用程序入口目录（通常包含 main.go）
│   └── app                     # 具体的应用主程序
├── config                      # 配置加载逻辑（如初始化配置结构体等）
├── configs                     # 配置文件目录（如 .yaml、.json 等）
├── internal                    # 内部模块（按领域或功能划分）
│   ├── bootstrap               # 项目启动流程，如初始化数据库、日志、依赖注入等
│   ├── controller              # 控制器层（HTTP 接口逻辑）
│   ├── model                   # 模型定义层
│   │   ├── base                # 通用基础模型（如 BaseModel）
│   │   ├── entity              # 与数据库结构对应的实体定义
│   │   ├── query               # 查询结构体定义（用于参数组合、查询构造）
│   │   ├── request             # 接收前端请求的结构体
│   │   └── response            # 返回给前端的响应结构体
│   ├── repository              # 仓储层（数据访问逻辑）
│   │   ├── cache               # 缓存访问逻辑（通常是 Redis）
│   │   └── db                  # 数据库访问逻辑（通常是 GORM/SQL）
│   └── service                 # 服务层（业务逻辑实现）
├── pkg                         # 可复用的通用模块（第三方或自研工具）
│   ├── ai                      # AI 模块（如模型推理、调用接口）
│   ├── aliyun                 # 阿里云 SDK 封装（如短信、OSS）
│   ├── authutils               # 鉴权工具包
│   ├── constant                # 常量定义
│   │   ├── rediskey            # Redis 键名常量
│   │   └── tablename           # 表名常量
│   ├── digit_client            # 数字识别客户端（如识图接口等）
│   ├── errors                  # 自定义错误类型
│   ├── jwt                     # JWT 生成与解析
│   ├── middleware              # Gin 中间件集合
│   ├── redisx                  # Redis 封装（连接池、通用方法）
│   ├── result                  # 通用响应结构体（如统一的 Response 封装）
│   └── utils                   # 工具函数（如字符串处理、时间格式化等）
└── scripts                     # 启动脚本、部署脚本、数据库初始化等
```

# 完整技术栈

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'primaryColor': '#ffdfd3', 'edgeLabelBackground':'#fff', 'fontFamily': '"Microsoft YaHei", sans-serif'}}}%%
graph TD
    subgraph 外部服务
        A[PostgreSQL]
        B[Redis]
        C[阿里云OSS]
        D[AI模型服务]
    end

    subgraph ArtDesignBackend
        subgraph 基础设施层
            Config[配置管理]
            Logger[Zap日志系统]
            DB[GORM数据库]
            Cache[Redis客户端]
            OSS[OSS客户端]
            AIClient[AI服务客户端]
        end

        subgraph 核心层
            subgraph 领域层
                Models[领域模型]
                Repositories[数据仓库]
                Services[业务服务]
            end

            subgraph 接口层
                Controllers[控制器]
                Routes[路由管理]
                Middleware[中间件]
            end
        end

        subgraph 支撑层
            Utils[工具包]
            Auth[JWT认证]
            Wire[依赖注入]
            Gen[wire生成器]
        end
    end

    %% 依赖关系
    Config -->|配置| DB
    Config -->|配置| Cache
    Config -->|配置| OSS
    Config -->|配置| AIClient
    Config -->|配置| Logger
    
    DB -->|持久化| A
    Cache -->|缓存| B
    OSS -->|文件存储| C
    AIClient -->|调用| D

    Controllers -->|调用| Services
    Services -->|依赖| Repositories
    Repositories -->|操作| DB
    Repositories -->|操作| Cache
    
    Middleware -->|鉴权| Auth
    Middleware -->|日志| Logger
    
    %% Wire依赖注入关系
    Wire -->|自动装配| Controllers
    Wire -->|自动装配| Services
    Wire -->|自动装配| Repositories
    Wire -->|自动装配| DB
    Wire -->|自动装配| Cache
    Wire -->|自动装配| OSS
    Wire -->|自动装配| AIClient
    Wire -->|自动装配| Auth
    
    Gen -.->|生成代码| Wire
    
    Utils --> 核心层
    Utils --> 基础设施层
    
    %% 框架依赖
    Routes --> Gin
    Controllers --> Gin
    Middleware --> Gin
    Auth --> JWT

    classDef external fill:#f9d5bb,stroke:#666;
    classDef infra fill:#d5e8d4,stroke:#338833;
    classDef core fill:#e1d5e7,stroke:#884488;
    classDef support fill:#fff2cc,stroke:#dd8800;
    classDef framework fill:#dae8fc,stroke:#3366cc;
    classDef generator fill:#ffcccc,stroke:#ff6666;
    
    class A,B,C,D external;
    class Config,Logger,DB,Cache,OSS,AIClient infra;
    class Models,Repositories,Services,Controllers,Routes,Middleware core;
    class Utils,Auth,Wire support;
    class Gin,JWT framework;
    class Gen generator;
```

# 注意事项

运行前记得生成依赖注入的wire代码

```shell
go get -u ./... && go mod tidy && go tool github.com/google/wire/cmd/wire ./...
```



