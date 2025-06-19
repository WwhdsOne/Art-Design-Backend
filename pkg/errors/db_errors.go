package errors

import (
	"fmt"
)

type DBError struct {
	Message string
	Err     error
}

func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("db error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("db error: %s", e.Message)
}

// Unwrap 支持 errors.Unwrap() 解包
func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError 创建基础 DBError（无嵌套错误）
func NewDBError(message string) error {
	return &DBError{
		Message: message,
	}
}

// WrapDBError 包装原始错误（推荐使用）
func WrapDBError(err error, message string) error {
	if err == nil {
		return NewDBError(message)
	}
	return &DBError{
		Message: message,
		Err:     err,
	}
}
