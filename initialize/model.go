package initialize

import (
	"Art-Design-Backend/model/entity"
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/pkg/middleware"
	"reflect"
)

// snowflakeIdFieldsMap 存储需要生成ID的模型和字段名
var snowflakeIdFieldsMap = make(map[interface{}]string)

// 注册模型相关

func InitModelInfo() {
	// 1. 注册雪花ID生成插件
	{
		registerIDField(&entity.User{}, "ID")
		registerIDField(&entity.OperationLog{}, "ID")
	}
	// 2. 校验器错误返回信息
	{
		registerValidator(&request.User{})
	}

}

// RegisterIDField 注册需要自动生成ID的模型和字段
// 参数：
//
//	model：需要生成ID的模型
//	fieldName：需要生成ID的字段名
func registerIDField(model interface{}, fieldName string) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	snowflakeIdFieldsMap[model] = fieldName
}

func registerValidator(model interface{}) {
	var t reflect.Type
	// 解读并缓存标签
	t = reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		middleware.FieldLabels[field.Name] = field.Tag.Get("label")
	}
}
