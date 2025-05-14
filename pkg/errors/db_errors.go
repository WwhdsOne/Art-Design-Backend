package errors

type DBError struct {
	Message string
}

func (e *DBError) Error() string {
	return e.Message
}

func NewDBError(message string) error {
	return &DBError{
		Message: message,
	}
}
