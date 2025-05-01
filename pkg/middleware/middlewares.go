package middleware

import (
	"Art-Design-Backend/model/entity"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"reflect"
	"regexp"
)

type Middlewares struct {
	Db    *gorm.DB             // 数据库
	Redis *redisx.RedisWrapper // redis
	Jwt   *jwt.JWT             // jwt
}

// registerValidator 注册校验器错误返回信息
func registerValidator(model interface{}) {
	var t reflect.Type
	// 解读并缓存标签
	t = reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldLabels[field.Name] = field.Tag.Get("label")
	}
}

// RegisterValidator 注册全局请求校验器
func RegisterValidator() {
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
			pass := fl.Field().String()
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pass)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(pass)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pass)
			return hasUpper && hasLower && hasNumber
		})
		if err != nil {
			zap.L().Fatal("自定义校验器注册失败")
			return
		}
	}
}

func NewMiddlewares(db *gorm.DB, redis *redisx.RedisWrapper, jwt *jwt.JWT) *Middlewares {
	RegisterValidator()
	fieldLabels = make(map[string]string)
	{
		// 注册校验器错误返回信息
		registerValidator(entity.User{})
	}
	return &Middlewares{
		Db:    db,
		Redis: redis,
		Jwt:   jwt,
	}
}
