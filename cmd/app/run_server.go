package main

import (
	"fmt"
)

func RunServer() {
	// 展示神兽
	displayGodAnimal()
	app := wireApp()
	app.InitGinServer()
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
