package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AutoMigrate 自动迁移
func AutoMigrate(_ *gorm.DB) {
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
	//// 7. AI模型供应商
	//db.AutoMigrate(&entity.AIProvider{})
	////8. 知识库
	//db.AutoMigrate(&entity.ChunkVector{})
	//db.AutoMigrate(&entity.FileChunk{})
	//db.AutoMigrate(&entity.KnowledgeBaseFile{})
	//db.AutoMigrate(&entity.KnowledgeBase{})
	//db.AutoMigrate(&entity.KnowledgeBaseFileRel{})
	//// 9. 会话
	//db.AutoMigrate(&entity.Conversation{})
	//db.AutoMigrate(&entity.Message{})
}

// snowflakeIDFieldsMap 存储类型和对应的ID字段名（缓存，提高效率）
var snowflakeIDFieldsMap sync.Map // key: reflect.Type, value: string

// snowflakeIDPlugin GORM插件实现，用于自动填充雪花ID
type snowflakeIDPlugin struct{}

// Name 插件名称（GORM要求实现）
func (p *snowflakeIDPlugin) Name() string {
	return "snowflake_id_plugin"
}

// initialize 注册插件钩子函数
// 在创建前（Before Create）触发 generateID
func (p *snowflakeIDPlugin) initialize(db *gorm.DB) (err error) {
	err = db.Callback().Create().
		Before("gorm:create").
		Register("generate_snowflake_id", p.generateID)
	return
}

// detectSnowflakeIDField 递归查找结构体中的主键字段名（如ID），用于自动生成雪花ID
func detectSnowflakeIDField(t reflect.Type) string {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 如果是匿名嵌套字段（如 BaseModel），递归查找
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if nestedField := detectSnowflakeIDField(field.Type); nestedField != "" {
				return nestedField
			}
		}

		// 查找 int64 类型的主键字段
		if field.Type.Kind() == reflect.Int64 {
			gormTag := field.Tag.Get("gorm")
			if strings.Contains(gormTag, "primaryKey") {
				return field.Name
			}
			// 如果字段名是 ID（大小写不敏感），也认定为主键
			if strings.EqualFold(field.Name, "ID") {
				return field.Name
			}
		}
	}
	return ""
}

// getFieldByName 获取字段值，支持递归嵌套结构体中查找字段
func getFieldByName(v reflect.Value, name string) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return reflect.Value{}
	}
	field := v.FieldByName(name)
	if field.IsValid() {
		return field
	}
	// 查找匿名字段中的目标字段
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		if structField.Anonymous && v.Field(i).Kind() == reflect.Struct {
			if f := getFieldByName(v.Field(i), name); f.IsValid() {
				return f
			}
		}
	}
	return reflect.Value{}
}

// setID 设置字段值为雪花ID，仅在字段为 int64 且值为 0 时设置
func (p *snowflakeIDPlugin) setID(db *gorm.DB, fieldName string) {
	field := getFieldByName(reflect.ValueOf(db.Statement.Model), fieldName)
	if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.Int64 {
		return
	}
	if field.Int() == 0 {
		field.SetInt(utils.GenerateSnowflakeID()) // 设置生成的雪花ID
	}
}

// generateID 插入前自动生成雪花ID（插件核心逻辑）
func (p *snowflakeIDPlugin) generateID(db *gorm.DB) {
	model := db.Statement.Model
	modelType := reflect.TypeOf(model)

	// 如果是切片类型（批量插入），取元素类型
	if modelType.Kind() == reflect.Slice {
		modelType = modelType.Elem()
	}
	// 如果是指针类型，取实际结构体类型
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	// 不是结构体则跳过
	if modelType.Kind() != reflect.Struct {
		return
	}

	var fieldName string

	// 从缓存中查找是否已记录字段名
	if fieldNameRaw, ok := snowflakeIDFieldsMap.Load(modelType); ok {
		fieldName = fieldNameRaw.(string)
	} else {
		// 首次分析结构体字段，找出主键字段名
		fieldName = detectSnowflakeIDField(modelType)
		snowflakeIDFieldsMap.Store(modelType, fieldName)
	}

	// 找不到主键ID字段，跳过处理（比如联表结构体）
	if fieldName == "" {
		return
	}

	// 设置主键字段为雪花ID
	p.setID(db, fieldName)
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

func (z *zapGormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Info {
		z.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (z *zapGormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Warn {
		z.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (z *zapGormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if z.logLevel >= logger.Error {
		z.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (z *zapGormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
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
	m := cfg.PostgreSQL
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
