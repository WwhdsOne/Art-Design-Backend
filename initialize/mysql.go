package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/global"
	"Art-Design-Backend/model/entity"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func AutoMigrate(db *gorm.DB) {
	// 自动迁移
	// 1. 用户
	db.AutoMigrate(&entity.User{})
	// 2. 操作日志
	db.AutoMigrate(&entity.OperationLog{})
}

func InitDB(cfg *config.Config) *gorm.DB {
	m := cfg.Mysql
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,     //用户名
		m.Password, //密码
		m.Host,     //地址
		m.Port,     //端口
		m.Database, //数据库
	)
	// 连接数据库
	Db, err := gorm.Open(mysql.Open(ds), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印执行sql
	})

	if err != nil {
		global.Logger.Error("数据库连接失败")
		return nil
	}
	// 自动迁移
	// AutoMigrate(Db)
	return Db
}
