package errors

import "fmt"

type CacheError struct {
	Message string // 语义信息
	Err     error  // 原始底层错误（可选）
}

// NewCacheError 创建简单消息的 CacheError
func NewCacheError(message string) error {
	return &CacheError{
		Message: message,
	}
}

// WrapCacheError 包装原始错误的 CacheError（推荐使用）
func WrapCacheError(err error, message string) error {
	return &CacheError{
		Message: message,
		Err:     err,
	}
}

func (e *CacheError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cache error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("cache error: %s", e.Message)
}

// Unwrap 支持 errors.Unwrap
func (e *CacheError) Unwrap() error {
	return e.Err
}
