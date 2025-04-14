package initialize

import (
	"Art-Design-Backend/global"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

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
			global.Logger.Fatal("自定义校验器注册失败")
			return
		}
	}
}
