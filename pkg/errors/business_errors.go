package errors

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(message string) error {
	return &BusinessError{
		Message: message,
	}
}
