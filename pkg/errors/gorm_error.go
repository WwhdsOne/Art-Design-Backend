package errors

type GormError struct {
	Message string
}

func (e *GormError) Error() string {
	return e.Message
}

func NewGormError(message string) error {
	return &GormError{
		Message: message,
	}
}
