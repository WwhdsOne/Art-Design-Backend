package main

import (
	"Art-Design-Backend/core"
)

//func init() {
//	// 方法5：手动创建东八区时区（固定偏移+8小时）
//	// 参数说明：
//	// - "CST"：时区名称（这里用China Standard Time缩写）
//	// - 8*60*60：东八区偏移秒数（8小时×60分钟×60秒）
//	loc := time.FixedZone("CST", 8*60*60)
//	time.Local = loc // 设置全局默认时区
//	log.Printf("Timezone set to: %s (UTC+8)\n", loc)
//}

func main() {
	core.RunServer()
}
