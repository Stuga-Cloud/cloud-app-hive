package errors

type InvalidApplicationScalabilitySpecificationsError struct {
	Message string
}

func (e InvalidApplicationScalabilitySpecificationsError) Error() string {
	return e.Message
}

func NewInvalidApplicationScalabilitySpecificationsError(message string) InvalidApplicationScalabilitySpecificationsError {
	return InvalidApplicationScalabilitySpecificationsError{
		Message: message,
	}
}
