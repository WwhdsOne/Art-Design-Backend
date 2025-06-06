package middleware

import (
	"Art-Design-Backend/internal/model/request"
	myerrors "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/result"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// fieldLabels 字段标签映射
var fieldLabels = make(map[string]string)

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

func getFieldLabel(obj interface{}, fieldName string) string {
	t := reflect.TypeOf(obj)

	// 确保处理的是结构体类型，无论是值类型还是指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "" // 如果不是结构体，返回空字符串或适当的错误处理
	}

	// 检查是否有预定义的标签映射
	if label, ok := fieldLabels[fieldName]; ok && label != "" {
		return label
	}

	// 获取结构体字段
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return fieldName // 如果字段不存在，返回字段名
	}

	// 获取标签
	label := field.Tag.Get("label")
	if label == "" {
		return fieldName // 如果标签为空，返回字段名
	}
	return label
}

// handleValidationErrors 处理验证错误
func handleValidationErrors(c *gin.Context, ginErr *gin.Error) bool {
	var errs validator.ValidationErrors
	if errors.As(ginErr.Err, &errs) {
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
		result.FailWithMessage(fmt.Sprintf("%s%s", fieldName, message), c)
		return true
	}
	return false
}

// handleDBErrors 处理Gorm错误
func handleDBErrors(c *gin.Context, ginErr *gin.Error) bool {
	var gormErr *myerrors.DBError
	if errors.As(ginErr.Err, &gormErr) {
		result.FailWithMessage(gormErr.Message, c)
		return true
	}
	return false
}

// handleGenericErrors 处理除验证错误和Gorm错误之外的所有其他错误
func handleGenericErrors(c *gin.Context, ginErr *gin.Error) bool {
	// 返回通用错误消息
	result.FailWithMessage(ginErr.Err.Error(), c)
	return true
}

// ErrorHandlerMiddleware 错误处理中间件
func (m *Middlewares) ErrorHandlerMiddleware() gin.HandlerFunc {
	{
		// 注册校验器错误返回信息
		registerValidator(request.ChangePassword{})
		registerValidator(request.DigitPredict{})
		registerValidator(request.Login{})
		registerValidator(request.Menu{})
		registerValidator(request.RegisterUser{})
		registerValidator(request.Role{})
		registerValidator(request.RoleMenuBinding{})
		registerValidator(request.User{})
		registerValidator(request.MenuAuth{})
	}
	return func(c *gin.Context) {
		c.Next() // 先调用后续的处理函数

		// 检查是否有错误发生
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				// 处理验证错误
				if handleValidationErrors(c, ginErr) {
					c.Abort()
					return
				}

				// 处理Gorm错误
				if handleDBErrors(c, ginErr) {
					c.Abort()
					return
				}

				// 处理其他所有错误
				if handleGenericErrors(c, ginErr) {
					c.Abort()
					return
				}
			}
		}
	}
}
