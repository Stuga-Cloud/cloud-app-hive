package errors

type InvalidApplicationContainerSpecificationsError struct {
	Message string
}

func (e InvalidApplicationContainerSpecificationsError) Error() string {
	return e.Message
}

func NewInvalidApplicationContainerSpecificationsError(message string) InvalidApplicationContainerSpecificationsError {
	return InvalidApplicationContainerSpecificationsError{
		Message: message,
	}
}
