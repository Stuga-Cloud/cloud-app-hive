package errors

type InvalidApplicationEnvironmentVariablesError struct {
	Message string
}

func (e *InvalidApplicationEnvironmentVariablesError) Error() string {
	return e.Message
}

func NewInvalidApplicationEnvironmentVariablesError(
	message string,
) *InvalidApplicationEnvironmentVariablesError {
	return &InvalidApplicationEnvironmentVariablesError{
		Message: message,
	}
}
