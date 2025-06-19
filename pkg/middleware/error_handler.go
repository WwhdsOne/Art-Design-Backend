package middleware

import (
	myerrors "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/result"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"sync"
)

var (
	fieldLabels = make(map[string]map[string]string)
	mu          sync.RWMutex
)

func getFieldLabel(obj interface{}, fieldName string) string {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typeName := t.PkgPath() + "." + t.Name()

	// 读锁读取
	mu.RLock()
	labelMap, ok := fieldLabels[typeName]
	mu.RUnlock()

	if !ok {
		// 初始化标签映射
		labelMap = make(map[string]string)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			labelMap[field.Name] = field.Tag.Get("label")
		}

		// 加写锁存入（双重检查防止重复写）
		mu.Lock()
		if _, doubleCheck := fieldLabels[typeName]; !doubleCheck {
			fieldLabels[typeName] = labelMap
		} else {
			labelMap = fieldLabels[typeName]
		}
		mu.Unlock()
	}

	if label, ok := labelMap[fieldName]; ok && label != "" {
		return label
	}
	return fieldName
}

// handleValidationErrors 处理验证错误
func handleValidationErrors(c *gin.Context, errs validator.ValidationErrors) {
	// 只返回第一个错误即可
	e := errs[0]
	obj := c.MustGet(gin.BindKey)
	fieldName := getFieldLabel(obj, e.Field())

	var message string
	switch e.Tag() {
	case "required":
		message = "不能为空"
	case "min":
		message = fmt.Sprintf("长度不能少于%s个字符", e.Param())
	case "max":
		message = fmt.Sprintf("长度不能超过%s个字符", e.Param())
	case "email":
		message = "必须是有效的邮箱格式"
	case "e164":
		message = "必须是国际电话号码格式（如+8613812345678）"
	case "alphanumunicode":
		message = "只能包含字母、数字或中文"
	case "oneof":
		message = fmt.Sprintf("必须是以下值之一: %s", strings.Replace(e.Param(), " ", " 或 ", -1))
	case "strongpassword":
		message = "必须包含大小写字母和数字"
	default:
		message = fmt.Sprintf("不符合校验规则(%s)", e.Tag())
	}
	result.FailWithMessage(fmt.Sprintf("%s %s", fieldName, message), c)
}

// handleDBErrors 处理Gorm错误
func handleDBErrors(c *gin.Context, gormErr *myerrors.DBError) {
	result.FailWithMessage(gormErr.Message, c)
}

// handleGenericErrors 处理除验证错误和Gorm错误之外的所有其他错误
func handleGenericErrors(c *gin.Context, err error) {
	zap.L().Error("request failed",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Error(err),
	)
	result.FailWithMessage(err.Error(), c)
}

// ErrorHandlerMiddleware 错误处理中间件
func (m *Middlewares) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		// 捕获错误后立即中止请求
		c.Abort()

		var veErr validator.ValidationErrors
		var dbErr *myerrors.DBError
		for _, ginErr := range c.Errors {

			switch {
			case errors.As(ginErr.Err, &veErr):
				handleValidationErrors(c, veErr)
				return
			case errors.As(ginErr.Err, &dbErr):
				handleDBErrors(c, dbErr)
				return
			default:
				handleGenericErrors(c, ginErr.Err)
				return
			}
		}
	}
}
