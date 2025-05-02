package config

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
	"time"
)

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var snowflakeIdFieldsMap = make(map[interface{}]string)

// AutoMigrate 自动迁移
func AutoMigrate(db *gorm.DB) {
	// 1. 操作日志
	//db.AutoMigrate(&entity.OperationLog{})
	//// 2. 用户
	//db.AutoMigrate(&entity.User{})
	//// 3. 角色
	//db.AutoMigrate(&entity.Role{})
	//// 4. 菜单
	db.AutoMigrate(&entity.Menu{})
}

// snowflakeIDPlugin GORM插件实现
type snowflakeIDPlugin struct{}

func (p *snowflakeIDPlugin) Name() string {
	return "snowflake_id_plugin"
}

// initialize 初始化数据库
// 雪花ID生成插件
func (p *snowflakeIDPlugin) initialize(db *gorm.DB) (err error) {
	err = db.Callback().Create().
		Before("gorm:create").
		Register("generate_snowflake_id", p.generateID)
	return
}

func (p *snowflakeIDPlugin) generateID(db *gorm.DB) {

	// 获取模型类型
	modelType := reflect.TypeOf(db.Statement.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 如果存在
	if loadedFieldName, exist := snowflakeIdFieldsMap[modelType]; exist {
		// 获取字段
		fieldName := loadedFieldName
		modelValue := reflect.ValueOf(db.Statement.Model).Elem()
		field := modelValue.FieldByName(fieldName)
		// 如果字段为0，则设置雪花ID
		if field.Int() == 0 {
			field.SetInt(utils.GenerateSnowflakeId())
		}
	}
	// 不存在直接返回
	return
}

// zapGormLogger 实现 gorm.Logger.Interface
type zapGormLogger struct {
	zapLogger *zap.Logger
	logLevel  logger.LogLevel
}

func (z *zapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *z
	newLogger.logLevel = level
	return &newLogger
}

func (z *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Info {
		z.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (z *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Warn {
		z.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (z *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Error {
		z.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (z *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if z.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("affected_rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		z.zapLogger.Error("SQL Trace", fields...)
	} else {
		z.zapLogger.Debug("SQL Trace", fields...)
	}
}

// RegisterIDField 注册需要自动生成ID的模型和字段
func registerIDField(model interface{}, fieldName string) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	snowflakeIdFieldsMap[model] = fieldName
}

func NewGorm(cfg *Config, log *zap.Logger) (DB *gorm.DB) {
	m := cfg.Mysql
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,     //用户名
		m.Password, //密码
		m.Host,     //地址
		m.Port,     //端口
		m.Database, //数据库
	)
	// 创建 Zap 日志适配器
	gormLogger := &zapGormLogger{
		zapLogger: log,
		logLevel:  logger.Info, // 设置默认日志级别
	}

	// 连接数据库
	DB, err := gorm.Open(mysql.Open(ds), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭自动迁移外键创建
	})

	if err != nil {
		log.Fatal("数据库连接失败")
		return
	}

	// 注册模型应当自动填充雪花ID的字段
	{
		registerIDField(&entity.User{}, "ID")
		registerIDField(&entity.OperationLog{}, "ID")
		registerIDField(&entity.Menu{}, "ID")
	}

	// 雪花ID插件
	snowflakeID := &snowflakeIDPlugin{}
	if err = snowflakeID.initialize(DB); err != nil {
		log.Fatal("雪花ID插件注册失败", zap.Error(err))
		return
	}
	// 自动迁移
	AutoMigrate(DB)
	return
}
