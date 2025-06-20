package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
	"strings"
	"sync"
	"time"
)

// AutoMigrate 自动迁移
func AutoMigrate(db *gorm.DB) {
	//// 1. 操作日志
	//db.AutoMigrate(&entity.OperationLog{})
	//// 2. 用户
	//db.AutoMigrate(&entity.User{})
	//// 3. 角色
	//db.AutoMigrate(&entity.Role{})
	//// 4. 菜单
	//db.AutoMigrate(&entity.Menu{})
	//// 5. 数字识别
	//db.AutoMigrate(&entity.DigitPredict{})
	//// 6. AI模型
	//db.AutoMigrate(&entity.AIModel{})
}

// snowflakeIdFieldsMap 存储类型和对应的ID字段名
var snowflakeIdFieldsMap sync.Map // key: reflect.Type, value: string

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

func detectSnowflakeIDField(t reflect.Type) string {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 优先识别 `gorm:"primaryKey"` 或 `gorm:"autoSnowflake"` 之类的标记
		if field.Type.Kind() == reflect.Int64 {
			if gormTag := field.Tag.Get("gorm"); strings.Contains(gormTag, "primaryKey") {
				return field.Name
			}
			if strings.EqualFold(field.Name, "ID") {
				return field.Name
			}
		}
	}
	return ""
}

func (p *snowflakeIDPlugin) setID(db *gorm.DB, fieldName string) {
	modelValue := reflect.ValueOf(db.Statement.Model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	field := modelValue.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.Int64 {
		return
	}
	if field.Int() == 0 {
		field.SetInt(utils.GenerateSnowflakeId())
	}
}

func (p *snowflakeIDPlugin) generateID(db *gorm.DB) {
	model := db.Statement.Model
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 缓存中查找字段
	if fieldNameRaw, ok := snowflakeIdFieldsMap.Load(modelType); ok {
		p.setID(db, fieldNameRaw.(string))
		return
	}

	// 没找到就首次分析结构体，记录字段
	fieldName := detectSnowflakeIDField(modelType)
	if fieldName != "" {
		snowflakeIdFieldsMap.Store(modelType, fieldName)
		p.setID(db, fieldName)
	}
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

func (z *zapGormLogger) Info(c context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Info {
		z.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (z *zapGormLogger) Warn(c context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Warn {
		z.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (z *zapGormLogger) Error(c context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Error {
		z.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (z *zapGormLogger) Trace(c context.Context, begin time.Time, fc func() (string, int64), err error) {
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

func InitGorm(cfg *config.Config, log *zap.Logger) (DB *gorm.DB) {
	m := cfg.PostgreSql
	// 构建PostgreSQL连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		m.Host,
		m.Port,
		m.User,
		m.Password,
		m.Database,
	)
	// 创建 Zap 日志适配器
	gormLogger := &zapGormLogger{
		zapLogger: log,
		logLevel:  logger.Info, // 设置默认日志级别
	}

	// 连接数据库
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭自动迁移外键创建
	})

	if err != nil {
		log.Fatal("数据库连接失败")
		return
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
