package errors

type InvalidApplicationCannotVerticallyScaleBecauseMaxSpecsError struct {
	Message string
}

func (e *InvalidApplicationCannotVerticallyScaleBecauseMaxSpecsError) Error() string {
	return e.Message
}

func NewInvalidApplicationCannotVerticallyScaleBecauseMaxSpecsError(
	message string,
) *InvalidApplicationCannotVerticallyScaleBecauseMaxSpecsError {
	return &InvalidApplicationCannotVerticallyScaleBecauseMaxSpecsError{
		Message: message,
	}
}