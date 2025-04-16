package initialize

import (
	"Art-Design-Backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func encoder(z *config.Zap) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		NameKey:       "name",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: z.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			// 为时间添加绿色
			encoder.AppendString("\033[33m" + z.Prefix + "\033[0m \033[32m" +
				t.Format("2006-01-02 15:04:05.000") + "\033[0m")
		},
		EncodeLevel:    levelEncoder(z.EncodeLevel),
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	if z.Format == "json" {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)

}

// LevelEncoder 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author [SliverHorn](https://github.com/SliverHorn)
func levelEncoder(encodeLevel string) zapcore.LevelEncoder {
	switch {
	case encodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case encodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case encodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case encodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// InitLogger 初始化 Zap Logger
func InitLogger(c *config.Config) (logger *zap.Logger) {

	z := c.Zap
	// 创建编码器
	zapEncoder := encoder(&z)

	// 创建输出目标，这里输出到控制台
	writer := zapcore.AddSync(os.Stdout)

	// 创建日志级别
	logLevel := zapcore.DebugLevel // 可以根据需要设置其他级别

	// 创建 Core
	core := zapcore.NewCore(zapEncoder, writer, logLevel)

	logger = zap.New(core)

	return
}
