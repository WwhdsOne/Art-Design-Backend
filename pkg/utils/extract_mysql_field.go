package utils

import "strings"

// ExtractMySQLUniqueField 从MySQL错误信息中提取唯一约束冲突的字段名
func ExtractMySQLUniqueField(errMsg string) string {
	// MySQL错误格式示例:
	// "Error 1062: Duplicate entry 'value' for key 'field_name'"
	// 或 "Error 1062: Duplicate entry 'value' for key 'table_name.field_name'"

	if !strings.Contains(errMsg, "Duplicate entry") {
		return ""
	}

	parts := strings.Split(errMsg, "for key ")
	if len(parts) < 2 {
		return ""
	}

	field := strings.Trim(parts[1], "'`\"")

	// 去除可能的表名前缀 (table_name.field_name -> field_name)
	if dotIndex := strings.LastIndex(field, "."); dotIndex != -1 {
		field = field[dotIndex+1:]
	}

	return field
}
