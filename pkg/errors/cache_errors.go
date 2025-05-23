package errors

type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

func NewCacheError(message string) error {
	return &CacheError{
		Message: message,
	}
}
