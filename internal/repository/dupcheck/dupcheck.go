package dupcheck

import (
	"Art-Design-Backend/pkg/errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Field 描述一个需要检查重复的字段
type Field struct {
	Column  string // 数据库列名
	ErrMsg  string // 重复时的错误信息
	Value   any    // 字段值（nil 跳过）
	IsNil   bool   // 为 true 时跳过该字段
}

// Check 在指定表中检查字段重复，excludeID > 0 时排除自身
func Check(db *gorm.DB, table string, excludeID int64, fields []Field) error {
	conditions := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))
	errMsgs := make(map[string]string) // alias → error message

	for _, f := range fields {
		if f.IsNil || f.Value == nil {
			continue
		}
		if v, ok := f.Value.(string); ok && v == "" {
			continue
		}

		alias := f.Column + "_exists"
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS(SELECT 1 FROM %s WHERE %s = ? %s) AS %s",
			table, f.Column, excludeClause(excludeID), alias,
		))
		args = append(args, f.Value)
		errMsgs[alias] = f.ErrMsg
	}

	if len(conditions) == 0 {
		return nil
	}

	// 动态构建 scan 目标
	scanTarget := make(map[string]bool, len(conditions))
	for _, c := range conditions {
		// 提取 alias（AS 后面的部分）
		parts := strings.SplitN(c, " AS ", 2)
		if len(parts) == 2 {
			scanTarget[strings.TrimSpace(parts[1])] = false
		}
	}

	query := "SELECT " + strings.Join(conditions, ", ")
	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return errors.WrapDBError(err, "检查重复属性失败")
	}
	defer rows.Close()

	if rows.Next() {
		columns, _ := rows.Columns()
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return errors.WrapDBError(err, "扫描重复检查结果失败")
		}

		for i, col := range columns {
			if val, ok := values[i].(bool); ok && val {
				if msg, exists := errMsgs[col]; exists {
					return errors.NewDBError(msg)
				}
			}
		}
	}

	return nil
}

func excludeClause(id int64) string {
	if id > 0 {
		return fmt.Sprintf("AND id != %d", id)
	}
	return ""
}
