package core

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/global"
	"Art-Design-Backend/initialize"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert/yaml"
	"log"
	"os"
)

func initGlobal(cfg *config.Config) {
	// 初始化日志
	global.Logger = initialize.InitLogger(cfg)
	// 初始化数据库
	global.DB = initialize.InitDB(cfg)
	// 初始化Redis
	global.Redis = initialize.InitRedis(cfg)
	// 初始化JWT
	global.JWT = initialize.InitJWT(cfg)
	// 初始化OSS客户端
	global.OSSClient = initialize.InitOSSClient(cfg)
}

func readConfig(isDev bool) (cfg *config.Config) {
	var data []byte
	var err error
	if !isDev {
		data, err = os.ReadFile("config.yaml")
	} else {
		data, err = os.ReadFile("conf/config.yaml")

	}
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("配置如下 : %v\n", cfg)
	return
}

func RunServer() {
	// 展示神兽
	displayGodAnimal()
	isDev := os.Getenv("ENV") == "DEV"
	// 读取配置文件
	cfg := readConfig(isDev)
	// 初始化全局变量
	initGlobal(cfg)
	// 设置GIN模式
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}
	// 初始化模型信息
	initialize.InitModelInfo()
	// 初始化GIN
	r := initialize.InitGin()
	// 初始化校验器
	initialize.RegisterValidator()
	// 初始化路由
	initialize.InitRouter(r)
	// 启动
	global.Logger.Info(fmt.Sprintf("%s服务启动成功，端口号为：%s", cfg.Server.App, cfg.Server.Port))
	r.Run(cfg.Server.Port)
}

func displayGodAnimal() {
	fmt.Println(`
                           ┏━┓     ┏━┓
                          ┏┛ ┻━━━━━┛ ┻┓
                          ┃　　　　　　 ┃
                          ┃　　　━　　　┃
                          ┃　┳┛　  ┗┳　┃
                          ┃　　　　　　 ┃
                          ┃　　　┻　　　┃
                          ┃　　　　　　 ┃
                          ┗━┓　　　┏━━━┛
                            ┃　　　┃   神兽保佑
                            ┃　　　┃   代码无BUG！
                            ┃　　　┗━━━━━━━━━┓
                            ┃　　　　　　　    ┣┓
                            ┃　　　　         ┏┛
                            ┗━┓ ┓ ┏━━━┳ ┓ ┏━┛
                              ┃ ┫ ┫   ┃ ┫ ┫
                              ┗━┻━┛   ┗━┻━┛`)
}
