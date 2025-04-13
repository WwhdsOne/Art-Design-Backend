package core

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/global"
	"Art-Design-Backend/initialize"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert/yaml"
	"os"
)

func initGlobal(cfg *config.Config) {
	// 初始化日志
	global.Logger = initialize.InitLogger()
	// 初始化数据库
	global.DB = initialize.InitDB(cfg)
	// 初始化Redis
	global.Redis = initialize.InitRedis(cfg)
	// 初始化JWT
	global.JWT = initialize.InitJWT(cfg)
	// 初始化OSS客户端
	global.OSSClient = initialize.InitOSSClient(cfg)
}

func readConfig() (cfg *config.Config) {
	var data []byte
	var err error
	if os.Getenv("ENV") == "DEV" {
		data, err = os.ReadFile("conf/config.yaml")
	} else {
		data, err = os.ReadFile("config.dev.yaml")
		gin.SetMode(gin.ReleaseMode)
	}
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Printf("配置如下 : %v\n", cfg)
	return
}

func RunServer() {
	// 展示神兽
	displayGodAnimal()
	// 读取配置文件
	cfg := readConfig()
	// 初始化全局变量
	initGlobal(cfg)
	// 初始化GIN
	r := initialize.InitGin()
	// 初始化路由
	initialize.InitRouter(r)
	// 启动
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
