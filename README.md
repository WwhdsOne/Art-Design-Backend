# Art-Design-Pro个人后端项目

[前端项目地址](https://github.com/Daymychen/art-design-pro)

# 项目结构

```shell
.
├── README.md
├── build.sh
├── cmd
│   └── app
│       ├── main.go
│       ├── run_server.go
│       ├── wire.go
│       └── wire_gen.go
├── config
│   ├── init_config.go
│   ├── init_gorm.go
│   ├── init_http_server.go
│   ├── init_jwt.go
│   ├── init_oss_client.go
│   ├── init_redis.go
│   ├── init_validator.go
│   └── init_zaplog.go
├── configs
│   ├── config-prod.yaml
│   └── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── controller
│   │   ├── auth.go
│   │   ├── controllers.go
│   │   ├── menu.go
│   │   ├── role.go
│   │   └── user.go
│   ├── model
│   │   ├── base
│   │   │   ├── base_model.go
│   │   │   ├── long_string_id.go
│   │   │   ├── page_request.go
│   │   │   └── page_response.go
│   │   ├── entity
│   │   │   ├── menu.go
│   │   │   ├── operation_log.go
│   │   │   ├── role.go
│   │   │   ├── role_menus.go
│   │   │   ├── user.go
│   │   │   └── user_roles.go
│   │   ├── query
│   │   │   └── user.go
│   │   ├── request
│   │   │   ├── change_password.go
│   │   │   ├── login.go
│   │   │   ├── menu.go
│   │   │   ├── register_user.go
│   │   │   ├── role.go
│   │   │   └── user.go
│   │   └── resp
│   │       ├── menu.go
│   │       ├── role.go
│   │       └── user.go
│   ├── repository
│   │   ├── gorm_transaction.go
│   │   ├── menu.go
│   │   ├── repos.go
│   │   ├── role.go
│   │   └── user.go
│   └── service
│       ├── auth.go
│       ├── menu.go
│       ├── role.go
│       └── user.go
└── pkg
    ├── aliyun
    │   └── oss_client.go
    ├── authutils
    │   └── auth_utils.go
    ├── constant
    │   ├── default_avatar.go
    │   ├── oss_directory.go
    │   ├── redis_key.go
    │   └── table_name.go
    ├── errors
    │   └── db_errors.go
    ├── jwt
    │   └── jwt.go
    ├── middleware
    │   ├── auth.go
    │   ├── error_handler.go
    │   ├── middlewares.go
    │   └── operation_log.go
    ├── redisx
    │   └── redisx.go
    ├── response
    │   └── response.go
    └── utils
        ├── http_param_parser.go
        ├── parse_duration.go
        ├── random_value.go
        ├── snowflake_id.go
        └── uuid.go

25 directories, 71 files
```

# 注意事项

运行前记得生成依赖注入的wire代码

```shell
go tool github.com/google/wire/cmd/wire ./...
```

