package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func RunServer() {
	// 展示神兽
	displayGodAnimal()
	isDev := os.Getenv("ENV") == "DEV"
	// 设置GIN模式
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}
	app := wireApp()
	app.GinServer()
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
