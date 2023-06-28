package errors

type InvalidApplicationSecretsError struct {
	Message string
}

func (e *InvalidApplicationSecretsError) Error() string {
	return e.Message
}

func NewInvalidApplicationSecretsError(
	message string,
) *InvalidApplicationSecretsError {
	return &InvalidApplicationSecretsError{
		Message: message,
	}
}
